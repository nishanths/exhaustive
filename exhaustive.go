/*
Package exhaustive defines an analyzer that checks exhaustiveness of
enum switch statements in Go source code. It can be configured to
additionally check exhaustiveness of map literals that have enum key
types.

# Definition of enum types and enum members

The Go language spec does not provide an explicit definition for an
enum. By convention, and for the purpose of this analyzer, an enum type
is any named type that meets these requirements:

 1. has underlying type float, string, or integer (includes byte and
    rune); and
 2. has at least one constant of its type defined in the same scope.

In the example below, Biome is an enum type. The 3 constants are its
enum members.

	package eco

	type Biome int

	const (
	    Tundra  Biome = 1
	    Savanna Biome = 2
	    Desert  Biome = 3
	)

Enum member constants for a particular enum type do not necessarily all
have to be declared in the same const block. The constant values may be
specified using iota, using literal values, or using any valid means for
declaring a Go constant. It is allowed for multiple enum member
constants for a particular enum type to have the same constant value.

# Definition of exhaustiveness

A switch statement that switches on a value of an enum type is
exhaustive if all of the enum members are listed in the switch
statement's cases. If multiple enum members have the same constant
value, it is sufficient for any one of these same-valued members to be
listed.

For an enum type defined in the same package as the switch statement,
both exported and unexported enum members must be listed to satisfy
exhaustiveness. For an enum type defined in an external package, it is
sufficient that only exported enum members are listed.

Only identifiers (e.g. Tundra) and qualilified identifiers (e.g.
somepkg.Grassland) that name constants listed in a switch statement's
cases may contribute towards satisfying exhaustiveness; literal values
or variables will not.

When using the default analyzer configuration, the existence of a
'default' case in a switch statement, on its own, does not automatically
make a switch statement exhaustive. See the
-default-signifies-exhaustive flag to adjust this behavior.

A similar definition of exhaustiveness applies to a map literal whose
key type is an enum type. To be exhaustive, the map literal must specify
keys corresponding to all of the enum members. Empty map literals are never
checked. Note that the -check flag must include "map" for map literals
to be checked.

# Type aliases

The analyzer handles type aliases in the following manner. In the
example, T2 is a enum type, and T1 is an alias for T2.  Note that we
don't call T1 itself an enum type; T1 is only an alias for an enum type.

	package pkg
	type T1 = newpkg.T2
	const (
	    A = newpkg.A
	    B = newpkg.B
	)

	package newpkg
	type T2 int
	const (
	    A T2 = 1
	    B T2 = 2
	)

A switch statement that switches on a value of type T1 (which, in
reality, is just an alternate spelling for type T2) is exhaustive if all
of T2's enum members are listed in the switch statement's cases. Recall
that only constants declared in the same scope as type T2's scope
can be T2's enum members.

The analyzer guarantees that introducing a type alias (such as type T1 =
newpkg.T2) will never result in new diagnostics from the analyzer, as
long as the set of enum member constant values of the RHS type is a
subset of the set of enum member constant values of the old LHS type.

The following switch statements are equally valid and exhaustive.

	// The type of v is effectively newpkg.T2 due to alias.
	var v pkg.T1

	// pkg.A is a valid substitute for newpkg.A (same constant value).
	// Similarly for pkg.B.
	switch v {
	case pkg.A:
	case pkg.B:
	}

	switch v {
	case newpkg.A:
	case newpkg.B:
	}

# Type parameters

A switch statement that switches on a value of a type-parameterized type
is checked for exhaustiveness iff each of the elements of its constraint
is an enum type. The following switch statement will be checked,
assuming M, N, and O are enum types. To satisfy exhaustiveness, all
members for each of M, N, and O must be listed in the switch statement's
cases.

	func bar[T M | N | O](v T) {
		switch v {
		}
	}

# Flags

Flags supported by the analyzer are described below. All flags are
optional.

	flag                            type    default value
	----                            ----    -------------
	-check                          string  switch
	-explicit-exhaustive-switch     bool    false
	-explicit-exhaustive-map        bool    false
	-check-generated                bool    false
	-default-signifies-exhaustive   bool    false
	-ignore-enum-members            string  (empty)
	-package-scope-only             bool    false

The -check flag specifies is a comma-separated list of program elements
that should be checked for exhaustiveness. Supported program elements
are "switch" and "map". By default, only switch statements are checked.
Specify -check=switch,map to additionally check map literals.

If the -explicit-exhaustive-switch flag is enabled, the analyzer only
checks enum switch statements associated with a comment beginning with
"//exhaustive:enforce". By default the flag is disabled, which means
that the analyzer checks every enum switch statement not associated with
a comment beginning with "//exhaustive:ignore". The
-explicit-exhaustive-map flag is the map literal counterpart of the
-explicit-exhaustive-switch flag.

	//exhaustive:ignore
	switch v {
	case A:
	case B:
	}

If the -check-generated flag is enabled, switch statements or map
literals in generated Go source files are also checked. Otherwise, by
default, generated files are not checked. Refer to
https://golang.org/s/generatedcode for the definition of generated
files.

If the -default-signifies-exhaustive flag is enabled, the presence of a
'default' case in a switch statement always satisfies exhaustiveness,
even if all enum members are not listed. It is recommended that you do
not enable this flag. Enabling it usually defeats the purpose of
exhaustiveness checking.

The -ignore-enum-members flag specifies a regular expression in Go
package regexp syntax. Enum members matching the regular expression
do not have to be listed in switch statement cases to satisfy
exhaustiveness. The specified regular expression is matched against an
enum member name inclusive of the enum package import path: for example,
if the enum package import path is "example.com/eco" and the member name
is "Tundra", the specified regular expression will be matched against
the string "example.com/eco.Tundra".

If the -package-scope-only flag is enabled, the analyzer only finds
enums defined in package-level scopes, and consequently only switch
statements and map literals that use package-level enums will be checked
for exhaustiveness. By default, the analyzer finds enums defined in all
scopes, and checks switch statements that switch on all these enums.
*/
package exhaustive

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func init() {
	Analyzer.Flags.Var(&fCheck, CheckFlag, "comma-separated list of program elements to check for exhaustiveness; supported elements are: switch, map")
	Analyzer.Flags.BoolVar(&fExplicitExhaustiveSwitch, ExplicitExhaustiveSwitchFlag, false, `run exhaustive check on switch statements with "//exhaustive:enforce" comment`)
	Analyzer.Flags.BoolVar(&fExplicitExhaustiveMap, ExplicitExhaustiveMapFlag, false, `run exhaustive check on map literals with "//exhaustive:enforce" comment`)
	Analyzer.Flags.BoolVar(&fCheckGenerated, CheckGeneratedFlag, false, "check generated files")
	Analyzer.Flags.BoolVar(&fDefaultSignifiesExhaustive, DefaultSignifiesExhaustiveFlag, false, `presence of "default" case in a switch statement unconditionally satisfies exhaustiveness`)
	Analyzer.Flags.Var(&fIgnoreEnumMembers, IgnoreEnumMembersFlag, "enum members matching `regexp` do not have to be listed to satisfy exhaustiveness")
	Analyzer.Flags.BoolVar(&fPackageScopeOnly, PackageScopeOnlyFlag, false, "find enums only in package-level scopes, not in inner scopes")

	var unused string
	Analyzer.Flags.StringVar(&unused, IgnorePatternFlag, "", "no effect (deprecated); see -"+IgnoreEnumMembersFlag+" instead")
	Analyzer.Flags.StringVar(&unused, CheckingStrategyFlag, "", "no effect (deprecated)")
}

