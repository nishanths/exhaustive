package exhaustive

import (
	"bytes"
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
)

// the main function that checks an expression switch.
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

	r := pass.ResultOf[enumerated.Analyzer].(enumerated.Result)
	_, cs, ok := isEnumeratedType(r, pass.TypesInfo.TypeOf(sw.Tag))
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
		pass.Reportf(sw.Pos(), "missing cases in expression switch: %s", formatMissingConstants(pass.Pkg, need))
	}

	checkDefault := func() bool {
		if v, ok := directives[dirDefrequire]; ok {
			return v
		}
		return opts.requireDefaultCase
	}()
	if checkDefault {
		checkHasDefaultCase(pass, sw)
	}
}

// checks whether the given expression switch includes a default case
func checkHasDefaultCase(pass *analysis.Pass, sw *ast.SwitchStmt) {
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

// shared checker state for all type switches in a particular package.
type typeSwitchState struct {
	impl map[*types.Interface]map[implementer]struct{}

	// Only those relationships relevant to type switch
	// analysis are recorded in this map. See func
	// computeTypeRelations for details.
	// At least one of the type names in the value slice will
	// be the type name in the key.
	identical map[tp][]tp
}

type tp struct {
	typename *types.TypeName // type name T
	pointer  bool            // whether T or *T
}

// An implementer represents a non-interface type that can
// implement an interface.
type implementer struct {
	named   *types.Named // non-interface defined type T; currently non-generic only.
	pointer bool         // whether T or *T
}

func (t *implementer) Type() types.Type {
	if t.pointer {
		return types.NewPointer(t.named)
	}
	return t.named
}

var universeError = types.Universe.Lookup("error")

// the main function that checks a type switch.
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
		// Go spec as of go1.27: "Interfaces that are not basic may
		// only be used as type constraints, or as elements of
		// other interfaces used as constraints. They cannot be
		// the types of values or variables, or components of
		// other, non-interface types."
		//
		// and Go spec on Type switches: "As with type
		// assertions, x must be of interface type, but
		// not a type parameter".
		//
		// Thus this branch is impossible as of go1.27.
		skip("unsupported expression in type switch guard")
		return
	}

	if !directives[dirEnforce] {
		tn, ok := hasTypeName(pass.TypesInfo.TypeOf(swExpr))
		if ok && !proceedInclExclPatterns(tn, opts) {
			return
		}
	}

	checkMissingCases := func() bool {
		// Special case: Type switches of the empty
		// interface are not checked for missing cases.
		// Note that all types implement the empty interface.
		if swInterface.Empty() {
			skip("unsupported: interface type of switch expression is the empty interface")
			return false
		}
		// Special case: Type switches of the error built-in interface
		// are not checked for missing cases.
		if !directives[dirEnforce] {
			tn, ok := hasTypeName(types.Unalias(pass.TypesInfo.TypeOf(swExpr)))
			if ok && tn == universeError {
				return false
			}
		}
		return true
	}()

	if checkMissingCases {
		checkTypeSwitchCases(pass, sw, swInterface, ch, opts)
	}

	checkTypeSwitchIneff(pass, sw)
}

// checks that necessary types are listed in the cases of a type switch
func checkTypeSwitchCases(pass *analysis.Pass, sw *ast.TypeSwitchStmt, swInterface *types.Interface, ch *typeSwitchState, opts *options) {
	r := pass.ResultOf[finder.Analyzer].(*finder.Result)
	computeImplements(ch, r, swInterface)
	computeTypeRelations(ch, r)
	need := needTypes(pass.Pkg, swInterface, ch)

	reportMissingNil := opts.requireNilCase

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

			if caseIntf, typeParam, ok := interfaceTypeOf(caseType); ok {
				if !typeParam {
					// Remove the types in the 'need' map that
					// implement the case interface type.
					for k := range need {
						if types.Implements(k.Type(), caseIntf) {
							delete(need, k)
						}
					}
				}
				continue
			}

			switch t := types.Unalias(caseType).(type) {
			case *types.Named:
				n := t
				if n.TypeParams() == nil {
					delete(need, implementer{n.Origin(), false})
				}
			case *types.Pointer:
				if n, ok := types.Unalias(t.Elem()).(*types.Named); ok {
					if n.TypeParams() == nil {
						delete(need, implementer{n.Origin(), true})
					}
				}
			}
		}
	}

	if reportMissingNil || len(need) != 0 {
		m := make(map[implementer][]tp, len(need))
		for t := range need {
			for _, tt := range ch.identical[tp{t.named.Obj(), t.pointer}] {
				if !canAccessObject(pass.Pkg, tt.typename) {
					continue
				}
				m[t] = append(m[t], tt)
			}
		}
		pass.Reportf(sw.Pos(), "missing cases in type switch: %s", formatMissingTypes(pass.Pkg, m, reportMissingNil))
	}
}

