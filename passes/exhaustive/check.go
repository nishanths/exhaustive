package exhaustive

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/types"
	"maps"
	"slices"
	"strings"

	"github.com/nishanths/exhaustive/passes/enumerated"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func checkSwitch(pass *analysis.Pass, opts *options) {
	var (
		in             = pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		enums          = pass.ResultOf[enumerated.Analyzer].(enumerated.Result)
		commentsByFile = make(map[*ast.File]ast.CommentMap)
		// Helper functions.
		printExpr = func(e ast.Expr) string { return printNode(pass.Fset, ast.Unparen(e)) }
		errorType = func(e ast.Expr) string { return fmt.Sprintf("%s: could not determine type", printExpr(e)) }
	)

	in.Root().Inspect([]ast.Node{(*ast.SwitchStmt)(nil)}, func(c inspector.Cursor) (descend bool) {
		file := head(c.Enclosing((*ast.File)(nil))).Node().(*ast.File)
		if !opts.checkGenerated && ast.IsGenerated(file) {
			return false
		}

		sw := c.Node().(*ast.SwitchStmt)

		// Parse comment directives.
		fileComments, ok := commentsByFile[file]
		if !ok {
			fileComments = ast.NewCommentMap(pass.Fset, file, file.Comments)
			commentsByFile[file] = fileComments
		}
		directives, err := parseDirectives(fileComments[sw])
		if err != nil {
			pass.Reportf(sw.Pos(), "error parsing comment directives: %s", err)
			return true
		}
		if opts.needEnforceDirective && !directives[dirEnforce] {
			return true
		}
		if directives[dirIgnore] {
			return true
		}

		reasonUnchecked := ""
		if directives[dirEnforce] {
			defer func() {
				if reasonUnchecked != "" {
					pass.Reportf(sw.Pos(), "enforce directive present and switch statement not checked: %s", reasonUnchecked)
				}
			}()
		}

		// The tag will be nil for switch statements of
		// the form "switch { ... }".
		if sw.Tag == nil {
			reasonUnchecked = "missing switch expression"
			return true
		}

		tagType := pass.TypesInfo.TypeOf(sw.Tag)
		if tagType == nil {
			pass.Reportf(sw.Tag.Pos(), "%s", errorType(sw.Tag))
			return true
		}

		t, cs, ok := knownEnumeratedType(enums, tagType)
		if !ok {
			reasonUnchecked = "switch expression type is not an enumerated type"
			return true
		}

		// Handle user-specified include/exclude pattern
		// flags.
		//
		// Note that comment directives are more specific
		// (they are specified at source code level) than
		// include/exclude pattern flags. Hence an
		// enforce directive, if present, should take
		// precedence over any include/exclude patterns.
		if !directives[dirEnforce] {
			if !proceedInclExclPattern(t, opts) {
				return true
			}
		}

		needDefault := opts.defaultCaseRequired
		if v, ok := directives[dirDefrequire]; ok {
			needDefault = v
		}
		foundDefault := false

		unsatisfied := mustSatisfy(pass.Pkg, cs, opts)
		for _, cc := range sw.Body.List {
			cc := cc.(*ast.CaseClause)
			if cc.List == nil {
				foundDefault = true
				if opts.defaultCaseExhaustive {
					clear(unsatisfied)
					break
				}
			}
			for _, expr := range cc.List {
				tv, ok := pass.TypesInfo.Types[expr]
				if !ok {
					pass.Reportf(expr.Pos(), "%s", errorType(expr))
					continue
				}
				if tv.Value != nil {
					delete(unsatisfied, tv.Value.ExactString())
				}
			}
		}
		if needDefault && !foundDefault {
			pass.Reportf(sw.Pos(), "missing default case")
		}
		if len(unsatisfied) != 0 {
			pass.Reportf(sw.Pos(), "switch not exhaustive: missing cases: %s", formatUnsatisfiedNames(pass.Pkg, unsatisfied))
		}
		return true
	})
}

