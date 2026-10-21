package exhaustive

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/types"
	"slices"
	"strings"

	"github.com/nishanths/exhaustive/internal/passes/finder"
	"github.com/nishanths/exhaustive/passes/enumerated"

	"golang.org/x/tools/go/analysis"
)

func checkExprSwitch(pass *analysis.Pass, sw *ast.SwitchStmt, comments []*ast.CommentGroup, opts *options) {
	directives, proceed, err := proceedCommentDirectives(comments, opts)
	if err != nil {
		pass.Reportf(sw.Pos(), "%s", err)
		return
	}
	if !proceed {
		return
	}

	skip := func(reason string) {
		if directives[dirEnforce] {
			pass.Reportf(sw.Pos(), "enforce directive present and expression switch not checked: %s", reason)
		}
	}

	// The tag will be nil for switch statements of
	// the form "switch { ... }".
	if sw.Tag == nil {
		skip("no switch expression")
		return
	}

	enums := pass.ResultOf[enumerated.Analyzer].(enumerated.Result)
	_, cs, ok := knownEnumeratedType(enums, pass.TypesInfo.TypeOf(sw.Tag))
	if !ok {
		skip("switch expression type is not an enumerated type")
		return
	}

	// Handle user-specified include/exclude pattern
	// flags. Note that comment directives take
	// precedence over any include/exclude patterns.
	if !directives[dirEnforce] {
		if tn, ok := hasTypeName(pass.TypesInfo.TypeOf(sw.Tag)); ok && !proceedInclExclPatterns(tn, opts) {
			return
		}
	}

	need := needValues(pass.Pkg, cs, opts)

	for _, cc := range sw.Body.List {
		cc := cc.(*ast.CaseClause)
		if cc.List == nil {
			if opts.defaultCaseExhaustive {
				clear(need)
				break
			}
			continue
		}
		for _, expr := range cc.List {
			tv, ok := pass.TypesInfo.Types[expr]
			if !ok {
				continue
			}
			if tv.Value != nil {
				delete(need, tv.Value.ExactString())
			}
		}
	}

	if len(need) != 0 {
		pass.Reportf(sw.Pos(), "missing cases in expression switch: %s", formatMissingNames(pass.Pkg, need))
	}


	requireDefault := opts.requireDefaultCase
	if v, ok := directives[dirDefrequire]; ok {
		requireDefault = v
	}
	if requireDefault {
		checkDefaultCase(pass, sw)
	}
}

func checkDefaultCase(pass *analysis.Pass, sw *ast.SwitchStmt) {
	foundDefault := false
	for _, cc := range sw.Body.List {
		if cc.(*ast.CaseClause).List == nil {
			foundDefault = true
			break
		}
	}
	if !foundDefault {
		pass.Reportf(sw.Pos(), "missing default case")
	}
}

// hasTypeName determines whether t is a type with a name.
func hasTypeName(t types.Type) (*types.TypeName, bool) {
	switch t := t.(type) {
	case *types.Alias:
		return t.Obj(), true
	case *types.Named: // also represents predeclared type error
		return t.Obj(), true
	case *types.TypeParam:
		return t.Obj(), true
	default:
		return nil, false
	}
}

var universeError = types.Universe.Lookup("error")

// shared checker state for all type switches in a particular package.
type typeSwitchState struct {
	// The 'map[tp]struct{}' value is the set of types T
	// or *T that implement the interface in the key,
	// where T is a non-interface defined type.
	impl map[*types.Interface]map[tp]struct{}

	// The '[]tp' value is the set of types
	// that are identical to the type in the key.
	// One of the type names in '[]tp' will be the
	// type name of the defined type in the key itself.
	//
	// The type name in the key is always that of a defined type.
	// See func computeTypeRelations for details.
	// Only those relationships relevant for type switch
	// analysis are recorded in this map.
	identical map[tp][]tp
}

// tp represents a type T or type *T,
// where T is a type name.
type tp struct {
	*types.TypeName
	pointer bool
}