// checks that types listed in the cases of a type switch are not ineffectual
func checkTypeSwitchIneff(pass *analysis.Pass, sw *ast.TypeSwitchStmt) {
	// This function does not check for issues that are
	// already reported as errors by "go build". As such, the
	// following categories of issues are not checked:
	// identical types in type switch cases; impossible type
	// in type switch case.

	type intfcase struct {
		intf      *types.Interface
		typeParam bool
		expr      ast.Expr
	}

	reported := make(map[ast.Expr]bool) // case expressions of types reported as ineffectual

	// Go spec: "A type parameter or a generic type may be
	// used as a type in a case. If upon instantiation that
	// type turns out to duplicate another entry in the
	// switch, the first matching case is chosen."
	//
	// Go spec: "The 'fallthrough' statement is not permitted
	// in a type switch."
	//
	// If the interface type C is a type constraint for a type parameter,
	// the actual type of C depends upon the type argument
	// used upon instantiation.
	// In this scenario, a true return value from types.Implements(T, C),
	// where T is the type in the current case, is insufficient to
	// say that the type T is ineffectual (it may
	// only be so for certain instantiations).
	//
	// For instance, given this code
	//
	//	type t1 struct{}
	//	func (t1) o() {}
	//	type t2 struct{}
	//	func (t2) o() {}
	//
	//	func f[P interface{o()}]() {
	//	     var x interface{o()} = t2{}
	//	     switch x.(type) {
	//	     case P:
	//	     case t2:
	//	     case interface{o()}:
	//	     }
	//	}
	//
	//	func main() {
	//		f[t1]()
	//	}
	//
	// note that for the "f[t1]()" call, the types
	// in "case t2" and "case interface{o()}" are not
	// rendered ineffectual by the earlier type in "case P".
	ineffectualOrder := func(e intfcase, currentCaseType types.Type) bool {
		if e.typeParam {
			// This scenario is not currently supported. But a more
			// rigorous analysis could report ineffectual types
			// in this scenario.
			//
			// For instance, given a type constraint interface C for
			// a type parameter P specified as the type in an
			// earlier case, if the type set of C consists only of
			// the current case type T (for simplicity, assume T is
			// not a type constraint interface for a type parameter),
			// then T is ineffectual. (The type parameter P can
			// be instantiated only with type argument T.)
			return false
		}
		// Note: if the type in a case is generic, it
		// will always be an instantiated generic type
		// in a valid Go program.
		return types.Implements(currentCaseType, e.intf)
	}

	var earlierIntfs []intfcase // interface types in earlier cases
	var currentIntfs []intfcase // interface types in current case

	for _, cc := range sw.Body.List {
		currentIntfs = currentIntfs[:0]
		for _, caseExpr := range cc.(*ast.CaseClause).List {
			caseType := pass.TypesInfo.TypeOf(caseExpr)
			if caseType == nil {
				continue
			}
			i, typeParam, ok := interfaceTypeOf(caseType)
			if ok {
				currentIntfs = append(currentIntfs, intfcase{i, typeParam, caseExpr})
			}
			for _, e := range earlierIntfs {
				// Assuming that the notation "Y -> X" means that type Y
				// in a case is ineffectual due to type X in an earlier case,
				// note that the !reported[e.expr] check below
				// prevents reporting a diagnostic for "C -> B"
				// in the following scenario.
				//
				//    B -> A; and
				//    C -> A; and
				//    C -> B
				//
				// The aim is to avoid confusing diagnostic messages
				// that refer to types that are already ineffectual.
				if ineffectualOrder(e, caseType) && !reported[e.expr] {
					var message string
					if testingDiagnosticShort {
						message = "ineffectual: %s -> %s"
					} else {
						message = "order of cases makes it impossible to match type %s in the case (due to type %s in earlier case)"
					}
					pass.Reportf(caseExpr.Pos(), message, formatExpr(caseExpr), formatExpr(e.expr))
					reported[caseExpr] = true
				}
			}
		}
		earlierIntfs = append(earlierIntfs, currentIntfs...)
	}
}

