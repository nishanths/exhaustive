package exhaustive

import (
	"flag"
	"go/ast"
	"go/token"
	"go/types"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func denotesPackage(ident *ast.Ident, info *types.Info) bool {
	obj := info.ObjectOf(ident)
	if obj == nil {
		return false
	}
	_, ok := obj.(*types.PkgName)
	return ok
}

// exprConstVal returns the constantValue for an expression if the
// expression is a constant value and if the expression is considered
// valid to satisfy exhaustiveness as defined by this program.
// Otherwise it returns (_, false).
func exprConstVal(e ast.Expr, info *types.Info) (constantValue, bool) {
	handleIdent := func(ident *ast.Ident) (constantValue, bool) {
		obj := info.Uses[ident]
		if obj == nil {
			return "", false
		}
		if _, ok := obj.(*types.Const); !ok {
			return "", false
		}
		// There are two scenarios.
		// See related test cases in typealias/quux/quux.go.
		//
		// Scenario 1
		//
		// Tag package and constant package are the same. This is
		// simple; we just use fs.ModeDir's value.
		//
		// Example:
		//
		//   var mode fs.FileMode
		//   switch mode {
		//   case fs.ModeDir:
		//   }
		//
		// Scenario 2
		//
		// Tag package and constant package are different. In this
		// scenario, too, we accept the case clause expr constant value,
		// as is. If the Go type checker is okay with the name being
		// listed in the case clause, we don't care much further.
		//
		// Example:
		//
		//   var mode fs.FileMode
		//   switch mode {
		//   case os.ModeDir:
		//   }
		//
		// Or equivalently:
		//
		//   // The type of mode is effectively fs.FileMode,
		//   // due to type alias.
		//   var mode os.FileMode
		//   switch mode {
		//   case os.ModeDir:
		//   }
		return determineConstVal(ident, info), true
	}

	e = stripTypeConversions(astutil.Unparen(e), info)

	switch e := e.(type) {
	case *ast.Ident:
		return handleIdent(e)

	case *ast.SelectorExpr:
		x := astutil.Unparen(e.X)
		// Ensure we only see the form pkg.Const, and not e.g.
		// structVal.f or structVal.inner.f.
		//
		// For this purpose, first we check that X, which is everything
		// except the rightmost field selector *ast.Ident (the Sel
		// field), is also an *ast.Ident.
		xIdent, ok := x.(*ast.Ident)
		if !ok {
			return "", false
		}
		// Second, check that it's a package. It doesn't matter which
		// package, just that it denotes some package.
		if !denotesPackage(xIdent, info) {
			return "", false
		}
		return handleIdent(e.Sel)

	default:
		// e.g. literal
		// these aren't considered towards satisfying exhaustiveness.
		return "", false
	}
}

func stripTypeConversions(e ast.Expr, info *types.Info) ast.Expr {
	c, ok := e.(*ast.CallExpr)
	if !ok {
		return e
	}
	if len(c.Args) != 1 {
		return e
	}
	typ := info.TypeOf(c.Fun)
	if typ == nil {
		// can never happen for a valid Go program?
		return e
	}
	// must not allow function calls.
	if _, ok := typ.Underlying().(*types.Signature); ok {
		return e
	}
	return stripTypeConversions(astutil.Unparen(c.Args[0]), info)
}

// member is a single member of an enum type.
type member struct {
	pos  token.Pos
	typ  enumType
	name string
	val  constantValue
}

// typeAndMembers combines an enumType and its members set.
type typeAndMembers struct {
	et enumType
	em enumMembers
}

var _ flag.Value = (*regexpFlag)(nil)
var _ flag.Value = (*stringsFlag)(nil)

// regexpFlag implements flag.Value for parsing
// regular expression flag inputs.
type regexpFlag struct{ rx *regexp.Regexp }

func (f *regexpFlag) String() string {
	if f == nil || f.rx == nil {
		return ""
	}
	return f.rx.String()
}

func (f *regexpFlag) Set(expr string) error {
	if expr == "" {
		f.rx = nil
		return nil
	}

	rx, err := regexp.Compile(expr)
	if err != nil {
		return err
	}

	f.rx = rx
	return nil
}

// stringsFlag implements flag.Value for parsing a comma-separated string
// list. Surrounding whitespace is stripped from the input and from each
// element. If filter is non-nil it is called for each element in the input.
type stringsFlag struct {
	elements []string
	filter   func(string) error
}

func (f *stringsFlag) String() string {
	if f == nil {
		return ""
	}
	return strings.Join(f.elements, ",")
}

func (f *stringsFlag) filterFunc() func(string) error {
	if f.filter != nil {
		return f.filter
	}
	return func(_ string) error { return nil }
}

func (f *stringsFlag) Set(input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		f.elements = nil
		return nil
	}

	for _, el := range strings.Split(input, ",") {
		el = strings.TrimSpace(el)
		if err := f.filterFunc()(el); err != nil {
			return err
		}
		f.elements = append(f.elements, el)
	}
	return nil
}

type checklist struct {
	info     map[enumType]enumMembers
	checkl   map[member]struct{}
	ignoreRx *regexp.Regexp
}