func checkTypeSwitch(pass *analysis.Pass, sw *ast.TypeSwitchStmt, comments []*ast.CommentGroup, ch *typeSwitchState, opts *options) {
	directives, proceed, err := proceedCommentDirectives(comments, opts)
	if err != nil {
		pass.Reportf(sw.Pos(), "%s", err)
		return
	}
	if !proceed {
		return
	}

	skip := func(reason string) {
		if directives[dirEnforce] {
			pass.Reportf(sw.Pos(), "enforce directive present and type switch not checked: %s", reason)
		}
	}

	var swExpr ast.Expr
	switch n := sw.Assign.(type) {
	case *ast.ExprStmt:
		swExpr = ast.Unparen(n.X).(*ast.TypeAssertExpr).X
	case *ast.AssignStmt:
		swExpr = ast.Unparen(n.Rhs[0]).(*ast.TypeAssertExpr).X
	default:
		pass.Reportf(sw.Assign.Pos(), "unsupported type switch syntax (%T)", n) // impossible as of go1.27
		return
	}

	swInterface, ok := pass.TypesInfo.TypeOf(swExpr).Underlying().(*types.Interface)
	if !ok || !swInterface.IsMethodSet() {
		// Go spec section on Type switches: "As with
		// type assertions, x must be of interface type,
		// but not a type parameter".
		//
		// This branch is impossible as of go1.27,
		// but be future-proof for correctness.
		// The IsMethodSet requirement is essential for
		// correctness because the logic below only
		// evaluates types that are valid method receiver types.
		skip("unsupported expression in type assertion")
		return
	}

	if !directives[dirEnforce] {
		tn, ok := hasTypeName(pass.TypesInfo.TypeOf(swExpr))
		if ok && !proceedInclExclPatterns(tn, opts) {
			return
		}
		// Special case: Type switches of the empty
		// interface are not checked. Note that all types
		// implement the empty interface.
		if swInterface.Empty() {
			return
		}
		// Special case: Type switches of the error
		// built-in interface are not checked.
		if ok && tn == universeError {
			return
		}
	}

	r := pass.ResultOf[finder.Analyzer].(*finder.Result)
	computeImplements(ch, r, swInterface)
	computeTypeRelations(ch, r)
	need := needTypes(pass.Pkg, ch, swInterface, opts)

	reportMissingNil := opts.requireCaseNil

	for _, cc := range sw.Body.List {
		cc := cc.(*ast.CaseClause)
		if cc.List == nil {
			if opts.defaultCaseExhaustive {
				reportMissingNil = false
				clear(need)
				break
			}
			continue
		}

		for _, caseExpr := range cc.List {
			caseType := pass.TypesInfo.TypeOf(caseExpr)
			if caseType == nil {
				continue
			}

			if caseType == types.Typ[types.UntypedNil] {
				reportMissingNil = false
			}

			// Go spec: "If T is an interface type, x.(T) asserts that
			// the dynamic type of x implements the interface T."
			//
			// Note: It is possible in the language for a type to
			// implement an interface T. It is impossible for a type
			// to implement *T where T is an interface. Thus it is
			// sufficient to only consider interface types (but not,
			// for example, pointers to interface types).
			//
			// The following code correctly handles both traditional
			// interfaces and general interfaces (from a type parameter).
			// It is not necessary to handle e.g. []P, where P is a
			// type parameter, because these are not of underlying
			// type interface and it is impossible for a type to
			// implement them.
			if i, ok := caseType.Underlying().(*types.Interface); ok {
				computeImplements(ch, r, i)
				for t := range ch.impl[i] {
					delete(need, t)
				}
			}

			switch t := types.Unalias(caseType).(type) {
			case *types.Named:
				n := t
				delete(need, tp{n.Obj(), false})
			case *types.Pointer:
				if n, ok := types.Unalias(t.Elem()).(*types.Named); ok {
					delete(need, tp{n.Obj(), true})
				}
			}
		}
	}

	if reportMissingNil || len(need) != 0 {
		pass.Reportf(sw.Pos(), "missing cases in type switch: %s", formatMissingTypes(pass.Pkg, need, reportMissingNil))
	}

	if opts.checkImpossible {
		checkUnreachableCases(pass, sw)
	}
}