func underlyingIs[T types.Type](t types.Type) bool {
	_, ok := t.Underlying().(T)
	return ok
}

// hasTypeName determines whether t is a type with a name.
func hasTypeName(t types.Type) (*types.TypeName, bool) {
	switch t := t.(type) {
	case *types.Alias:
		return t.Obj(), true
	case *types.Named: // Note: also represents predeclared type error
		return t.Obj(), true
	case *types.TypeParam:
		return t.Obj(), true
	default:
		return nil, false
	}
}

func interfaceTypeOf(t types.Type) (_ *types.Interface, typeParam bool, ok bool) {
	t = types.Unalias(t)
	switch t := t.(type) {
	case *types.TypeParam:
		return t.Constraint().Underlying().(*types.Interface), true, true
	default:
		if v, ok := t.Underlying().(*types.Interface); ok {
			return v, false, true
		}
		return nil, false, false
	}
}

// Returns the defined type, if any, denoted by the alias type a.
// At most maxpointer pointer types may be present in the alias
// chain until a defined type is reached.
func aliasDenotesDefinedType(a *types.Alias, maxpointer int) (success bool, defined *types.Named, pointers int) {
	var t types.Type = a
	npointer := 0
	for {
		switch s := t.(type) {
		case *types.Alias:
			t = s.Rhs()
		case *types.Pointer:
			npointer += 1
			if npointer > maxpointer {
				return false, nil, 0
			}
			t = s.Elem()
		case *types.Named:
			return true, s, npointer
		default:
			return false, nil, 0
		}
	}
}

