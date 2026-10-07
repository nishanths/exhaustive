package exhaustive

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/types"
	"maps"
	"slices"
	"strings"

	"github.com/nishanths/exhaustive/internal/passes/finder"
	"github.com/nishanths/exhaustive/passes/enumerated"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

func checkSwitch(pass *analysis.Pass, sw *ast.SwitchStmt, comments []*ast.CommentGroup, opts *options) {
	directives, err := parseDirectives(comments)
	if err != nil {
		pass.Reportf(sw.Pos(), "error parsing comment directives: %s", err)
		return
	}
	if opts.checkEnforceOnly && !directives[dirEnforce] {
		return
	}
	if directives[dirIgnore] {
		return
	}

	skip := func(reason string) {
		if directives[dirEnforce] {
			pass.Reportf(sw.Pos(), "enforce directive present and switch not checked: %s", reason)
		}
	}

	// The tag will be nil for switch statements of
	// the form "switch { ... }".
	if sw.Tag == nil {
		skip("no switch expression")
		return
	}

	enums := pass.ResultOf[enumerated.Analyzer].(enumerated.Result)
	t, cs, ok := knownEnumeratedType(enums, pass.TypesInfo.TypeOf(sw.Tag))
	if !ok {
		skip("switch expression type is not an enumerated type")
		return
	}

	// Handle user-specified include/exclude pattern
	// flags. Note that comment directives take
	// precedence over any include/exclude patterns.
	if !directives[dirEnforce] {
		if !proceedInclExclPatterns(t, opts) {
			return
		}
	}

	needDefault := opts.requireDefaultCase
	if v, ok := directives[dirDefrequire]; ok {
		needDefault = v
	}
	foundDefault := false

	need := needValues(pass.Pkg, cs, opts)
	for _, cc := range sw.Body.List {
		cc := cc.(*ast.CaseClause)
		if cc.List == nil {
			foundDefault = true
			if opts.defaultCaseExhaustive {
				clear(need)
				break
			}
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
	if needDefault && !foundDefault {
		pass.Reportf(sw.Pos(), "missing default case")
	}
	if len(need) != 0 {
		pass.Reportf(sw.Pos(), "switch not exhaustive: missing cases: %s", formatMissingNames(pass.Pkg, need))
	}
}

var universeError = types.Universe.Lookup("error")
var universeNil = types.Universe.Lookup("nil")

// shared state for type switch checks in a given package.
type typeSwitchState struct {
	impl map[*types.Interface]map[typename]struct{} // interface -> defined types that implement the interface
	rel  map[*types.TypeName][]typename             // defined type -> type declarations that denote the defined type
}

type typename struct {
	*types.TypeName
	pointer bool
}

func (v typename) String() string {
	return fmt.Sprintf("{%v %v}", v.TypeName.Type(), v.pointer)
}

func checkTypeSwitch(pass *analysis.Pass, sw *ast.TypeSwitchStmt, comments []*ast.CommentGroup, ch *typeSwitchState, opts *options) {
	directives, err := parseDirectives(comments)
	if err != nil {
		pass.Reportf(sw.Pos(), "error parsing comment directives: %s", err)
		return
	}
	if opts.checkEnforceOnly && !directives[dirEnforce] {
		return
	}
	if directives[dirIgnore] {
		return
	}

	skip := func(reason string) {
		if directives[dirEnforce] {
			pass.Reportf(sw.Pos(), "enforce directive present and type switch not checked: %s", reason)
		}
	}

	var x ast.Expr
	switch n := sw.Assign.(type) {
	case *ast.ExprStmt:
		x = ast.Unparen(n.X).(*ast.TypeAssertExpr).X
	case *ast.AssignStmt:
		x = ast.Unparen(n.Rhs[0]).(*ast.TypeAssertExpr).X
	default:
		pass.Reportf(sw.Assign.Pos(), "unsupported type switch syntax (type %T)", n) // not possible as of go1.27
		return
	}

	// Go spec section on Type switches: "As with type
	// assertions, x must be of interface type".
	swInterface, ok := pass.TypesInfo.TypeOf(x).Underlying().(*types.Interface)
	if !ok {
		skip("expression in type assertion is not interface")
		return
	}

	if !directives[dirEnforce] {
		// Special case: Type switches of the empty interface
		// are not checked.
		// Note that all types implement the empty interface.
		if swInterface.Empty() {
			return
		}
		// Special case: Type switches of the error built-in interface
		// are not checked.
		if n, ok := types.Unalias(pass.TypesInfo.TypeOf(x)).(*types.Named); ok && n.Obj() == universeError {
			return
		}
	}

	r := pass.ResultOf[finder.Analyzer].(*finder.Result)
	computeImplements(ch, r, swInterface)
	computeTypeRelations(ch, r)
	need := needTypes(pass.Pkg, ch, swInterface, opts)

	needDefault := opts.requireDefaultCase
	if v, ok := directives[dirDefrequire]; ok {
		needDefault = v
	}
	foundDefault := false
	foundNil := false

	for _, cc := range sw.Body.List {
		cc := cc.(*ast.CaseClause)
		if cc.List == nil {
			foundDefault = true
			if opts.defaultCaseExhaustive {
				clear(need)
				break
			}
		}
		for _, expr := range cc.List {
			t := pass.TypesInfo.TypeOf(expr)
			if t == nil {
				continue
			}
			if !foundNil {
				if obj, ok := t.(types.Object); ok && obj == universeNil {
					foundNil = true
				}
			}
			// Note: Go spec as of go1.27: "[The receiver's] type must be a defined type T
			// or a pointer to a defined type T".
			// See also: Related notes in func computeImplements.
			switch t := types.Unalias(t).(type) {
			case *types.Named:
				n := t
				delete(need, typename{n.Obj(), false})
			case *types.Pointer:
				if n, ok := t.Elem().(*types.Named); ok {
					delete(need, typename{n.Obj(), true})
				}
			}
		}
	}

	if needDefault && !foundDefault {
		pass.Reportf(sw.Pos(), "missing default case")
	}

	if len(need) != 0 {
		matchpointer := func(defined typename) func(typename) typename {
			return func(name typename) typename {
				return typename{name.TypeName, cmp.Or(defined.pointer, name.pointer)}
			}
		}
		m := make(map[typename][]typename, len(need))
		for k := range need {
			m[k] = slicemap(ch.rel[k.TypeName], matchpointer(k))
		}
		pass.Reportf(sw.Pos(), "type switch not exhaustive: missing cases: %s", formatMissingTypes(pass.Pkg, m, opts.requireNilCase && !foundNil))
	}
}

func computeImplements(ch *typeSwitchState, r *finder.Result, swInterface *types.Interface) {
	if _, ok := ch.impl[swInterface]; ok {
		return
	}
	ch.impl[swInterface] = make(map[typename]struct{})

	add := func(n *types.Named, pointer bool) {
		ch.impl[swInterface][typename{n.Obj(), pointer}] = struct{}{}
	}

	// Note: None of the non-interface predeclared types have methods
	// as of go1.27 and thus they cannot implement a
	// non-empty basic interface.
	// If this does not hold true in future Go versions,
	// then those predeclared types must be considered
	// for the following "implements" checks.
	// A similar note applies to func computeTypeRelations.
	//
	// Note: Go spec as of go1.27: "[The receiver's] type must be a defined type T
	// or a pointer to a defined type T" and "T is called the receiver base type.
	// A receiver base type cannot be a pointer or interface type".
	//
	// Examples
	//
	//   type S int   // S and *S: valid receiver types
	//   type U *S    // U: invalid receiver type (U is a pointer type)
	//   type V *int  // V: invalid receiver type (V is a pointer type)
	//   type A = *S  // A: valid receiver type
	//   type B = *A  // B: invalid receiver type (B is a 2x-pointer to a defined type)
	//
	// Note: See also: Go spec section 'Method sets'.
	for tn := range r.TypeDecls {
		switch t := types.Unalias(tn.Type()).(type) {
		case *types.Named:
			n := t
			if !isInterface(n) {
				if types.Implements(types.NewPointer(n), swInterface) {
					add(n, true)
				}
				if types.Implements(n, swInterface) {
					add(n, false)
				}
			}
		case *types.Pointer:
			if n, ok := t.Elem().(*types.Named); ok && !isInterface(n) && types.Implements(t, swInterface) {
				add(n, true)
			}
		}
	}
}

func computeTypeRelations(ch *typeSwitchState, r *finder.Result) {
	if ch.rel != nil {
		return
	}
	ch.rel = make(map[*types.TypeName][]typename, 128)

	for tn := range r.TypeDecls {
		switch t := tn.Type().(type) {
		case *types.Alias:
			switch tt := types.Unalias(t).(type) {
			case *types.Named:
				n := tt
				ch.rel[n.Obj()] = append(ch.rel[n.Obj()], typename{tn, false})
			case *types.Pointer:
				if n, ok := tt.Elem().(*types.Named); ok {
					ch.rel[n.Obj()] = append(ch.rel[n.Obj()], typename{tn, true})
				}
			}
		case *types.Named:
			n := t
			ch.rel[n.Obj()] = append(ch.rel[n.Obj()], typename{tn, false})
		default:
			panic(fmt.Sprintf("internal error: unexpected type %T", t))
		}
	}
}

func isInterface(n *types.Named) bool {
	_, ok := n.Underlying().(*types.Interface)
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
	ret := make(map[string][]*types.Const)
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

func needTypes(currentPkg *types.Package, ch *typeSwitchState, swInterface *types.Interface, opts *options) map[typename]struct{} {
	ret := make(map[typename]struct{}, len(ch.impl[swInterface]))
	for defined := range ch.impl[swInterface] {
		for _, name := range ch.rel[defined.TypeName] {
			if name.Pkg() != currentPkg && !name.Exported() {
				continue
			}
			ret[defined] = struct{}{}
			break
		}
	}
	return ret
}

func formatMissingNames(currentPkg *types.Package, missing map[string][]*types.Const) string {
	// Note that the constants are grouped by value in the input.
	// The logic below does the following: The constants in
	// each group are sorted by AST position. The groups
	// themselves are then sorted by AST position (of the
	// first constant in the group).
	groups := make([][]*types.Const, len(missing))
	i := 0
	for _, cs := range missing {
		assert(len(cs) != 0)         // unexpectedly zero constants for value
		groups[i] = slices.Clone(cs) // do not modify input
		slices.SortFunc(groups[i], func(a, b *types.Const) int { return cmp.Compare(a.Pos(), b.Pos()) })
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

func formatMissingTypes(currentPkg *types.Package, missing map[typename][]typename, missingNil bool) string {
	type group struct {
		defined typename
		names   []typename
	}

	cmpfn := func(a, b typename) int {
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

	typenameString := func(t typename) string {
		if t.pointer {
			return "*" + nameString(currentPkg, t)
		}
		return nameString(currentPkg, t)
	}

	groups := make([]group, len(missing))
	i := 0
	for defined, names := range missing {
		assert(len(names) != 0)      // unexpectedly zero type names for defined type
		names := slices.Clone(names) // do not modify input
		slices.SortFunc(names, cmpfn)
		groups[i] = group{defined, names}
		i++
	}
	slices.SortFunc(groups, func(a, b group) int { return cmpfn(a.defined, b.defined) })

	// format to string
	var buf strings.Builder
	if missingNil {
		buf.WriteString("nil")
		if len(groups) > 0 {
			buf.WriteString(", ")
		}
	}
	for i, g := range groups {
		buf.WriteString(typenameString(g.defined))
		if len(g.names) != 0 && g.names[0] != g.defined {
			buf.WriteString(" (")
			for i, tn := range g.names {
				buf.WriteString(typenameString(tn))
				if i != len(g.names)-1 {
					buf.WriteString("|")
				}
			}
			buf.WriteString(")")
		}
		if i != len(groups)-1 {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

func proceedInclExclPatterns(t *types.Named, opts *options) bool {
	if len(opts.includeTypePatterns) != 0 {
		if isPkgLevel(t.Obj()) && matchAny(opts.includeTypePatterns, fullname(t.Obj())) {
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
	if isPkgLevel(t.Obj()) && matchAny(opts.excludeTypePatterns, fullname(t.Obj())) {
		return false
	}
	return true // default is to include
}

func checkMapLiteral(pass *analysis.Pass, compLit *ast.CompositeLit, comments []*ast.CommentGroup, opts *options) {
	mapType, ok := unpointer(pass.TypesInfo.Types[compLit].Type.Underlying()).(*types.Map)
	if !ok {
		return
	}

	directives, err := parseDirectives(comments)
	if err != nil {
		pass.Reportf(compLit.Pos(), "error parsing comment directives: %s", err)
		return
	}
	if opts.checkEnforceOnly && !directives[dirEnforce] {
		return
	}
	if directives[dirIgnore] {
		return
	}

	skip := func(reason string) {
		if directives[dirEnforce] {
			pass.Reportf(compLit.Pos(), "enforce directive present and map literal not checked: %s", reason)
		}
	}

	enums := pass.ResultOf[enumerated.Analyzer].(enumerated.Result)
	t, cs, ok := knownEnumeratedType(enums, mapType.Key())
	if !ok {
		skip("key type is not an enumerated type")
		return
	}

	if !directives[dirEnforce] {
		if !proceedInclExclPatterns(t, opts) {
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
		pass.Reportf(compLit.Pos(), "map literal not exhaustive: missing keys: %s", formatMissingNames(pass.Pkg, need))
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

func nameString(currentPkg *types.Package, obj types.Object) string {
	if currentPkg != obj.Pkg() && obj.Pkg() != nil {
		return obj.Pkg().Name() + "." + obj.Name()
	} else {
		return obj.Name()
	}
}
