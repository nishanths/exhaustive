package exhaustive

import (
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

func isDefaultCase(c *ast.CaseClause) bool {
	return c.List == nil
}

func checkSwitchStatements(pass *analysis.Pass, inspect *inspector.Inspector, comments map[*ast.File]ast.CommentMap) {
	inspect.WithStack([]ast.Node{&ast.SwitchStmt{}}, func(n ast.Node, _ bool, stack []ast.Node) bool {
		sw := n.(*ast.SwitchStmt)
		if sw.Tag == nil {
			return false
		}
		t := pass.TypesInfo.Types[sw.Tag]
		if !t.IsValue() {
			return false
		}
		named, ok := t.Type.(*types.Named)
		if !ok {
			return false
		}

		tagPkg := named.Obj().Pkg()
		if tagPkg == nil {
			// doc comment: nil for labels and objects in the Universe scope
			//
			// happens for the `error` type.
			// continuing would mean that ImportPackageFact panics.
			return false
		}

		if sw.Body == nil {
			return false
		}

		var enums enumsFact
		if !pass.ImportPackageFact(tagPkg, &enums) {
			// can't do anything further
			return false
		}

		enumMembers, isEnum := enums.entries[named]
		if !isEnum {
			// tag's type is not a known enum
			return false
		}

		// Get comment map.
		file := stack[0].(*ast.File)

		var allComments ast.CommentMap
		if cm, ok := comments[file]; ok {
			allComments = cm
		} else {
			allComments = ast.NewCommentMap(pass.Fset, file, file.Comments)
		}

		specificComments := allComments.Filter(sw)
		for _, group := range specificComments.Comments() {
			if containsIgnoreDirective(group.List) {
				return false // skip checking due to ignore directive
			}
		}

		samePkg := tagPkg == pass.Pkg
		checkUnexported := samePkg

		hitlist := make(map[string]struct{})
		for _, m := range enumMembers {
			if m.Exported() || checkUnexported {
				hitlist[m.Name()] = struct{}{}
			}
		}

		for _, stmt := range sw.Body.List {
			caseCl := stmt.(*ast.CaseClause)
			if isDefaultCase(caseCl) && fDefaultSuffices {
				return false
			}
			for _, e := range caseCl.List {
				e = removeParens(e)
				if samePkg {
					ident, ok := e.(*ast.Ident)
					if !ok {
						continue
					}
					delete(hitlist, ident.Name)
				} else {
					selExpr, ok := e.(*ast.SelectorExpr)
					if !ok {
						continue
					}
					delete(hitlist, selExpr.Sel.Name)
				}
			}
		}

		reportSwitch(pass, sw, samePkg, named, hitlist)
		return false
	})
}

func reportSwitch(pass *analysis.Pass, rng analysis.Range, samePkg bool, enumType *types.Named, missingMembers map[string]struct{}) {
	var enumTypeName string
	if samePkg {
		enumTypeName = enumType.Obj().Name()
	} else {
		enumTypeName = enumType.Obj().Pkg().Name() + "." + enumType.Obj().Name()
	}

	var missing []string
	for m := range missingMembers {
		missing = append(missing, m)
	}
	sort.Strings(missing)

	pass.ReportRangef(rng, "missing cases in switch of type %s: %s", enumTypeName, strings.Join(missing, ", "))
}