// computeImplements computes the set of non-interface defined
// types (in the set of predeclared types and r.TypeDecls)
// that implement the given interface type.
//
// If the defined type is a generic type, note that the
// "implements" relationship can have meaning only for an
// instantiated generic type.
// computeImplements currently does not support generic types.
//
// The interface type must be the type of either a type assertion
// expression or a type switch expression. As of
// go1.27, this means that the interface type must not be a
// general interface, i.e., an interface that cannot be
// used outside a type constraint.
// In other words, the above requirement implies that the
// interface type must be fully described by its method set.
//
// Additionally, the interface must not be the empty interface.
//
// [general interface]: https://go.dev/ref/spec#General_interfaces
func computeImplements(ch *typeSwitchState, r *finder.Result, intf *types.Interface) {
	if _, ok := ch.impl[intf]; ok {
		return
	}

	if !intf.IsMethodSet() {
		panic("interface not fully described by method set")
	}

	if intf.Empty() {
		panic("empty interface")
	}

	if ch.impl == nil {
		ch.impl = make(map[*types.Interface]map[implementer]struct{})
	}
	ch.impl[intf] = make(map[implementer]struct{})

	add := func(d implementer) {
		ch.impl[intf][d] = struct{}{}
	}

	maybeAdd := func(n *types.Named, pointer bool) {
		if underlyingIs[*types.Interface](n) {
			return
		}
		if n.TypeParams() != nil {
			return // generic types are not supported
		}
		var n2 types.Type
		if pointer {
			n2 = types.NewPointer(n)
		} else {
			n2 = n
		}
		// Note: types.Implements is documented to have
		// unspecified behavior if its V argument (the
		// first argument) is an uninstantiated generic type.
		if !types.Implements(n2, intf) {
			return
		}
		add(implementer{n, pointer})
	}

	// Go spec, as of go1.27, section 'Method declarations':
	// "[The receiver's] type must be a defined type T or a
	// pointer to a defined type T" and "T is called the
	// receiver base type. A receiver base type cannot be a
	// pointer or interface type".
	//
	// (continued) "If the receiver base type is a generic
	// type, the receiver specification must declare
	// corresponding type parameters for the method to use.
	// ... one for each type parameter of the receiver base type."
	//
	// (continued) "If the receiver type is denoted by (a
	// pointer to) an alias, the alias must not be generic
	// and it must not denote an instantiated generic type,
	// neither directly nor indirectly via another alias, and
	// irrespective of pointer indirections."
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
	// See also: Go spec section 'Method sets'.

	// Given the above descriptions from the Go spec and
	// given the current restrictions for the input interface
	// argument (method set only, non-empty) to this
	// function, it is sufficient to evaluate types that are
	// valid method receiver types.
	//
	// Thus the following code considers type T or type *T,
	// where T is a defined type.
	//
	// None of the non-interface predeclared types, as of
	// go1.27, have methods associated with them. Therefore
	// it is not necessary to evaluate the predeclared types.
	for tn := range r.TypeDecls {
		switch t := types.Unalias(tn.Type()).(type) {
		case *types.Named: // defined type
			n := t
			maybeAdd(n, true)
			maybeAdd(n, false)
		case *types.Pointer: // pointer to a defined type
			if n, ok := types.Unalias(t.Elem()).(*types.Named); ok {
				maybeAdd(n, true)
			}
		}
	}
}

func computeTypeRelations(ch *typeSwitchState, r *finder.Result) {
	if ch.identical != nil {
		return
	}
	ch.identical = make(map[tp][]tp, len(r.TypeDecls))

	add := func(t0 tp, t1 tp) {
		ch.identical[t0] = append(ch.identical[t0], t1)
	}

	for tn := range r.TypeDecls {
		switch t := tn.Type().(type) {
		case *types.Alias:
			success, defined, pointers := aliasDenotesDefinedType(t, 1)
			if !success {
				continue
			}
			if underlyingIs[*types.Interface](defined) {
				continue
			}
			if defined.TypeParams() != nil {
				continue
			}
			if pointers == 0 {
				add(tp{defined.Obj(), false}, tp{tn, false})
				add(tp{defined.Obj(), true}, tp{tn, true})
			} else { // pointers == 1
				add(tp{defined.Obj(), true}, tp{tn, false})
			}
		case *types.Named:
			n := t
			if underlyingIs[*types.Interface](n) {
				continue
			}
			if n.TypeParams() != nil {
				continue
			}
			add(tp{n.Obj(), false}, tp{tn, false})
			add(tp{n.Obj(), true}, tp{tn, true})
		default:
			panic(fmt.Sprintf("internal error: unexpected type %T", t))
		}
	}
}

// isEnumeratedType reports whether the given type (which
// typically is the type for an arbitrary expression) is an
// enumerated type. It returns the enumerated type and the
// enumerated constants of that type.
func isEnumeratedType(enums enumerated.Result, t types.Type) (*types.Named, []*types.Const, bool) {
	for {
		switch s := t.(type) {
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
			t = s.Rhs()
		case *types.Named:
			v, ok := enums[s]
			if ok {
				return s, v, true
			}
			// Package go/types documentation for (*types.Named).Underlying:
			// "Underlying types are never Named, TypeParam,
			// or Alias types."
			//
			// If this were not the case, then for the
			// following declarations, for the defined
			// type B, the value of "B".Underlying()
			// could be the defined type A, not int.
			//
			// 	type A int
			// 	type B A
			//
			// The guarantee that (*types.Named).Underlying() will
			// not be a *types.Named value is essential for
			// correctness of the current function. The goal
			// of the function is to determine whether the
			// type (say, B) represented by the argument t is
			// an enumerated type. The function must not consider
			// whether a different defined type (i.e., A) is
			// an enumerated type.
			if _, ok := s.Underlying().(*types.Named); ok {
				panic("(*types.Named).Underlying() unexpectedly a *types.Named type")
			}
			t = s.Underlying()
		default:
			// The possible types, as of go1.27, are handled
			// in the cases above and the default case should
			// not run. Nevertheless returning false is the appropriate
			// result for any unhandled type.
			return nil, nil, false
		}
	}
}