// knownEnumeratedType reports whether the given type (which
// typically is the type of an arbitrary expression) is an
// enumerated type. It returns the enumerated type and the
// enumerated constants of that type.
func knownEnumeratedType(enums enumerated.Result, t types.Type) (*types.Named, []*types.Const, bool) {
	for {
		switch concrete := t.(type) {
		case *types.Array,
			*types.Interface,
			*types.TypeParam,
			*types.Union,
			*types.Basic,
			*types.Chan,
			*types.Map,
			*types.Pointer,
			*types.Signature,
			*types.Slice,
			*types.Struct,
			*types.Tuple:
			return nil, nil, false
		case *types.Alias:
			t = concrete.Rhs()
		case *types.Named:
			v, ok := enums[concrete]
			if ok {
				return concrete, v, true
			}
			// Package go/types documentation for (*types.Named).Underlying:
			// "Underlying types are never Named, TypeParam,
			// or Alias types."
			//
			// If this were not the case, then for the
			// following declarations, for the named
			// type B, the value of "B".Underlying()
			// could be the named type A, not int.
			//
			// 	type A int
			// 	type B A
			//
			// The guarantee that (*types.Named).Underlying() will
			// not be a *types.Named value is critical for
			// correctness of the current function. The goal
			// of the function is to determine whether a
			// type (say, B) represented by the argument t is
			// an enumerated type. The function must not consider
			// whether a different named type (i.e., A) is an enumerated type.
			_, isNamed := concrete.Underlying().(*types.Named)
			assert(!isNamed)
			t = concrete.Underlying()
		default:
			// All possible types, as of go1.26, are handled
			// in the cases above. This default case should
			// never run. Nevertheless returning false
			// is the appropriate result in this scenario.
			return nil, nil, false
		}
	}
}

// The result from mustSatisfy is the set of constant values that
// must satisfied to be exhaustive. The result is prepared in
// the context of the current package and the given options.
//
// The keys of the returned map are the set of constant values in
// (go/constant.Value).ExactString representation.
// Note that it is possible for the same constant value to have
// multiple names, i.e. multiple *types.Const.
func mustSatisfy(currentPkg *types.Package, cs []*types.Const, opts *options) map[string][]*types.Const {
	ret := make(map[string][]*types.Const)
	for _, c := range cs {
		if c.Pkg() != currentPkg && !ast.IsExported(c.Name()) {
			continue
		}
		if isPkgLevel(c) && matchAny(opts.excludeConstPatterns, fullname(c)) {
			continue
		}
		ret[c.Val().ExactString()] = append(ret[c.Val().ExactString()], c)
	}
	return ret
}

