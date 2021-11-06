package exhaustive

import (
	"fmt"
	"go/ast"
	"go/types"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/ast/inspector"
)

// nodeVisitor is similar to the visitor function used by Inspector.WithStack,
// except that it returns an additional value: a short description of
// the result of this node visit.
//
// The result is typically useful in debugging or in unit tests to check
// that the nodeVisitor function took the expected code path.
type nodeVisitor func(n ast.Node, push bool, stack []ast.Node) (proceed bool, result string)

// Result values returned by a node visitor constructed via switchStmtChecker.
const (
	resultNotPush              = "not push"
	resultGeneratedFile        = "generated file"
	resultNoSwitchTag          = "no switch tag"
	resultTagNotValue          = "switch tag not value type"
	resultTagNotNamed          = "switch tag not named type"
	resultTagNoPkg             = "switch tag does not belong to regular package"
	resultTagNotEnum           = "switch tag not known enum type"
	resultSwitchIgnoreComment  = "switch statement has ignore comment"
	resultEnumMembersAccounted = "requisite enum members accounted for"
	resultDefaultCaseSuffices  = "default case presence satisfies exhaustiveness"
	resultReportedDiagnostic   = "reported diagnostic"
)

// switchStmtChecker returns a node visitor that checks exhaustiveness
// of enum switch statements for the supplied pass, and reports diagnostics for
// switch statements that are non-exhaustive.
func switchStmtChecker(pass *analysis.Pass, cfg config) nodeVisitor {
	comments := make(map[*ast.File]ast.CommentMap)
	generated := make(map[*ast.File]bool)

	return func(n ast.Node, push bool, stack []ast.Node) (bool, string) {
		if !push {
			// we only inspect things on the way down, not up.
			return true, resultNotPush
		}

		file := stack[0].(*ast.File)
		sw := n.(*ast.SwitchStmt)

		// Determine if the file is a generated file, and save the result.
		// If it is a generated file, don't check the file.
		if _, ok := generated[file]; !ok {
			generated[file] = isGeneratedFile(file)
		}
		if generated[file] && !cfg.checkGeneratedFiles {
			// don't check this file.
			return true, resultGeneratedFile
		}

		if _, ok := comments[file]; !ok {
			comments[file] = ast.NewCommentMap(pass.Fset, file, file.Comments)
		}
		if containsIgnoreDirective(comments[file].Filter(sw).Comments()) {
			// skip checking due to ignore directive
			return true, resultSwitchIgnoreComment
		}

		if sw.Tag == nil {
			return true, resultNoSwitchTag
		}

		t := pass.TypesInfo.Types[sw.Tag]
		if !t.IsValue() {
			return true, resultTagNotValue
		}

		tagType, ok := t.Type.(*types.Named)
		if !ok {
			return true, resultTagNotNamed
		}

		tagPkg := tagType.Obj().Pkg()
		if tagPkg == nil {
			// The Go documentation says: nil for labels and objects in the Universe scope.
			// This happens for the `error` type, for example.
			return true, resultTagNoPkg
		}

		enumTyp := enumType{tagType}

		members, ok := importFact(pass, enumTyp)
		if !ok {
			// switch tag's type is not a known enum type.
			return true, resultTagNotEnum
		}

		samePkg := tagPkg == pass.Pkg // do the switch statement and the switch tag type (i.e. enum type) live in the same package?
		checkUnexported := samePkg    // we want to include unexported members in the exhaustiveness check only if we're in the same package
		checklist := makeChecklist(members, tagPkg, checkUnexported, cfg.ignoreEnumMembers)

		hasDefaultCase := analyzeSwitchClauses(sw, pass.TypesInfo, samePkg, func(memberName string) {
			checklist.found(memberName)
		})

		if len(checklist.remaining()) == 0 {
			// All enum members accounted for.
			// Nothing to report.
			return true, resultEnumMembersAccounted
		}
		if hasDefaultCase && cfg.defaultSignifiesExhaustive {
			// Though enum members are not accounted for,
			// the existence of the default case signifies exhaustiveness.
			// So don't report.
			return true, resultDefaultCaseSuffices
		}
		pass.Report(makeDiagnostic(sw, samePkg, enumTyp, members, toSlice(checklist.remaining())))
		return true, resultReportedDiagnostic
	}
}

// config is configuration for checkSwitchStatements.
type config struct {
	defaultSignifiesExhaustive bool
	checkGeneratedFiles        bool
	ignoreEnumMembers          *regexp.Regexp
}