// The result is the set of constant values that must be
// satisfied in an expression switch.
// The result is prepared additionally in the context of the
// current package and the given options.
//
// The keys of the returned map are constant values in
// (go/constant.Value).ExactString representation.
// Note that it is possible for the same constant value to have
// multiple names, i.e. multiple *types.Const.
func needValues(currpkg *types.Package, cs []*types.Const, opts *options) map[string][]*types.Const {
	ret := make(map[string][]*types.Const, len(cs)) // usually one name per value
	for _, c := range cs {
		if !canAccessObject(currpkg, c) {
			continue
		}
		if isPkgLevel(c) && matchAny(opts.excludeConstPatterns, fullname(c)) {
			continue
		}
		ret[c.Val().ExactString()] = append(ret[c.Val().ExactString()], c)
	}
	return ret
}

// The result is the set of types that must be satisfied in a
// type switch. The result corresponds to the
// term "knowable type set" in the package comment.
func needTypes(currpkg *types.Package, swInterface *types.Interface, ch *typeSwitchState) map[implementer]struct{} {
	ret := make(map[implementer]struct{})
	for t := range ch.impl[swInterface] {
		for _, tt := range ch.identical[tp{t.named.Obj(), t.pointer}] {
			if !canAccessObject(currpkg, tt.typename) {
				continue
			}
			ret[t] = struct{}{}
			break
		}
	}
	return ret
}