func checkUnreachableCases(pass *analysis.Pass, sw *ast.TypeSwitchStmt) {
	// Note: This function does not check for issues that are
	// already reported by "go build".
	// As such, the following categories of issues
	// are not checked: identical types in type switch cases;
	// impossible type in type switch case.

	earlier := make(map[*types.Interface]ast.Expr) // interfaces in earlier case clauses
	current := make(map[*types.Interface]ast.Expr) // interfaces in current case clause

	for _, cc := range sw.Body.List {
		clear(current)

		for _, caseExpr := range cc.(*ast.CaseClause).List {
			caseType := pass.TypesInfo.TypeOf(caseExpr)
			if caseType == nil {
				continue
			}
			if i, ok := caseType.Underlying().(*types.Interface); ok {
				if _, ok := current[i]; !ok {
					current[i] = caseExpr
				}
			}
			for k, v := range earlier {
				if types.Implements(caseType, k) {
					pass.Reportf(caseExpr.Pos(), "unreachable type switch case %s (matched by earlier case %s)", formatExpr(caseExpr), formatExpr(v))
				}
			}
		}

		for k, v := range current {
			if _, ok := earlier[k]; !ok {
				earlier[k] = v
			}
		}
	}
}

func computeImplements(ch *typeSwitchState, r *finder.Result, intf *types.Interface) {
	if _, ok := ch.impl[intf]; ok {
		return
	}

	if ch.impl == nil {
		ch.impl = make(map[*types.Interface]map[tp]struct{})
	}
	ch.impl[intf] = make(map[tp]struct{})

	add := func(d tp) {
		ch.impl[intf][d] = struct{}{}
	}

	// Note: None of the non-interface predeclared types have
	// methods as of go1.27 and thus they cannot implement a
	// non-empty interface. If this does not hold true
	// in future Go versions, then those predeclared types
	// must be evaluated for the following "implements" checks.
	// A similar note applies in func computeTypeRelations.
	//
	// Note: Go spec as of go1.27: "[The receiver's] type
	// must be a defined type T or a pointer to a defined
	// type T" and "T is called the receiver base type. A
	// receiver base type cannot be a pointer or interface type".
	//
	// Examples
	//
	//   type S int   // S and *S: valid receiver types
	//   type U *S    // U: invalid receiver type (U is a pointer type)
	//   type V *int  // V: invalid receiver type (V is a pointer type)
	//   type R V     // R: invalid receiver type (R is a pointer type)
	//   type A = *S  // A: valid receiver type
	//   type B = *A  // B: invalid receiver type (B is a 2x-pointer to a defined type)
	//
	// Note: See also: Go spec section 'Method sets'.
	for tn := range r.TypeDecls {
		switch t := types.Unalias(tn.Type()).(type) {
		case *types.Named: // defined type
			n := t
			// The following check for an interface type is not due
			// to the "receiver based type cannot be a pointer or
			// interface type" requirement in the Go spec. (That
			// requirement will be handled internally
			// anyway by func types.Implements.)
			// Rather, the interface check is specific to the logic
			// of the analysis, which should not consider interfaces
			// as a possible implementing type in this context.
			if !underlyingIs[*types.Interface](n) && n.Origin().TypeParams() == nil {
				if types.Implements(types.NewPointer(n), intf) {
					add(tp{n.Obj(), true})
				}
				if types.Implements(n, intf) {
					add(tp{n.Obj(), false})
				}
			}
		case *types.Pointer: // pointer to a defined type
			if n, ok := types.Unalias(t.Elem()).(*types.Named); ok {
				if !underlyingIs[*types.Interface](n) && n.Origin().TypeParams() == nil {
					if types.Implements(types.NewPointer(n), intf) {
						add(tp{n.Obj(), true})
					}
				}
			}
		}
	}
}

func computeTypeRelations(ch *typeSwitchState, r *finder.Result) {
	if ch.identical != nil {
		return
	}

	ch.identical = make(map[tp][]tp, len(r.TypeDecls))
	add := func(root, t2 tp) {
		ch.identical[root] = append(ch.identical[root], t2)
	}

	for tn := range r.TypeDecls {
		switch t := tn.Type().(type) {
		case *types.Alias:
			switch tt := types.Unalias(t).(type) {
			case *types.Named:
				n := tt
				if n.Origin().TypeParams() == nil {
					add(tp{n.Obj(), false}, tp{tn, false})
					add(tp{n.Obj(), true}, tp{tn, true})
				}
			case *types.Pointer:
				if n, ok := types.Unalias(tt.Elem()).(*types.Named); ok {
					if n.Origin().TypeParams() == nil {
						add(tp{n.Obj(), true}, tp{tn, false})
					}
				}
			}
		case *types.Named:
			n := t
			if n.Origin().TypeParams() == nil {
				add(tp{n.Obj(), false}, tp{tn, false})
				add(tp{n.Obj(), true}, tp{tn, true})
			}
		default:
			panic(fmt.Sprintf("internal error: unexpected type %T", t))
		}
	}
}

