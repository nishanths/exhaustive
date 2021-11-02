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

func analyzeSwitchClauses(sw *ast.SwitchStmt, typesInfo *types.Info, samePkg bool, found func(identName string)) (defaultCase *ast.CaseClause) {
	for _, stmt := range sw.Body.List {
		caseCl := stmt.(*ast.CaseClause)
		if isDefaultCase(caseCl) {
			defaultCase = caseCl
			continue // nothing more to do if it's the default case
		}
		for _, e := range caseCl.List {
			analyzeCaseClauseExpr(e, typesInfo, samePkg, found)
		}
	}
	return defaultCase
}

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

	if !isPackageNameIdentifier(typesInfo, ident) {
		return
	}
	// TODO(next): possible to check if ident is the package name of the enum package name?

	found(selExpr.Sel.Name)
}

type config struct {
	defaultSignifiesExhaustive bool
	checkGeneratedFiles        bool
	ignoreMembers              *regexp.Regexp
	hitlistStrategy            hitlistStrategy
}

func checkSwitchStatements(pass *analysis.Pass, inspect *inspector.Inspector, cfg config) error {
	comments := make(map[*ast.File]ast.CommentMap)
	generated := make(map[*ast.File]bool)

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
		if generated[file] && !cfg.checkGeneratedFiles {
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

		hitlist := makeHitlist(em, tagPkg, checkUnexported, cfg.ignoreMembers)

		defaultCase := analyzeSwitchClauses(sw, pass.TypesInfo, samePkg, func(name string) {
			hitlist.found(name, cfg.hitlistStrategy)
		})

		defaultSuffices := cfg.defaultSignifiesExhaustive && defaultCase != nil
		shouldReport := len(hitlist.remaining()) != 0 && !defaultSuffices

		if shouldReport {
			report(pass, sw, samePkg, tagType, em, hitlist.remaining())
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

func report(
	pass *analysis.Pass,
	sw *ast.SwitchStmt,
	samePkg bool,
	enumType *types.Named,
	em *enumMembers,
	missingMembers map[string]struct{},
) {
	message := fmt.Sprintf("missing cases in switch of type %s: %s",
		enumTypeName(enumType, samePkg),
		strings.Join(missingCasesOutput(missingMembers, em), ", "))

	pass.Report(analysis.Diagnostic{
		Pos:            sw.Pos(),
		End:            sw.End(),
		Message:        message,
		SuggestedFixes: nil,
	})
}
