package exhaustive

import (
	"fmt"
	"go/ast"
	"go/types"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
)

// nodeVisitor is like the visitor function used by Inspector.WithStack,
// except that it returns an additional value: a short description of
// the result of this node visit.
//
// The result is typically useful in debugging or in unit tests to check
// that the nodeVisitor function took the expected code path.
type nodeVisitor func(n ast.Node, push bool, stack []ast.Node) (proceed bool, result string)

// Result values returned by a node visitors.
const (
	resultEmptyMapLiteral = "empty map literal"
	resultNotMapLiteral   = "not map literal"
	resultKeyNotNamed     = "map key not named type"
	resultKeyNilPkg       = "nil map key package"
	resultKeyNotEnum      = "map key not known enum type"

	resultNoSwitchTag = "no switch tag"
	resultTagNotValue = "switch tag not value type"
	resultTagNotNamed = "switch tag not named type"
	resultTagNilPkg   = "nil switch tag package"
	resultTagNotEnum  = "switch tag not known enum type"

	resultNotPush              = "not push"
	resultGeneratedFile        = "generated file"
	resultIgnoreComment        = "has ignore comment"
	resultNoEnforceComment     = "has no enforce comment"
	resultEnumMembersAccounted = "required enum members accounted for"
	resultDefaultCaseSuffices  = "default case presence satisfies exhaustiveness"
	resultReportedDiagnostic   = "reported diagnostic"
)

// switchChecker returns a node visitor that checks exhaustiveness of
// enum switch statements for the supplied pass, and reports
// diagnostics. The node visitor expects only *ast.SwitchStmt nodes.
func switchChecker(pass *analysis.Pass, cfg switchConfig, generated boolCache, comments commentCache) nodeVisitor {
	return func(n ast.Node, push bool, stack []ast.Node) (bool, string) {
		if !push {
			// The proceed return value should not matter; it is ignored by
			// inspector package for pop calls.
			// Nevertheless, return true to be on the safe side for the future.
			return true, resultNotPush
		}

		file := stack[0].(*ast.File)

		if !cfg.checkGenerated && generated.get(file) {
			// Don't check this file.
			// Return false because the children nodes of node `n` don't have to be checked.
			return false, resultGeneratedFile
		}

		sw := n.(*ast.SwitchStmt)

		switchComments := comments.get(pass.Fset, file)[sw]
		if !cfg.explicit && hasComment(switchComments, ignoreComment) {
			// Skip checking of this switch statement due to ignore directive comment.
			// Still return true because there may be nested switch statements
			// that are not to be ignored.
			return true, resultIgnoreComment
		}
		if cfg.explicit && !hasComment(switchComments, enforceComment) {
			// Skip checking of this switch statement due to missing enforce directive comment.
			return true, resultNoEnforceComment
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
			return true, resultTagNilPkg
		}

		enumTyp := enumType{tagType.Obj()}
		members, ok := importFact(pass, enumTyp)
		if !ok {
			// switch tag's type is not a known enum type.
			return true, resultTagNotEnum
		}

		samePkg := tagPkg == pass.Pkg // do the switch statement and the switch tag type (i.e. enum type) live in the same package?
		checkUnexported := samePkg    // we want to include unexported members in the exhaustiveness check only if we're in the same package
		checklist := makeChecklist(members, tagPkg, checkUnexported, cfg.ignoreEnumMembers)

		hasDefaultCase := analyzeSwitchClauses(sw, pass.TypesInfo, checklist.found)

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
		pass.Report(makeSwitchDiagnostic(sw, samePkg, enumTyp, members, checklist.remaining()))
		return true, resultReportedDiagnostic
	}
}

// switchConfig is configuration for switchChecker.
type switchConfig struct {
	explicit                   bool
	defaultSignifiesExhaustive bool
	checkGenerated             bool
	ignoreEnumMembers          *regexp.Regexp // can be nil
}