// formats constant names for a "missing cases" diagnostic for an
// expression switch.
func formatMissingConstants(currpkg *types.Package, missing map[string][]*types.Const) string {
	// Note that the constants in the input are grouped by value.
	// The logic below does the following: The constants in
	// each group are sorted by AST position. The groups
	// themselves are then sorted by AST position (of the
	// first constant in the group).
	groups := make([][]*types.Const, 0, len(missing))
	for _, cs := range missing {
		cs := slices.Clone(cs) // do not modify input during sort
		slices.SortFunc(cs, func(a, b *types.Const) int { return cmp.Compare(a.Pos(), b.Pos()) })
		groups = append(groups, cs)
	}
	slices.SortFunc(groups, func(a, b []*types.Const) int { return cmp.Compare(a[0].Pos(), b[0].Pos()) })

	// format to string
	var buf strings.Builder
	for i, g := range groups {
		for i, c := range g {
			buf.WriteString(formatObject(currpkg, c))
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

// formats type names for a "missing cases" diagnostic for a type switch.
func formatMissingTypes(currpkg *types.Package, missing map[implementer][]tp, missingNil bool) string {
	makeGroups := func(missing map[implementer][]tp) [][]tp {
		cmp0 := func(a, b *types.TypeName) int {
			// aliases should be sorted after defined types.
			switch {
			case !a.IsAlias() && b.IsAlias():
				return -1
			case a.IsAlias() && !b.IsAlias():
				return 1
			}
			// predeclared types should be sorted first,
			// then types in the current package,
			// then types in external packages ordered by package path.
			switch {
			case a.Pkg() == nil && b.Pkg() != nil:
				return -1
			case a.Pkg() != nil && b.Pkg() == nil:
				return 1
			case a.Pkg() == currpkg && b.Pkg() != currpkg:
				return -1
			case a.Pkg() != currpkg && b.Pkg() == currpkg:
				return 1
			case a.Pkg() != b.Pkg():
				return cmp.Compare(a.Pkg().Path(), b.Pkg().Path())
			}
			// types in the same package should be ordered by name.
			if a.Name() != b.Name() {
				return cmp.Compare(a.Name(), b.Name())
			}
			return 0
		}
		cmp1 := func(a, b tp) int {
			if v := cmp0(a.typename, b.typename); v != 0 {
				return v
			}
			// all else being the same (above),
			// a non-pointer type should be sorted first.
			switch {
			case !a.pointer && b.pointer:
				return -1
			case a.pointer && !b.pointer:
				return 1
			}
			return 0
		}
		cmp2 := func(a, b implementer) int {
			if v := cmp0(a.named.Obj(), b.named.Obj()); v != 0 {
				return v
			}
			switch {
			case !a.pointer && b.pointer:
				return -1
			case a.pointer && !b.pointer:
				return 1
			}
			return 0
		}

		s := slices.SortedFunc(maps.Keys(missing), cmp2)
		groups := make([][]tp, 0, len(missing))
		for _, t := range s {
			g := missing[t]
			g = slices.Clone(g) // do not modify input during sort
			slices.SortFunc(g, cmp1)
			groups = append(groups, g)
		}
		slices.SortFunc(groups, func(g1, g2 []tp) int { return cmp1(g1[0], g2[0]) })
		return groups
	}

	// helper function that formats a single type.
	formatType := func(tt tp) string {
		var buf bytes.Buffer
		if tt.pointer {
			buf.WriteByte('*')
		}
		buf.WriteString(formatObject(currpkg, tt.typename))

		switch s := tt.typename.Type().(type) {
		case *types.Named:
			// Generic defined types are not supported and
			// should not have been produced earlier.
		case *types.Alias:
			// The alias type may be generic but the
			// type arguments used to instantiate it
			// should be inconsequential, since the
			// defined type that the alias denotes
			// will not be a generic type.
			if s.TypeParams() != nil {
				buf.WriteString("[…]")
			}
		default:
			panic(fmt.Sprintf("internal error: unexpected type %T", tt))
		}

		return buf.String()
	}

	groups := makeGroups(missing)

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
		for _, tt := range g {
			if wrote {
				buf.WriteString(" or ")
			}
			buf.WriteString(formatType(tt))
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

// the main function that checks a map literal.
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
	_, cs, ok := isEnumeratedType(enums, mapType.Key())
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
		pass.Reportf(compLit.Pos(), "missing keys in map literal: %s", formatMissingConstants(pass.Pkg, need))
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

func formatObject(currpkg *types.Package, obj types.Object) string {
	if obj.Pkg() == nil || currpkg == obj.Pkg() {
		return obj.Name()
	}
	if !isPkgLevel(obj) {
		panic("not a package-level object")
	}
	return obj.Pkg().Name() + "." + obj.Name()
}

// fullname returns the object name qualified with the package
// path (always). The result is meaningful only for package-level
// objects.
func fullname(obj types.Object) string {
	if !isPkgLevel(obj) {
		panic("not a package-level object")
	}
	var b strings.Builder
	b.WriteString(obj.Pkg().Path())
	b.WriteString(".")
	b.WriteString(obj.Name())
	return b.String()
}

func canAccessObject(current *types.Package, obj types.Object) bool {
	other := obj.Pkg()
	if other == nil {
		return true // object in universe scope
	}
	if other == current {
		return true // same package
	}
	// Go spec section 'Exported identifiers': "An identifier
	// may be exported to permit access to it from another
	// package. An identifier is exported if both:
	// * the first character of the identifier's name is a
	// Unicode uppercase letter (Unicode character category
	// Lu); and
	// * the identifier is declared in the package block or
	// it is a field name or method name."
	//
	// Note: The "field name or method name" part of
	// the quoted spec is not handled here. However that part
	// is not applicable to current callers.
	return obj.Exported() && isPkgLevel(obj)
}

func isPkgLevel(obj types.Object) bool {
	return obj.Pkg() != nil && obj.Pkg().Scope().Lookup(obj.Name()) == obj
}

func formatExpr(e ast.Expr) string {
	return types.ExprString(e)
}
