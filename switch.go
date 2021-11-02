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

type checkingStrategy int

const (
	strategyValue checkingStrategy = iota
	strategyName
)

// nodeVisitor is similar to the visitor function used by Inspector.WithStack,
// except that it returns two additional values: a short description of
// the result of this node visit, and an error.
//
// The result is typically useful in debugging or in unit tests to check
// that the nodeVisitor function took the expected code path.
//
// A returned non-nil error does not stop further calls to the visitor; this is
// solely controlled by the proceed value. The error however allows callers
// to e.g. record errors encountered during visits.
type nodeVisitor func(n ast.Node, push bool, stack []ast.Node) (proceed bool, result string, err error)

// Result values returned by a node visitor constructed via switchStmtChecker.
const (
	resultNotPush              = "not push"
	resultGeneratedFile        = "generated file"
	resultNoSwitchTag          = "no switch tag"
	resultTagNotValue          = "switch tag not value type"
	resultTagNotNamed          = "switch tag not named type"
	resultTagNoPkg             = "switch tag does not belong to regular package"
	resultPassImportFailed     = "pass.ImportPackageFact failed"
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

	return func(n ast.Node, push bool, stack []ast.Node) (bool, string, error) {
		if !push {
			// we only inspect things on the way down, not up.
			return true, resultNotPush, nil
		}

		file := stack[0].(*ast.File)

		// Determine if the file is a generated file, and save the result.
		// If it is a generated file, don't check the file.
		if _, ok := generated[file]; !ok {
			generated[file] = isGeneratedFile(file)
		}
		if generated[file] && !cfg.checkGeneratedFiles {
			// don't check this file.
			return true, resultGeneratedFile, nil
		}

		sw := n.(*ast.SwitchStmt)

		if sw.Tag == nil {
			return true, resultNoSwitchTag, nil
		}

		t := pass.TypesInfo.Types[sw.Tag]
		if !t.IsValue() {
			return true, resultTagNotValue, nil
		}

		tagType, ok := t.Type.(*types.Named)
		if !ok {
			return true, resultTagNotNamed, nil
		}

		tagPkg := tagType.Obj().Pkg()
		if tagPkg == nil {
			// The Go documentation says: nil for labels and objects in the Universe scope.
			// This happens for the `error` type, for example.
			// Continuing would mean that pass.ImportPackageFact panics.
			return true, resultTagNoPkg, nil
		}

		var enums enumsFact
		if !pass.ImportPackageFact(tagPkg, &enums) {
			return true, resultPassImportFailed, fmt.Errorf("pass.ImportPackageFact returned false for %s", tagPkg)
		}

		em, isEnum := enums.Enums[tagType.Obj().Name()]
		if !isEnum {
			// switch tag's type is not a known enum type.
			return true, resultTagNotEnum, nil
		}

		if _, ok := comments[file]; !ok {
			comments[file] = ast.NewCommentMap(pass.Fset, file, file.Comments)
		}
		if containsIgnoreDirective(comments[file].Filter(sw).Comments()) {
			// skip checking due to ignore directive
			return true, resultSwitchIgnoreComment, nil
		}

		samePkg := tagPkg == pass.Pkg // do the switch statement and the switch tag type (i.e. enum type) live in the same package?
		checkUnexported := samePkg    // we want to include unexported members in the exhaustiveness check only if we're in the same package
		hitlist := makeHitlist(em, tagPkg, checkUnexported, cfg.ignoreEnumMembers)

		hasDefaultCase := analyzeSwitchClauses(sw, pass.TypesInfo, samePkg, func(memberName string) {
			hitlist.found(memberName, cfg.checkingStrategy)
		})

		if len(hitlist.remaining()) == 0 {
			// All enum members accounted for.
			// Nothing to report.
			return true, resultEnumMembersAccounted, nil
		}
		if hasDefaultCase && cfg.defaultSignifiesExhaustive {
			// Though enum members are not accounted for,
			// the existence of the default case signifies exhaustiveness.
			// So don't report.
			return true, resultDefaultCaseSuffices, nil
		}
		pass.Report(makeDiagnostic(sw, samePkg, tagType, em, toSlice(hitlist.remaining()), cfg.checkingStrategy))
		return true, resultReportedDiagnostic, nil
	}
}

