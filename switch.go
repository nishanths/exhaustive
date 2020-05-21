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
	return c.List == nil // see doc comment on field
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
		tagType, ok := t.Type.(*types.Named)
		if !ok {
			return false
		}

		tagPkg := tagType.Obj().Pkg()
		if tagPkg == nil {
			// Doc comment: nil for labels and objects in the Universe scope.
			// This happens for the `error` type, for example.
			// Continuing would mean that ImportPackageFact panics.
			return false
		}

		var enums enumsFact
		if !pass.ImportPackageFact(tagPkg, &enums) {
			// Can't do anything further.
			return false
		}

		enumMembers, isEnum := enums.entries[tagType]
		if !isEnum {
			// Tag's type is not a known enum.
			return false
		}

		// Get comment map.
		file := stack[0].(*ast.File)
		var allComments ast.CommentMap
		if cm, ok := comments[file]; ok {
			allComments = cm
		} else {
			allComments = ast.NewCommentMap(pass.Fset, file, file.Comments)
			comments[file] = allComments
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

		if sw.Body == nil {
			// TODO: Is this even syntactically valid?
			//
			// Either way, nothing is deleted from hitlist in this case (all
			// members are reported as missing).
			reportSwitch(pass, sw, samePkg, tagType, hitlist)
			return false
		}

		for _, stmt := range sw.Body.List {
			caseCl := stmt.(*ast.CaseClause)
			if isDefaultCase(caseCl) {
				if fDefaultSuffices {
					return false
				}
				continue // nothing more to do if it's the default case
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

		reportSwitch(pass, sw, samePkg, tagType, hitlist)
		return false
	})
}

func reportSwitch(pass *analysis.Pass, sw *ast.SwitchStmt, samePkg bool, enumType *types.Named, missingMembers map[string]struct{}) {
	missing := make([]string, 0, len(missingMembers))
	for m := range missingMembers {
		missing = append(missing, m)
	}
	sort.Strings(missing)

	pass.ReportRangef(sw, "missing cases in switch of type %s: %s",
		enumTypeName(enumType, samePkg), strings.Join(missing, ", "))
}

func removeParens(e ast.Expr) ast.Expr {
	for {
		parenExpr, ok := e.(*ast.ParenExpr)
		if !ok {
			break
		}
		e = parenExpr.X
	}
	return e
}
