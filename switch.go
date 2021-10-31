package exhaustive

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/ast/inspector"
)

func isDefaultCase(c *ast.CaseClause) bool {
	return c.List == nil // see doc comment on field
}

func isPackageNameIdentifier(typesInfo *types.Info, ident *ast.Ident) bool {
	obj := typesInfo.ObjectOf(ident)
	if obj == nil {
		return false
	}
	_, ok := obj.(*types.PkgName)
	return ok
}

func enumTypeName(e *types.Named, samePkg bool) string {
	if samePkg {
		return e.Obj().Name()
	}
	return e.Obj().Pkg().Name() + "." + e.Obj().Name()
}

func analyzeSwitchClauses(sw *ast.SwitchStmt, typesInfo *types.Info, hitlist *hitlist, strategy hitlistStrategy, samePkg bool) (defaultCase *ast.CaseClause) {
	for _, stmt := range sw.Body.List {
		caseCl := stmt.(*ast.CaseClause)

		if isDefaultCase(caseCl) {
			defaultCase = caseCl
			continue // nothing more to do if it's the default case
		}

		analyzeCaseClause(caseCl, typesInfo, hitlist, strategy, samePkg)
	}

	return
}

func analyzeCaseClause(caseCl *ast.CaseClause, typesInfo *types.Info, hitlist *hitlist, strategy hitlistStrategy, samePkg bool) {
	for _, e := range caseCl.List {
		analyzeCaseClauseExpr(e, typesInfo, hitlist, strategy, samePkg)
	}
}

func analyzeCaseClauseExpr(e ast.Expr, typesInfo *types.Info, hitlist *hitlist, strategy hitlistStrategy, samePkg bool) {
	e = astutil.Unparen(e)

	if samePkg {
		ident, ok := e.(*ast.Ident)
		if !ok {
			return
		}
		hitlist.found(ident.Name, strategy)
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
	if !isPackageNameIdentifier(typesInfo, ident) {
		return
	}

	hitlist.found(selExpr.Sel.Name, strategy)
}

func checkSwitchStatements(pass *analysis.Pass, inspect *inspector.Inspector, strategy hitlistStrategy) error {
	return checkSwitchStatements_(pass, inspect, strategy, make(map[*ast.File]ast.CommentMap), make(map[*ast.File]bool))
}

func checkSwitchStatements_(pass *analysis.Pass, inspect *inspector.Inspector, strategy hitlistStrategy, comments map[*ast.File]ast.CommentMap, generated map[*ast.File]bool) error {
	inspect.WithStack([]ast.Node{&ast.SwitchStmt{}}, func(n ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			// we only inspect things on the way down, not up.
			return true
		}

		file := stack[0].(*ast.File)

		// Determine if the file is a generated file, and save the result.
		// If it is a generated file, don't check the file.
		if _, ok := generated[file]; !ok {
			generated[file] = isGeneratedFile(file)
		}
		if generated[file] && !fCheckGeneratedFiles {
			// don't check this file.
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
			// The Go documentation says: nil for labels and objects in the Universe scope.
			// This happens for the `error` type, for example.
			// Continuing would mean that pass.ImportPackageFact panics.
			return true
		}

		var enums enumsFact
		if !pass.ImportPackageFact(tagPkg, &enums) {
			panic(fmt.Sprintf("pass.ImportPackageFact returned false for %s", tagPkg))
		}

		em, isEnum := enums.Enums[tagType.Obj().Name()]
		if !isEnum {
			// switch tag's type is not a known enum type.
			return true
		}

		if _, ok := comments[file]; !ok {
			comments[file] = ast.NewCommentMap(pass.Fset, file, file.Comments)
		}
		if containsIgnoreDirectiveGroups(comments[file].Filter(sw).Comments()) {
			return true // skip checking due to ignore directive
		}

		samePkg := tagPkg == pass.Pkg
		checkUnexported := samePkg

		hitlist := makeHitlist(em, tagPkg, checkUnexported, fIgnorePattern.Get().(*regexp.Regexp))

		defaultCase := analyzeSwitchClauses(sw, pass.TypesInfo, hitlist, strategy, samePkg)

		defaultSuffices := fDefaultSignifiesExhaustive && defaultCase != nil
		shouldReport := len(hitlist.remaining()) != 0 && !defaultSuffices

		if shouldReport {
			reportSwitch(pass, sw, defaultCase, samePkg, tagType, em, hitlist.remaining(), file)
		}
		return true
	})

	return nil
}