// config is configuration for checkSwitchStatements.
type config struct {
	defaultSignifiesExhaustive bool
	checkGeneratedFiles        bool
	ignoreEnumMembers          *regexp.Regexp
	checkingStrategy           checkingStrategy
}

// checkSwitchStatements checks exhaustiveness of enum switch statements for the supplied
// pass. It reports switch statements that are not exhaustive via pass.Report.
func checkSwitchStatements(pass *analysis.Pass, inspect *inspector.Inspector, cfg config) error {
	f := switchStmtChecker(pass, cfg)

	var firstErr error
	setErr := func(err error) {
		if firstErr == nil {
			firstErr = err
		}
	}

	inspect.WithStack([]ast.Node{&ast.SwitchStmt{}}, func(n ast.Node, push bool, stack []ast.Node) bool {
		proceed, _, err := f(n, push, stack)
		if err != nil {
			setErr(err)
		}
		return proceed
	})

	return firstErr
}

func isDefaultCase(c *ast.CaseClause) bool {
	return c.List == nil // see doc comment on List field
}

// isPackageNameIdent returns whether ident represents an imported Go package.
func isPackageNameIdent(ident *ast.Ident, typesInfo *types.Info) bool {
	obj := typesInfo.ObjectOf(ident)
	if obj == nil {
		return false
	}
	_, ok := obj.(*types.PkgName)
	return ok
}

// analyzeSwitchClauses analyzes the clauses in the supplied switch statement.
//
// The typesInfo param should typically be pass.TypesInfo. The samePkg param
// indicates whether the switch tag type and the switch statement live in the
// same package. The found function is called for each enum member name found in
// the switch statement.
//
// The hasDefaultCase return value indicates whether the switch statement has a
// default clause.
func analyzeSwitchClauses(sw *ast.SwitchStmt, typesInfo *types.Info, samePkg bool, found func(identName string)) (hasDefaultCase bool) {
	for _, stmt := range sw.Body.List {
		caseCl := stmt.(*ast.CaseClause)
		if isDefaultCase(caseCl) {
			hasDefaultCase = true
			continue // nothing more to do if it's the default case
		}
		for _, expr := range caseCl.List {
			analyzeCaseClauseExpr(expr, typesInfo, samePkg, found)
		}
	}
	return hasDefaultCase
}

// Helper for analyzeSwitchClauses. See docs there.
func analyzeCaseClauseExpr(e ast.Expr, typesInfo *types.Info, samePkg bool, found func(identName string)) {
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

	if !isPackageNameIdent(ident, typesInfo) {
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
func diagnosticMissingMembers(missingMembers []string, em *enumMembers, strategy checkingStrategy) []string {
	switch strategy {
	case strategyValue:
		var out []string

		constValMembers := make(map[string][]string) // constant value -> member name
		var otherMembers []string                    // non-constant value member names

		for _, m := range missingMembers {
			if constVal, ok := em.NameToValue[m]; ok {
				constValMembers[constVal] = append(constValMembers[constVal], m)
			} else {
				otherMembers = append(otherMembers, m)
			}
		}

		for _, names := range constValMembers {
			sort.Strings(names)
			out = append(out, strings.Join(names, "|"))
		}
		out = append(out, otherMembers...)
		sort.Strings(out)
		return out

	case strategyName:
		out := make([]string, len(missingMembers))
		copy(out, missingMembers)
		sort.Strings(out)
		return out

	default:
		panic(fmt.Sprintf("unknown strategy %v", strategy))
	}
}

// diagnosticEnumTypeName returns a string representation of an enum type for
// use in reported diagnostics.
func diagnosticEnumTypeName(enumType *types.Named, samePkg bool) string {
	if samePkg {
		return enumType.Obj().Name()
	}
	return enumType.Obj().Pkg().Name() + "." + enumType.Obj().Name()
}

func makeDiagnostic(sw *ast.SwitchStmt, samePkg bool, enumType *types.Named, allMembers *enumMembers, missingMembers []string, strategy checkingStrategy) analysis.Diagnostic {
	message := fmt.Sprintf("missing cases in switch of type %s: %s",
		diagnosticEnumTypeName(enumType, samePkg),
		strings.Join(diagnosticMissingMembers(missingMembers, allMembers, strategy), ", "))

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
