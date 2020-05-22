package exhaustive

import (
	"fmt"
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
	inspect.WithStack([]ast.Node{&ast.SwitchStmt{}}, func(n ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}
		sw := n.(*ast.SwitchStmt)
		if sw.Tag == nil {
			return true
		}
		t := pass.TypesInfo.Types[sw.Tag]
		if !t.IsValue() {
			return true
		}
		tagType, ok := t.Type.(*types.Named)
		if !ok {
			return true
		}

		tagPkg := tagType.Obj().Pkg()
		if tagPkg == nil {
			// Doc comment: nil for labels and objects in the Universe scope.
			// This happens for the `error` type, for example.
			// Continuing would mean that ImportPackageFact panics.
			return true
		}

		var enums enumsFact
		if !pass.ImportPackageFact(tagPkg, &enums) {
			// Can't do anything further.
			return true
		}

		enumMembers, isEnum := enums.entries[tagType]
		if !isEnum {
			// Tag's type is not a known enum.
			return true
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
				return true // skip checking due to ignore directive
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

		if len(hitlist) == 0 {
			// can happen if external package and enum consists only of
			// unexported members
			return true
		}

		defaultCaseExists := false
		for _, stmt := range sw.Body.List {
			caseCl := stmt.(*ast.CaseClause)
			if isDefaultCase(caseCl) {
				defaultCaseExists = true
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

		defaultSuffices := fDefaultSuffices && defaultCaseExists
		shouldReport := len(hitlist) > 0 && !defaultSuffices

		if shouldReport {
			reportSwitch(pass, sw, samePkg, tagType, hitlist, defaultCaseExists, file)
		}
		return true
	})
}

func reportSwitch(pass *analysis.Pass, sw *ast.SwitchStmt, samePkg bool, enumType *types.Named, missingMembers map[string]struct{}, defaultCaseExists bool, f *ast.File) {
	missing := make([]string, 0, len(missingMembers))
	for m := range missingMembers {
		missing = append(missing, m)
	}
	sort.Strings(missing)

	var fixes []analysis.SuggestedFix
	if !defaultCaseExists {
		fixes = []analysis.SuggestedFix{
			computeFix(pass, f, sw, enumType, samePkg, missingMembers),
		}
	}

	pass.Report(analysis.Diagnostic{
		Pos:            sw.Pos(),
		End:            sw.End(),
		Message:        fmt.Sprintf("missing cases in switch of type %s: %s", enumTypeName(enumType, samePkg), strings.Join(missing, ", ")),
		SuggestedFixes: fixes,
	})
}

func computeFix(pass *analysis.Pass, f *ast.File, sw *ast.SwitchStmt, enumType *types.Named, samePkg bool, missingMembers map[string]struct{}) analysis.SuggestedFix {
	lastIdx := len(sw.Body.List) - 1
	lastCase := sw.Body.List[lastIdx]
	// log.Println(lastCase.End())
	// log.Println(sw.Body.Rbrace)

	missing := make([]string, 0, len(missingMembers))
	for m := range missingMembers {
		missing = append(missing, m) // TODO package name
	}
	sort.Strings(missing)

	textEdit := analysis.TextEdit{
		Pos:     lastCase.End(),
		NewText: []byte(`case ` + strings.Join(missing, ", ") + ":"),
	}
	return analysis.SuggestedFix{
		Message:   "apply?",
		TextEdits: []analysis.TextEdit{textEdit},
	}
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