func (c *checklist) ignore(pattern *regexp.Regexp) {
	c.ignoreRx = pattern
}

func (c *checklist) add(et enumType, em enumMembers, includeUnexported bool) {
	addOne := func(name string) {
		if name == "_" {
			// Blank identifier is often used to skip entries in iota
			// lists.  Also, it can't be referenced anywhere (e.g. can't
			// be referenced in switch statement cases) It doesn't make
			// sense to include it as required member to satisfy
			// exhaustiveness.
			return
		}
		if !ast.IsExported(name) && !includeUnexported {
			return
		}
		if c.ignoreRx != nil && c.ignoreRx.MatchString(et.Pkg().Path()+"."+name) {
			return
		}
		mem := member{
			em.NameToPos[name],
			et,
			name,
			em.NameToValue[name],
		}
		if c.checkl == nil {
			c.checkl = make(map[member]struct{})
		}
		c.checkl[mem] = struct{}{}
	}

	if c.info == nil {
		c.info = make(map[enumType]enumMembers)
	}
	c.info[et] = em

	for _, name := range em.Names {
		addOne(name)
	}
}

func (c *checklist) found(val constantValue) {
	// delete all same-valued items.
	for et, em := range c.info {
		for _, name := range em.ValueToNames[val] {
			delete(c.checkl, member{
				em.NameToPos[name],
				et,
				name,
				em.NameToValue[name],
			})
		}
	}
}

func (c *checklist) remaining() map[member]struct{} {
	return c.checkl
}

// group is a collection of same-valued members, possibly from
// different enum types.
type group []member

func groupMissing(missing map[member]struct{}, types []enumType) []group {
	// indices maps each element in the input slice to its index.
	indices := func(vs []enumType) map[enumType]int {
		ret := make(map[enumType]int, len(vs))
		for i, v := range vs {
			ret[v] = i
		}
		return ret
	}

	typesOrder := indices(types) // for quick lookup
	astBefore := func(x, y member) bool {
		if typesOrder[x.typ] < typesOrder[y.typ] {
			return true
		}
		if typesOrder[x.typ] > typesOrder[y.typ] {
			return false
		}
		return x.pos < y.pos
	}

	// byConstVal groups member names by constant value.
	byConstVal := func(members map[member]struct{}) map[constantValue][]member {
		ret := make(map[constantValue][]member)
		for m := range members {
			ret[m.val] = append(ret[m.val], m)
		}
		return ret
	}

	var groups []group
	for _, members := range byConstVal(missing) {
		groups = append(groups, group(members))
	}

	// sort members within each group in AST order.
	for i := range groups {
		g := groups[i]
		sort.Slice(g, func(i, j int) bool { return astBefore(g[i], g[j]) })
		groups[i] = g
	}
	// sort groups themselves in AST order.
	// the index [0] access is safe, because there will be at least one
	// element per group.
	sort.Slice(groups, func(i, j int) bool { return astBefore(groups[i][0], groups[j][0]) })

	return groups
}

func diagnosticEnumType(enumType *types.TypeName) string {
	return enumType.Pkg().Name() + "." + enumType.Name()
}

func dedupEnumTypes(types []enumType) []enumType {
	// TODO(nishanths) this function is a candidate for type parameterization

	m := make(map[enumType]struct{})
	var ret []enumType
	for _, t := range types {
		_, ok := m[t]
		if ok {
			continue
		}
		m[t] = struct{}{}
		ret = append(ret, t)
	}
	return ret
}

func diagnosticEnumTypes(types []enumType) string {
	var buf strings.Builder
	for i := range types {
		buf.WriteString(diagnosticEnumType(types[i].TypeName))
		if i != len(types)-1 {
			buf.WriteByte('|')
		}
	}
	return buf.String()
}

func diagnosticMember(m member) string {
	return m.typ.Pkg().Name() + "." + m.name
}

func diagnosticGroups(gs []group) string {
	out := make([]string, len(gs))
	for i := range gs {
		var buf strings.Builder
		for j := range gs[i] {
			buf.WriteString(diagnosticMember(gs[i][j]))
			if j != len(gs[i])-1 {
				buf.WriteByte('|')
			}
		}
		out[i] = buf.String()
	}
	return strings.Join(out, ", ")
}

// TODO(nishanths) If dropping pre-go1.18 support, the following
// types and functions are candidates to be type parameterized.

type boolCache struct {
	m     map[*ast.File]bool
	value func(*ast.File) bool
}

func (c boolCache) get(file *ast.File) bool {
	if c.m == nil {
		c.m = make(map[*ast.File]bool)
	}
	if _, ok := c.m[file]; !ok {
		c.m[file] = c.value(file)
	}
	return c.m[file]
}

type commentCache struct {
	m     map[*ast.File]ast.CommentMap
	value func(*token.FileSet, *ast.File) ast.CommentMap
}

func (c commentCache) get(fset *token.FileSet, file *ast.File) ast.CommentMap {
	if c.m == nil {
		c.m = make(map[*ast.File]ast.CommentMap)
	}
	if _, ok := c.m[file]; !ok {
		c.m[file] = c.value(fset, file)
	}
	return c.m[file]
}
