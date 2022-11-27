/*
Package exhaustive defines an analyzer that checks exhaustiveness of switch
statements of enum-like constants in Go source code. The analyzer can be
configured to additionally check exhaustiveness of map literals whose key type
is enum-like.

# Definition of enum

The Go [language spec] does not provide an explicit definition for enums. For
the purpose of this analyzer, and by convention, an enum type is any named
type that:

  - has underlying type float, string, or integer (includes byte and
    rune, which are aliases for uint8 and int32, respectively); and
  - has at least one constant of the type defined in the same scope.

In the example below, Biome is an enum type. The three constants are its
enum members.

	package eco

	type Biome int

	const (
		Tundra Biome = iota
		Savanna
		Desert
	)

Enum member constants for a particular enum type do not necessarily all
have to be declared in the same const block. The constant values may be
specified using iota, using literal values, or using any valid means for
declaring a Go constant. It is allowed for multiple enum member
constants for a particular enum type to have the same constant value.

# Definition of exhaustiveness

A switch statement that switches on a value of an enum type is exhaustive if
all of the enum members are listed in the switch statement's cases. If
multiple members have the same constant value, it is sufficient for any one of
these same-valued members to be listed.

For an enum type defined in the same package as the switch statement, both
exported and unexported enum members must be listed to satisfy exhaustiveness.
For an enum type defined in an external package, it is sufficient that only
exported enum members are listed. Only identifiers (e.g. Tundra) and qualified
identifiers (e.g. somepkg.Grassland) that name constants may contribute
towards satisfying exhaustiveness; other expressions such as literal values
and function calls will not.

When using the default analyzer configuration, the existence of a
default case in a switch statement, on its own, does not immediately
make a switch statement exhaustive. See the
-default-signifies-exhaustive flag to adjust this behavior.

A similar definition of exhaustiveness applies to a map literal whose key type
is an enum type. To be exhaustive, all of the enum members must be listed in
the map literal's keys. Empty map literals will not be checked. Note that the
-check flag must include "map" for map literals to be checked.

# Type aliases

The analyzer handles type aliases as shown in the following example. T2
is a enum type, and T1 is an alias for T2. Note that we don't call T1
itself an enum type; T1 is only an alias for an enum type.

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

A switch statement that switches on a value of type T1 (which, in reality, is
just an alternate spelling for type T2) is exhaustive if all of T2's enum
members are listed in the switch statement's cases. (Recall that
only constants declared in the same scope as type T2's scope can be T2's enum
members.) The following switch statements are valid Go code and and are
exhaustive.

	// Note that the type of v is effectively newpkg.T2 due to alias.
	func f(v pkg.T1) {
		switch v {
		case newpkg.A:
		case newpkg.B:
		}
		switch v {
		case pkg.A:
		case pkg.B:
		}
	}

The analyzer guarantees that introducing a type alias (such as type T1 =
newpkg.T2) will not result in new diagnostics from the analyzer, as long as
the set of enum member constant values of the alias RHS type is a subset of
the set of enum member constant values of the LHS type.

# Type parameters

A switch statement that switches on a value whose type is a type parameter is
checked for exhaustiveness iff each type element in the type constraint is an
enum type and shares the same underlying basic kind (e.g. uint8, string). In
the following example, the switch statement will be checked, provided M, N,
and O are enum types with the same underlying basic kind. To satisfy
exhaustiveness, all members for each of the types M, N, and O must be listed
in the switch statement's cases.

	func bar[T M | I](v T) {
		switch v {
			...
		}
	}
	type I interface{ N | J }
	type J interface{ O }

# Flags

Flags used by the analyzer are described below.

	flag                           type                     default value
	----                           ----                     -------------
	-check                         comma-separated strings  switch
	-explicit-exhaustive-switch    bool                     false
	-explicit-exhaustive-map       bool                     false
	-check-generated               bool                     false
	-default-signifies-exhaustive  bool                     false
	-ignore-enum-members           regexp pattern           (none)
	-ignore-enum-types             regexp pattern           (none)
	-package-scope-only            bool                     false

The -check flag specifies is a comma-separated list of program elements
that should be checked for exhaustiveness. Supported program elements
are "switch" and "map". By default, only switch statements are checked.
Specify -check=switch,map to additionally check map literals.

If the -explicit-exhaustive-switch flag is enabled, the analyzer checks a
switch statement only if it associated with a comment beginning with
"//exhaustive:enforce". By default the flag is disabled, which means that the
analyzer checks every enum switch statement not associated with a comment
beginning with "//exhaustive:ignore".

The -explicit-exhaustive-map flag is the map literal counterpart of the
-explicit-exhaustive-switch flag.

If the -check-generated flag is enabled, switch statements and map
literals in generated Go source files are checked. Otherwise, by
default, generated files are ignored. Refer to
https://golang.org/s/generatedcode for the definition of generated
files.

If the -default-signifies-exhaustive flag is enabled, the presence of a
default case in a switch statement unconditionally satisfies exhaustiveness
(all enum members do not have to be listed). Enabling this flag usually tends
to counter the purpose of exhaustiveness checking, so it is not recommended
that you do so.

The -ignore-enum-members flag specifies a regular expression in Go package
regexp syntax. Constants matching the regular expression do not have to be
listed in switch statement cases or map literals in order to satisfy
exhaustiveness. The specified regular expression is matched against the
constant name inclusive of the enum package import path. For example, if the
package import path of the constant is "example.com/eco" and the constant name
is "Tundra", the specified regular expression will be matched against the
string "example.com/eco.Tundra".

The -ignore-enum-types flag is similar to the -ignore-enum-members flag,
except that it applies to types.

If the -package-scope-only flag is enabled, the analyzer only finds enums
defined in package scope, but not in inner scopes such as functions.
Consequently only switch statements and map literals that use these enums will
be checked for exhaustiveness. By default, the analyzer finds enums defined in
all scopes, including in inner scopes such as functions.

# Skip analysis

To skip analysis of a switch statement or a map literal, associate it with a
comment that begins with "//exhaustive:ignore". For example:

	//exhaustive:ignore
	switch v {
	case A:
	case B:
	}

To ignore specific constants in exhaustiveness checks, use the
-ignore-enum-members flag. Similarly, to ignore specific types, use the
-ignore-enum-types flag. For example:

	exhaustive -ignore-enum-types '^time\.Duration$|^example.org/measure\.Unit$'

[language spec]: https://golang.org/ref/spec
*/
package exhaustive

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func init() {
	Analyzer.Flags.Var(&fCheck, CheckFlag, "comma-separated list of program `elements` that should be checked for exhaustiveness; supported elements are: switch, map")
	Analyzer.Flags.BoolVar(&fExplicitExhaustiveSwitch, ExplicitExhaustiveSwitchFlag, false, `check switch statement only if associated with "//exhaustive:enforce" comment`)
	Analyzer.Flags.BoolVar(&fExplicitExhaustiveMap, ExplicitExhaustiveMapFlag, false, `check map literal only if associated with "//exhaustive:enforce" comment`)
	Analyzer.Flags.BoolVar(&fCheckGenerated, CheckGeneratedFlag, false, "check generated files")
	Analyzer.Flags.BoolVar(&fDefaultSignifiesExhaustive, DefaultSignifiesExhaustiveFlag, false, "presence of default case in switch statement unconditionally satisfies exhaustiveness")
	Analyzer.Flags.Var(&fIgnoreEnumMembers, IgnoreEnumMembersFlag, "constants matching `regexp` are ignored for exhaustiveness checks")
	Analyzer.Flags.Var(&fIgnoreEnumTypes, IgnoreEnumTypesFlag, "types matching `regexp` are ignored for exhaustiveness checks")
	Analyzer.Flags.BoolVar(&fPackageScopeOnly, PackageScopeOnlyFlag, false, "find enums only in package scopes, not inner scopes")

	var unused string
	Analyzer.Flags.StringVar(&unused, IgnorePatternFlag, "", "no effect (deprecated); use -"+IgnoreEnumMembersFlag)
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
	IgnoreEnumTypesFlag            = "ignore-enum-types"
	PackageScopeOnlyFlag           = "package-scope-only"

	IgnorePatternFlag    = "ignore-pattern"    // Deprecated: use IgnoreEnumMembersFlag.
	CheckingStrategyFlag = "checking-strategy" // Deprecated.
)

// checkElement is a program element supported by the -check flag.
type checkElement string

const (
	elementSwitch checkElement = "switch"
	elementMap    checkElement = "map"
)

func validCheckElement(s string) error {
	switch checkElement(s) {
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
	fIgnoreEnumTypes            regexpFlag
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
	fIgnoreEnumTypes = regexpFlag{}
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
		ignoreConstant:             fIgnoreEnumMembers.re,
		ignoreType:                 fIgnoreEnumTypes.re,
	}
	mapConf := mapConfig{
		explicit:       fExplicitExhaustiveMap,
		checkGenerated: fCheckGenerated,
		ignoreConstant: fIgnoreEnumMembers.re,
		ignoreType:     fIgnoreEnumTypes.re,
	}
	swChecker := switchChecker(pass, swConf, generated, comments)
	mapChecker := mapChecker(pass, mapConf, generated, comments)

	// NOTE: should not share the same inspect.WithStack call for different
	// program elements: the visitor function for a program element may
	// exit traversal early, but this shouldn't affect traversal for
	// other program elements.
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
