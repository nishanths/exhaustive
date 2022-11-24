package exhaustive

import (
	"fmt"
	"go/ast"
	"go/types"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// mapConfig is configuration for mapChecker.
type mapConfig struct {
	explicit          bool
	checkGenerated    bool
	ignoreEnumMembers *regexp.Regexp // can be nil
}

// mapChecker returns a node visitor that checks for exhaustiveness of
// map literals for the supplied pass, and reports diagnostics. The
// node visitor expects only *ast.CompositeLit nodes.
func mapChecker(pass *analysis.Pass, cfg mapConfig, generated boolCache, comments commentCache) nodeVisitor {
	return func(n ast.Node, push bool, stack []ast.Node) (bool, string) {
		if !push {
			return true, resultNotPush
		}

		file := stack[0].(*ast.File)

		if !cfg.checkGenerated && generated.get(file) {
			return false, resultGeneratedFile
		}

		lit := n.(*ast.CompositeLit)

		mapType, ok := pass.TypesInfo.Types[lit.Type].Type.(*types.Map)
		if !ok {
			namedType, ok2 := pass.TypesInfo.Types[lit.Type].Type.(*types.Named)
			if !ok2 {
				return true, resultNotMapLiteral
			}
			mapType, ok = namedType.Underlying().(*types.Map)
			if !ok {
				return true, resultNotMapLiteral
			}
		}

		if len(lit.Elts) == 0 {
			return false, resultEmptyMapLiteral
		}

		keyType, ok := mapType.Key().(*types.Named)
		if !ok {
			return true, resultKeyNotNamed
		}

		fileComments := comments.get(pass.Fset, file)
		var relatedComments []*ast.CommentGroup
		for i := range stack {
			// iterate over stack in the reverse order (from bottom to top)
			node := stack[len(stack)-1-i]
			switch node.(type) {
			// need to check comments associated with following nodes,
			// because logic of ast package doesn't allow to associate comment with *ast.CompositeLit
			case *ast.CompositeLit, // stack[len(stack)-1]
				*ast.ReturnStmt, // return ...
				*ast.IndexExpr,  // map[enum]...{...}[key]
				*ast.CallExpr,   // myfunc(map...)
				*ast.UnaryExpr,  // &map...
				*ast.AssignStmt, // variable assignment (without var keyword)
				*ast.DeclStmt,   // var declaration, parent of *ast.GenDecl
				*ast.GenDecl,    // var declaration, parent of *ast.ValueSpec
				*ast.ValueSpec:  // var declaration
				relatedComments = append(relatedComments, fileComments[node]...)
				continue
			}
			// stop iteration on the first inappropriate node
			break
		}

		if !cfg.explicit && hasComment(relatedComments, ignoreComment) {
			// Skip checking of this map literal due to ignore
			// comment. Still return true because there may be nested
			// map literals that are not to be ignored.
			return true, resultIgnoreComment
		}
		if cfg.explicit && !hasComment(relatedComments, enforceComment) {
			return true, resultNoEnforceComment
		}

		keyPkg := keyType.Obj().Pkg()
		if keyPkg == nil {
			return true, resultKeyNilPkg
		}

		enumTyp := enumType{keyType.Obj()}
		members, ok := importFact(pass, enumTyp)
		if !ok {
			return true, resultKeyNotEnum
		}

		samePkg := keyPkg == pass.Pkg // do the map literal and the map key type (i.e. enum type) live in the same package?
		checkUnexported := samePkg    // we want to include unexported members in the exhaustiveness check only if we're in the same package
		checklist := makeChecklist(members, keyPkg, checkUnexported, cfg.ignoreEnumMembers)

		for _, e := range lit.Elts {
			expr, ok := e.(*ast.KeyValueExpr)
			if !ok {
				continue // is it possible for valid map literal?
			}
			analyzeCaseClauseExpr(expr.Key, pass.TypesInfo, checklist.found)
		}

		if len(checklist.remaining()) == 0 {
			return true, resultEnumMembersAccounted
		}

		pass.Report(makeMapDiagnostic(lit, samePkg, enumTyp, members, checklist.remaining()))
		return true, resultReportedDiagnostic
	}
}

func makeMapDiagnostic(lit *ast.CompositeLit, samePkg bool, enumTyp enumType, all enumMembers, missing map[string]struct{}) analysis.Diagnostic {
	typeName := diagnosticEnumTypeName(enumTyp.TypeName, samePkg)
	members := strings.Join(diagnosticMissingMembers(missing, all), ", ")
	return analysis.Diagnostic{
		Pos:     lit.Pos(),
		End:     lit.End(),
		Message: fmt.Sprintf("missing keys in map of key type %s: %s", typeName, members),
	}
}