func underlyingIs[T types.Type](n *types.Named) bool {
	_, ok := n.Underlying().(T)
	return ok
}

// knownEnumeratedType reports whether the given type (which
// typically is the type for an arbitrary expression) is an
// enumerated type. It returns the enumerated type and the
// enumerated constants of that type.
func knownEnumeratedType(enums enumerated.Result, t types.Type) (*types.Named, []*types.Const, bool) {
	for {
		switch concrete := t.(type) {
		case *types.Array,
			*types.Basic,
			*types.Chan,
			*types.Interface,
			*types.Map,
			*types.Pointer,
			*types.Signature,
			*types.Slice,
			*types.Struct,
			*types.Tuple,
			*types.TypeParam,
			*types.Union:
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
			// not be a *types.Named value is essential for
			// correctness of the current function. The goal
			// of the function is to determine whether a
			// type (say, B) represented by the argument t is
			// an enumerated type. The function must not consider
			// whether a different named type (i.e., A) is an enumerated type.
			if _, ok := concrete.Underlying().(*types.Named); ok {
				return nil, nil, false
			}
			t = concrete.Underlying()
		default:
			// The possible types, as of go1.27, are handled
			// in the cases above and the default case should
			// not run. Nevertheless returning false is the appropriate
			// result for any unhandled type.
			return nil, nil, false
		}
	}
}

// The result from needValues is the set of constant values that
// must be satisfied to be exhaustive. The result is prepared in
// the context of the current package and the given options.
//
// The keys of the returned map are the set of constant values in
// (go/constant.Value).ExactString representation.
// Note that it is possible for the same constant value to have
// multiple names, i.e. multiple *types.Const.
func needValues(currentPkg *types.Package, cs []*types.Const, opts *options) map[string][]*types.Const {
	ret := make(map[string][]*types.Const, len(cs)) // usually one name per value
	for _, c := range cs {
		if c.Pkg() != currentPkg && !c.Exported() {
			continue
		}
		if isPkgLevel(c) && matchAny(opts.excludeConstPatterns, fullname(c)) {
			continue
		}
		ret[c.Val().ExactString()] = append(ret[c.Val().ExactString()], c)
	}
	return ret
}

func needTypes(currentPkg *types.Package, ch *typeSwitchState, swInterface *types.Interface, opts *options) map[tp][]tp {
	ret := make(map[tp][]tp, len(ch.impl[swInterface]))
	for t := range ch.impl[swInterface] {
		for _, name := range ch.identical[t] {
			if name.Pkg() != currentPkg && !name.Exported() {
				continue
			}
			ret[t] = append(ret[t], name)
		}
	}
	return ret
}