// Flag names used by the analyzer. They are exported for use by analyzer
// driver programs.
const (
	CheckFlag                      = "check"
	ExplicitExhaustiveSwitchFlag   = "explicit-exhaustive-switch"
	ExplicitExhaustiveMapFlag      = "explicit-exhaustive-map"
	CheckGeneratedFlag             = "check-generated"
	DefaultSignifiesExhaustiveFlag = "default-signifies-exhaustive"
	IgnoreEnumMembersFlag          = "ignore-enum-members"
	PackageScopeOnlyFlag           = "package-scope-only"

	IgnorePatternFlag    = "ignore-pattern"    // Deprecated: see IgnoreEnumMembersFlag instead.
	CheckingStrategyFlag = "checking-strategy" // Deprecated.
)

// checkElement is a program element supported by the -check flag.
type checkElement string

const (
	elementSwitch checkElement = "switch"
	elementMap    checkElement = "map"
)

func validCheckElement(s string) error {
	e := checkElement(s) // temporarily for check
	switch e {
	case elementSwitch:
		return nil
	case elementMap:
		return nil
	default:
		return fmt.Errorf("invalid program element %q", s)
	}
}

var defaultCheckElements = []string{
	string(elementSwitch),
}

// Flag values.
var (
	fCheck                      = stringsFlag{elements: defaultCheckElements, filter: validCheckElement}
	fExplicitExhaustiveSwitch   bool
	fExplicitExhaustiveMap      bool
	fCheckGenerated             bool
	fDefaultSignifiesExhaustive bool
	fIgnoreEnumMembers          regexpFlag
	fPackageScopeOnly           bool
)