// checkSwitchStatements checks exhaustiveness of enum switch statements for the supplied
// pass. It reports switch statements that are not exhaustive via pass.Report.
func checkSwitchStatements(pass *analysis.Pass, inspect *inspector.Inspector, cfg config) {
	f := switchStmtChecker(pass, cfg)

	inspect.WithStack([]ast.Node{&ast.SwitchStmt{}}, func(n ast.Node, push bool, stack []ast.Node) bool {
		proceed, _ := f(n, push, stack)
		return proceed
	})
}

func isDefaultCase(c *ast.CaseClause) bool {
	return c.List == nil // see doc comment on List field
}

// isPackageNameIdent returns whether ident represents an imported Go package.
func isPackageNameIdent(ident *ast.Ident, info *types.Info) bool {
	obj := info.ObjectOf(ident)
	if obj == nil {
		return false
	}
	_, ok := obj.(*types.PkgName)
	return ok
}

// analyzeSwitchClauses analyzes the clauses in the supplied switch statement.
//
// The info param should typically be pass.TypesInfo. The samePkg param
// indicates whether the switch tag type and the switch statement live in the
// same package. The found function is called for each enum member name found in
// the switch statement.
//
// The hasDefaultCase return value indicates whether the switch statement has a
// default clause.
func analyzeSwitchClauses(sw *ast.SwitchStmt, info *types.Info, samePkg bool, found func(identName string)) (hasDefaultCase bool) {
	for _, stmt := range sw.Body.List {
		caseCl := stmt.(*ast.CaseClause)
		if isDefaultCase(caseCl) {
			hasDefaultCase = true
			continue // nothing more to do if it's the default case
		}
		for _, expr := range caseCl.List {
			analyzeCaseClauseExpr(expr, info, samePkg, found)
		}
	}
	return hasDefaultCase
}

// Helper for analyzeSwitchClauses. See docs there.
func analyzeCaseClauseExpr(e ast.Expr, info *types.Info, samePkg bool, found func(identName string)) {
	e = astutil.Unparen(e)

	if samePkg {
		ident, ok := e.(*ast.Ident)
		if !ok {
			return
		}
		found(ident.Name)
		return
	}

	selExpr, ok := e.(*ast.SelectorExpr)
	if !ok {
		return
	}

	// Check that X (which is everything except the rightmost *ast.Ident, or
	// the Sel) is also an *ast.Ident, and particularly that it is a package
	// name identifier.
	x := astutil.Unparen(selExpr.X)
	ident, ok := x.(*ast.Ident)
	if !ok {
		return
	}

	if !isPackageNameIdent(ident, info) {
		return
	}

	// TODO: ident represents a package at this point; check if it represents
	// the enum package? (Is this additional check necessary? Wouldn't the type
	// checker have already failed if this wasn't the case?)
	// This may need additional thought for type aliases, too.

	found(selExpr.Sel.Name)
}

// diagnosticMissingMembers constructs the list of missing enum members,
// suitable for use in a reported diagnostic message.
func diagnosticMissingMembers(missingMembers []string, em *enumMembers) []string {
	missingByConstVal := make(map[constantValue][]string) // missing members, keyed by constant value.
	for _, m := range missingMembers {
		val := em.NameToValue[m]
		missingByConstVal[val] = append(missingByConstVal[val], m)
	}

	var out []string
	for _, names := range missingByConstVal {
		sort.Strings(names)
		out = append(out, strings.Join(names, "|"))
	}
	sort.Strings(out)
	return out
}

// diagnosticEnumTypeName returns a string representation of an enum type for
// use in reported diagnostics.
func diagnosticEnumTypeName(enumType *types.Named, samePkg bool) string {
	if samePkg {
		return enumType.Obj().Name()
	}
	return enumType.Obj().Pkg().Name() + "." + enumType.Obj().Name()
}

func makeDiagnostic(sw *ast.SwitchStmt, samePkg bool, enumTyp enumType, allMembers *enumMembers, missingMembers []string) analysis.Diagnostic {
	message := fmt.Sprintf("missing cases in switch of type %s: %s",
		diagnosticEnumTypeName(enumTyp.named, samePkg),
		strings.Join(diagnosticMissingMembers(missingMembers, allMembers), ", "))

	return analysis.Diagnostic{
		Pos:     sw.Pos(),
		End:     sw.End(),
		Message: message,
	}
}

func toSlice(m map[string]struct{}) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	return out
}