// formats constant names for an expression switch diagnostic.
func formatMissingNames(currentPkg *types.Package, missing map[string][]*types.Const) string {
	// Note that the constants are grouped by value in the input.
	// The logic below does the following: The constants in
	// each group are sorted by AST position. The groups
	// themselves are then sorted by AST position (of the
	// first constant in the group).
	groups := make([][]*types.Const, len(missing))
	i := 0
	for _, cs := range missing {
		cs := slices.Clone(cs) // do not modify input
		slices.SortFunc(cs, func(a, b *types.Const) int { return cmp.Compare(a.Pos(), b.Pos()) })
		groups[i] = cs
		i++
	}
	slices.SortFunc(groups, func(a, b []*types.Const) int { return cmp.Compare(a[0].Pos(), b[0].Pos()) })

	// format to string
	var buf strings.Builder
	for i, g := range groups {
		for i, c := range g {
			buf.WriteString(nameString(currentPkg, c))
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

// formats type names for a type switch diagnostic.
func formatMissingTypes(currentPkg *types.Package, missing map[tp][]tp, missingNil bool) string {
	cmpfn := func(a, b tp) int {
		switch {
		case !a.IsAlias() && b.IsAlias():
			return -1
		case a.IsAlias() && !b.IsAlias():
			return 1
		}
		switch {
		case a.Pkg() == nil && b.Pkg() != nil:
			return -1
		case a.Pkg() != nil && b.Pkg() == nil:
			return 1
		case a.Pkg() == currentPkg && b.Pkg() != currentPkg:
			return -1
		case a.Pkg() != currentPkg && b.Pkg() == currentPkg:
			return 1
		case a.Pkg() != b.Pkg():
			return cmp.Compare(a.Pkg().Path(), b.Pkg().Path())
		}
		if a.Name() != b.Name() {
			return cmp.Compare(a.Name(), b.Name())
		}
		switch {
		case !a.pointer && b.pointer:
			return -1
		case a.pointer && !b.pointer:
			return 1
		}
		return 0
	}

	tpString := func(t tp) string {
		if t.pointer {
			return "*" + nameString(currentPkg, t)
		}
		return nameString(currentPkg, t)
	}

	type group struct {
		t     tp
		names []tp
	}
	groups := make([]group, len(missing))
	i := 0
	for t, names := range missing {
		names := slices.Clone(names) // do not modify input
		slices.SortFunc(names, cmpfn)
		groups[i] = group{t, names}
		i++
	}
	slices.SortFunc(groups, func(a, b group) int { return cmpfn(a.t, b.t) })

	// format to string
	var buf strings.Builder
	if missingNil {
		buf.WriteString("nil")
		if len(groups) > 0 {
			buf.WriteString(", ")
		}
	}
	for i, g := range groups {
		wrote := false
		for _, tt := range g.names {
			if wrote {
				buf.WriteString(" or ")
			}
			buf.WriteString(tpString(tt))
			wrote = true
		}
		if i != len(groups)-1 {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

func proceedCommentDirectives(comments []*ast.CommentGroup, opts *options) (map[directive]bool, bool, error) {
	directives, err := parseDirectives(comments)
	if err != nil {
		return nil, false, fmt.Errorf("error parsing comment directives: %s", err)
	}
	if opts.checkEnforceOnly && !directives[dirEnforce] {
		return directives, false, nil
	}
	if directives[dirIgnore] {
		return directives, false, nil
	}
	return directives, true, nil
}

func proceedInclExclPatterns(obj types.Object, opts *options) bool {
	if len(opts.includeTypePatterns) != 0 {
		if isPkgLevel(obj) && matchAny(opts.includeTypePatterns, fullname(obj)) {
			return true
		}
		// Include patterns were provided and none of
		// them matched.
		return false
	}
	// Exclusion should be checked for and performed only if a
	// positive match for explicit inclusion did not happen
	// earlier.
	// This way, if an object name is matched by both include
	// and exclude patterns, the inclusion match wins.
	// This implementation matches the behavior described
	// in the package comment.
	if isPkgLevel(obj) && matchAny(opts.excludeTypePatterns, fullname(obj)) {
		return false
	}
	return true // default is to include
}

func checkMapLit(pass *analysis.Pass, compLit *ast.CompositeLit, comments []*ast.CommentGroup, opts *options) {
	mapType, ok := unpointer(pass.TypesInfo.Types[compLit].Type.Underlying()).(*types.Map)
	if !ok {
		return
	}

	directives, proceed, err := proceedCommentDirectives(comments, opts)
	if err != nil {
		pass.Reportf(compLit.Pos(), "%s", err)
		return
	}
	if !proceed {
		return
	}

	skip := func(reason string) {
		if directives[dirEnforce] {
			pass.Reportf(compLit.Pos(), "enforce directive present and map literal not checked: %s", reason)
		}
	}

	enums := pass.ResultOf[enumerated.Analyzer].(enumerated.Result)
	_, cs, ok := knownEnumeratedType(enums, mapType.Key())
	if !ok {
		skip("key type is not an enumerated type")
		return
	}

	if !directives[dirEnforce] {
		if tn, ok := hasTypeName(mapType.Key()); ok && !proceedInclExclPatterns(tn, opts) {
			return
		}
	}

	need := needValues(pass.Pkg, cs, opts)
	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		tv, ok := pass.TypesInfo.Types[kv.Key]
		if !ok {
			continue
		}
		if tv.Value != nil {
			delete(need, tv.Value.ExactString())
		}
	}
	if len(need) != 0 {
		pass.Reportf(compLit.Pos(), "missing keys in map literal: %s", formatMissingNames(pass.Pkg, need))
	}
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

func nameString(currentPkg *types.Package, obj types.Object) string {
	if currentPkg != obj.Pkg() && obj.Pkg() != nil {
		return obj.Pkg().Name() + "." + obj.Name()
	} else {
		return obj.Name()
	}
}

func formatExpr(e ast.Expr) string {
	return types.ExprString(e)
}