func formatUnsatisfiedNames(passPkg *types.Package, unsatisfied map[string][]*types.Const) string {
	// Note that the constants are grouped by value in the input.
	// The logic below does the following: The constants in
	// each group are sorted by AST position. The groups
	// themselves are then sorted by AST position (of the
	// first constant in the group).
	groups := make([][]*types.Const, len(unsatisfied))
	i := 0
	for _, cs := range unsatisfied {
		if len(cs) == 0 {
			panic("zero constants for value")
		}
		groups[i] = slices.Clone(cs) // do not modify input
		slices.SortFunc(groups[i], func(a, b *types.Const) int { return cmp.Compare(a.Pos(), b.Pos()) })
		i++
	}
	slices.SortFunc(groups, func(a, b []*types.Const) int { return cmp.Compare(a[0].Pos(), b[0].Pos()) })

	// format to string
	var buf strings.Builder
	for i, g := range groups {
		for i, c := range g {
			var name string
			if passPkg != c.Pkg() && c.Pkg() != nil {
				name = c.Pkg().Name() + "." + c.Name()
			} else {
				name = c.Name()
			}
			buf.WriteString(name)
			if i != len(g)-1 {
				buf.WriteString("|")
			}
		}
		if i != len(groups)-1 {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

func proceedInclExclPattern(t *types.Named, opts *options) bool {
	if len(opts.includeTypePatterns) != 0 {
		if isPkgLevel(t.Obj()) && matchAny(opts.includeTypePatterns, fullname(t.Obj())) {
			return true
		}
		// Include patterns provided and none of
		// them matched.
		return false
	}
	// Exclusion should be checked and performed only if a
	// positive match for explicit inclusion did not happen
	// earlier.
	// This way, if an object name is matched by both include
	// and exclude patterns, the inclusion match wins.
	//
	// Note that the behavior described above is not
	// guaranteed by the analyzer. Users must expect
	// undefined behavior in this scenario.
	if isPkgLevel(t.Obj()) && matchAny(opts.excludeTypePatterns, fullname(t.Obj())) {
		return false
	}
	return true // default is to include
}

func checkMapLiteral(pass *analysis.Pass, opts *options) {
	var (
		in             = pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		enums          = pass.ResultOf[enumerated.Analyzer].(enumerated.Result)
		commentsByFile = make(map[*ast.File]ast.CommentMap)
		printExpr      = func(e ast.Expr) string { return printNode(pass.Fset, ast.Unparen(e)) }
		errorType      = func(e ast.Expr) string { return fmt.Sprintf("%s: could not determine type", printExpr(e)) }
	)

	in.Root().Inspect([]ast.Node{(*ast.CompositeLit)(nil)}, func(c inspector.Cursor) (descend bool) {
		file := head(c.Enclosing((*ast.File)(nil))).Node().(*ast.File)
		if !opts.checkGenerated && ast.IsGenerated(file) {
			return false
		}

		compLit := c.Node().(*ast.CompositeLit)

		mapType, ok := unpointer(pass.TypesInfo.Types[compLit].Type.Underlying()).(*types.Map)
		if !ok {
			return true
		}

		fileComments, ok := commentsByFile[file]
		if !ok {
			fileComments = ast.NewCommentMap(pass.Fset, file, file.Comments)
			commentsByFile[file] = fileComments
		}
		directives, err := parseDirectives(compositeLitComments(pass, fileComments, c))
		if err != nil {
			pass.Reportf(compLit.Pos(), "error parsing comment directives: %s", err)
			return true
		}
		if opts.needEnforceDirective && !directives[dirEnforce] {
			return true
		}
		if directives[dirIgnore] {
			return true
		}

		reasonUnchecked := ""
		if directives[dirEnforce] {
			defer func() {
				if reasonUnchecked != "" {
					pass.Reportf(compLit.Pos(), "enforce directive present and map literal not checked: %s", reasonUnchecked)
				}
			}()
		}

		t, cs, ok := knownEnumeratedType(enums, mapType.Key())
		if !ok {
			reasonUnchecked = "key type is not an enumerated type"
			return true
		}

		if !directives[dirEnforce] {
			if !proceedInclExclPattern(t, opts) {
				return true
			}
		}

		unsatisfied := mustSatisfy(pass.Pkg, cs, opts)
		for _, elt := range compLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				pass.Reportf(elt.Pos(), "%s: not key-value expression", printExpr(elt))
				continue
			}
			tv, ok := pass.TypesInfo.Types[kv.Key]
			if !ok {
				pass.Reportf(kv.Key.Pos(), "%s", errorType(kv.Key))
				continue
			}
			if tv.Value != nil {
				delete(unsatisfied, tv.Value.ExactString())
			}
		}
		if len(unsatisfied) != 0 {
			pass.Reportf(compLit.Pos(), "map literal not exhaustive: missing keys: %s", formatUnsatisfiedNames(pass.Pkg, unsatisfied))
		}
		return true
	})
}

func unpointer(t types.Type) types.Type {
	for {
		switch a := t.(type) {
		case *types.Pointer:
			t = a.Elem()
		default:
			return t
		}
	}
}

func compositeLitComments(pass *analysis.Pass, comments ast.CommentMap, c inspector.Cursor) []*ast.CommentGroup {
	ret := make(map[*ast.CommentGroup]struct{})
	add := func(gs ...*ast.CommentGroup) {
		for _, g := range gs {
			if _, ok := ret[g]; ok {
				continue
			}
			ret[g] = struct{}{}
		}
	}
loop:
	for c := range c.Enclosing() {
		switch n := c.Node().(type) {
		case *ast.AssignStmt:
			add(comments[n]...)
			break loop
		case *ast.ValueSpec:
			add(n.Doc)
			add(n.Comment)
			add(comments[n]...)
			// note: do not break the loop here; also need to
			// see the parent *ast.GenDecl.
		case *ast.GenDecl:
			add(n.Doc)
			add(comments[n]...)
			break loop
		}
	}
	return slices.Collect(maps.Keys(ret)) // current callers do not care about the order
}