func missingCasesOutput(missingMembers map[string]struct{}, em *enumMembers) []string {
	constValMembers := make(map[string][]string) // constant value -> member name
	var otherMembers []string                    // non-constant value member names

	for m := range missingMembers {
		if constVal, ok := em.NameToValue[m]; ok {
			constValMembers[constVal] = append(constValMembers[constVal], m)
		} else {
			otherMembers = append(otherMembers, m)
		}
	}

	ret := make([]string, 0, len(constValMembers)+len(otherMembers))
	for _, names := range constValMembers {
		sort.Strings(names)
		ret = append(ret, strings.Join(names, "|"))
	}
	ret = append(ret, otherMembers...)
	sort.Strings(ret)
	return ret
}

func reportSwitch(
	pass *analysis.Pass,
	sw *ast.SwitchStmt,
	defaultCase *ast.CaseClause,
	samePkg bool,
	enumType *types.Named,
	em *enumMembers,
	missingMembers map[string]struct{},
	f *ast.File,
) {
	missingOutput := missingCasesOutput(missingMembers, em)

	var fixes []analysis.SuggestedFix
	if fix, ok := computeFix(pass, pass.Fset, f, sw, defaultCase, enumType, samePkg, missingMembers); ok {
		fixes = append(fixes, fix)
	}

	pass.Report(analysis.Diagnostic{
		Pos:            sw.Pos(),
		End:            sw.End(),
		Message:        fmt.Sprintf("missing cases in switch of type %s: %s", enumTypeName(enumType, samePkg), strings.Join(missingOutput, ", ")),
		SuggestedFixes: fixes,
	})
}

// --- move the below to fix.go

func computeFix(pass *analysis.Pass, fset *token.FileSet, f *ast.File, sw *ast.SwitchStmt, defaultCase *ast.CaseClause, enumType *types.Named, samePkg bool, missingMembers map[string]struct{}) (analysis.SuggestedFix, bool) {
	// Function and method calls may be mutative, so we don't want to reuse the
	// call expression in the about-to-be-inserted case clause body. So we just
	// don't suggest a fix in such situations.
	//
	// However, we need to make an exception for type conversions, which are
	// also call expressions in the AST.
	//
	// We'll need to lookup type information for this, and can't rely solely
	// on the AST.
	if containsFuncCall(pass.TypesInfo, sw.Tag) {
		return analysis.SuggestedFix{}, false
	}

	textEdits := []analysis.TextEdit{missingCasesTextEdit(fset, f, samePkg, sw, defaultCase, enumType, missingMembers)}

	// need to add "fmt" import if "fmt" import doesn't already exist
	if !hasImportWithPath(flattenImportSpec(astutil.Imports(fset, f)), `"fmt"`) {
		textEdits = append(textEdits, fmtImportTextEdit(fset, f))
	}

	missing := make([]string, 0, len(missingMembers))
	for m := range missingMembers {
		missing = append(missing, m)
	}
	sort.Strings(missing)

	return analysis.SuggestedFix{
		Message:   fmt.Sprintf("add case clause for: %s", strings.Join(missing, ", ")),
		TextEdits: textEdits,
	}, true
}

func missingCasesTextEdit(fset *token.FileSet, f *ast.File, samePkg bool, sw *ast.SwitchStmt, defaultCase *ast.CaseClause, enumType *types.Named, missingMembers map[string]struct{}) analysis.TextEdit {
	// ... Construct insertion text for case clause and its body ...

	var tag bytes.Buffer
	printer.Fprint(&tag, fset, sw.Tag)

	// If possible and if necessary, determine the package identifier based on
	// the AST of other `case` clauses.
	var pkgIdent *ast.Ident
	if !samePkg {
		for _, stmt := range sw.Body.List {
			caseCl := stmt.(*ast.CaseClause)
			if len(caseCl.List) != 0 { // guard against default case
				if sel, ok := caseCl.List[0].(*ast.SelectorExpr); ok {
					pkgIdent = sel.X.(*ast.Ident)
					break
				}
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

	// ... Create the text edit ...

	pos := sw.Body.Rbrace - 1 // put it as last case
	if defaultCase != nil {
		pos = defaultCase.Case - 2 // put it before the default case (why -2?)
	}

	return analysis.TextEdit{
		Pos:     pos,
		End:     pos,
		NewText: []byte(insert),
	}
}
