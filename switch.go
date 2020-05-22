package exhaustive

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
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
				e = astutil.Unparen(e)
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
		if fix, ok := computeFix(pass, f, sw, enumType, samePkg, missingMembers); ok {
			fixes = append(fixes, fix)
		}
	}

	pass.Report(analysis.Diagnostic{
		Pos:            sw.Pos(),
		End:            sw.End(),
		Message:        fmt.Sprintf("missing cases in switch of type %s: %s", enumTypeName(enumType, samePkg), strings.Join(missing, ", ")),
		SuggestedFixes: fixes,
	})
}

func computeFix(pass *analysis.Pass, f *ast.File, sw *ast.SwitchStmt, enumType *types.Named, samePkg bool, missingMembers map[string]struct{}) (analysis.SuggestedFix, bool) {
	// Calls may be mutative, so we don't want to reuse the call expression in the
	// about-to-be-inserted case clause body.
	//
	// So we just don't suggest a fix in such situations.
	if _, ok := sw.Tag.(*ast.CallExpr); ok {
		return analysis.SuggestedFix{}, false
	}

	// Construct insertion text for case clause and its body.

	var tag bytes.Buffer
	printer.Fprint(&tag, pass.Fset, sw.Tag)

	// If possible, determine the package identifier based on the AST of other case clauses.
	var pkgIdent *ast.Ident
	if !samePkg {
		for _, caseCl := range sw.Body.List {
			caseCl = caseCl.(*ast.CaseClause)
			// At least one expression must exist in List at this point.
			// List cannot be nil because we only arrive here if the "default" clause
			// does not exist. Additionally, a syntactically valid case clause must
			// have at least one expression.
			if sel, ok := caseCl.List[0].(*ast.SelectorExpr); ok {
				pkgIdent = sel.X.(*ast.Ident)
				break
			}
		}
	}

	missing := make([]string, 0, len(missingMembers))
	for m := range missingMembers {
		if !samePkg {
			if pkgIdent != nil {
				// we were able to determine package identifier
				missing = append(missing, pkgIdent.Name+"."+m)
			} else {
				// use the package name (may not be correct always)
				//
				// TODO: May need to also add import if the package isn't imported
				// elsewhere. This (ie, a switch with zero case clauses) should
				// happen rarely, so don't implement this for now.
				missing = append(missing, enumType.Obj().Pkg().Name()+"."+m)
			}
		} else {
			missing = append(missing, m)
		}
	}
	sort.Strings(missing)

	insert := `case ` + strings.Join(missing, ", ") + `:
	panic(fmt.Sprintf("unhandled value: %v",` + tag.String() + `))`

	// TODO may need to add "fmt" import
	//
	// if fmt exists in any one of import GenDecl nothing to do.
	//
	// else: get the first import GenDecl.
	// determine if it has parens.
	// print it out.
	// if it has no parens, convert the printed form to have parens.
	// find the insertion position and insert "fmt".
	// done.

	pos := sw.Body.Lbrace + 1
	if len(sw.Body.List) != 0 {
		pos = sw.Body.List[len(sw.Body.List)-1].End()
	}
	textEdit := analysis.TextEdit{
		Pos:     pos,
		End:     pos,
		NewText: []byte(insert),
	}

	return analysis.SuggestedFix{
		Message:   fmt.Sprintf("add case clause for: %s?", strings.Join(missing, ", ")),
		TextEdits: []analysis.TextEdit{textEdit},
	}, true
}