// resetFlags resets the flag variables to their default values.
// Useful in tests.
func resetFlags() {
	fCheck = stringsFlag{elements: defaultCheckElements, filter: validCheckElement}
	fExplicitExhaustiveSwitch = false
	fExplicitExhaustiveMap = false
	fCheckGenerated = false
	fDefaultSignifiesExhaustive = false
	fIgnoreEnumMembers = regexpFlag{}
	fPackageScopeOnly = false
}

var Analyzer = &analysis.Analyzer{
	Name:      "exhaustive",
	Doc:       "check exhaustiveness of enum switch statements",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&enumMembersFact{}},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	for typ, members := range findEnums(
		fPackageScopeOnly,
		pass.Pkg,
		inspect,
		pass.TypesInfo,
	) {
		exportFact(pass, typ, members)
	}

	generated := boolCache{value: isGeneratedFile}
	comments := commentCache{value: fileCommentMap}
	swConf := switchConfig{
		explicit:                   fExplicitExhaustiveSwitch,
		defaultSignifiesExhaustive: fDefaultSignifiesExhaustive,
		checkGenerated:             fCheckGenerated,
		ignoreEnumMembers:          fIgnoreEnumMembers.regexp(),
	}
	mapConf := mapConfig{
		explicit:          fExplicitExhaustiveMap,
		checkGenerated:    fCheckGenerated,
		ignoreEnumMembers: fIgnoreEnumMembers.regexp(),
	}
	swChecker := switchChecker(pass, swConf, generated, comments)
	mapChecker := mapChecker(pass, mapConf, generated, comments)

	for _, e := range fCheck.elements {
		switch checkElement(e) {
		case elementSwitch:
			inspect.WithStack([]ast.Node{&ast.SwitchStmt{}}, toVisitor(swChecker))
		case elementMap:
			inspect.WithStack([]ast.Node{&ast.CompositeLit{}}, toVisitor(mapChecker))
		default:
			panic(fmt.Sprintf("unknown checkElement %v", e))
		}
	}
	return nil, nil
}

// toVisitor converts a nodeVisitor to a function suitable for use
// with inspect.WithStack.
func toVisitor(v nodeVisitor) func(ast.Node, bool, []ast.Node) bool {
	return func(node ast.Node, push bool, stack []ast.Node) bool {
		proceed, _ := v(node, push, stack)
		return proceed
	}
}

// TODO(nishanths): When dropping pre go1.19 support, the following
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