func isDefaultCase(c *ast.CaseClause) bool {
	return c.List == nil // see doc comment on List field
}

func denotesPackage(ident *ast.Ident, info *types.Info) (*types.Package, bool) {
	obj := info.ObjectOf(ident)
	if obj == nil {
		return nil, false
	}
	n, ok := obj.(*types.PkgName)
	if !ok {
		return nil, false
	}
	return n.Imported(), true
}

// analyzeSwitchClauses analyzes the clauses in the supplied switch
// statement. The info param typically is pass.TypesInfo. The found
// function is called for each enum member name found in the switch
// statement. The hasDefaultCase return value indicates whether the
// switch statement has a default clause.
func analyzeSwitchClauses(sw *ast.SwitchStmt, info *types.Info, found func(val constantValue)) (hasDefaultCase bool) {
	for _, stmt := range sw.Body.List {
		caseCl := stmt.(*ast.CaseClause)
		if isDefaultCase(caseCl) {
			hasDefaultCase = true
			continue
		}
		for _, expr := range caseCl.List {
			analyzeCaseClauseExpr(expr, info, found)
		}
	}
	return hasDefaultCase
}

func analyzeCaseClauseExpr(e ast.Expr, info *types.Info, found func(val constantValue)) {
	handleIdent := func(ident *ast.Ident) {
		obj := info.Uses[ident]
		if obj == nil {
			return
		}
		if _, ok := obj.(*types.Const); !ok {
			return
		}

		// There are two scenarios.
		// See related test cases in typealias/quux/quux.go.
		//
		// ## Scenario 1
		//
		// Tag package and constant package are the same. This is
		// simple; we just use fs.ModeDir's value.
		//
		// Example:
		//
		//   var mode fs.FileMode
		//   switch mode {
		//   case fs.ModeDir:
		//   }
		//
		// ## Scenario 2
		//
		// Tag package and constant package are different. In this
		// scenario, too, we accept the case clause expr constant value,
		// as is. If the Go type checker is okay with the name being
		// listed in the case clause, we don't care much further.
		//
		// Example:
		//
		//   var mode fs.FileMode
		//   switch mode {
		//   case os.ModeDir:
		//   }
		//
		// Or equivalently:
		//
		//   // The type of mode is effectively fs.FileMode,
		//   // due to type alias.
		//   var mode os.FileMode
		//   switch mode {
		//   case os.ModeDir:
		//   }
		found(determineConstVal(ident, info))
	}

	e = astutil.Unparen(e)
	switch e := e.(type) {
	case *ast.Ident:
		handleIdent(e)

	case *ast.SelectorExpr:
		x := astutil.Unparen(e.X)
		// Ensure we only see the form pkg.Const, and not e.g.
		// structVal.f or structVal.inner.f.
		//
		// For this purpose, first we check that X, which is everything
		// except the rightmost field selector *ast.Ident (the Sel
		// field), is also an *ast.Ident.
		xIdent, ok := x.(*ast.Ident)
		if !ok {
			return
		}
		// Second , check that it's a package. It doesn't matter which
		// package, just that it denotes some package.
		if _, ok := denotesPackage(xIdent, info); !ok {
			return
		}
		handleIdent(e.Sel)
	}
}

// Makes a diagnostic for a non-exhaustive switch statement. samePkg
// should be true if the enum type and the switch statement are defined
// in the same package.
func makeSwitchDiagnostic(sw *ast.SwitchStmt, samePkg bool, enumTyp enumType, all enumMembers, missing map[string]struct{}) analysis.Diagnostic {
	typeName := diagnosticEnumTypeName(enumTyp.TypeName, samePkg)
	members := strings.Join(diagnosticMissingMembers(missing, all), ", ")

	return analysis.Diagnostic{
		Pos:     sw.Pos(),
		End:     sw.End(),
		Message: fmt.Sprintf("missing cases in switch of type %s: %s", typeName, members),
	}
}
