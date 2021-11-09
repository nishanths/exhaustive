package exhaustive

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Flag names used by the analyzer. They are exported for use by analyzer
// driver programs.
const (
	CheckGeneratedFlag             = "check-generated"
	DefaultSignifiesExhaustiveFlag = "default-signifies-exhaustive"
	IgnoreEnumMembersFlag          = "ignore-enum-members"
	PackageScopeOnly               = "package-scope-only"

	excludeTypeAliasFlag = "exclude-type-alias"

	IgnorePatternFlag    = "ignore-pattern"    // Deprecated: see IgnoreEnumMembersFlag instead.
	CheckingStrategyFlag = "checking-strategy" // Deprecated.
)

var (
	// Public flags.
	fCheckGeneratedFiles        bool
	fDefaultSignifiesExhaustive bool
	fIgnoreEnumMembers          regexpFlag
	fPackageScopeOnly           bool

	// Internal flags.
	fExcludeTypeAlias bool
)

// resetFlags resets the flag variables to their default values.
// Useful in tests.
func resetFlags() {
	// Public flags.
	fCheckGeneratedFiles = false
	fDefaultSignifiesExhaustive = false
	fIgnoreEnumMembers = regexpFlag{}
	fPackageScopeOnly = false

	// Internal flags.
	fExcludeTypeAlias = false
}

func init() {
	var unused string

	// Public flags.
	Analyzer.Flags.BoolVar(&fCheckGeneratedFiles, CheckGeneratedFlag, false, "check switch statements in generated files")
	Analyzer.Flags.BoolVar(&fDefaultSignifiesExhaustive, DefaultSignifiesExhaustiveFlag, false, "presence of \"default\" case in switch statements satisfies exhaustiveness, even if all enum members are not listed")
	Analyzer.Flags.Var(&fIgnoreEnumMembers, IgnoreEnumMembersFlag, "enum members matching `regex` do not have to be listed in switch statements to satisfy exhaustiveness")
	Analyzer.Flags.BoolVar(&fPackageScopeOnly, PackageScopeOnly, false, "consider enums only in package scopes, not in inner scopes")

	// Internal flags.
	Analyzer.Flags.BoolVar(&fExcludeTypeAlias, excludeTypeAliasFlag, false, "don't check switch statements in which the switch tag's type name is an alias to an enum type")

	// Deprecated flags.
	Analyzer.Flags.StringVar(&unused, IgnorePatternFlag, "", "no effect (deprecated); see -"+IgnoreEnumMembersFlag+" instead")
	Analyzer.Flags.StringVar(&unused, CheckingStrategyFlag, "", "no effect (deprecated)")
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

	for typ, members := range findEnums(fPackageScopeOnly, fExcludeTypeAlias, pass.Pkg, inspect, pass.TypesInfo) {
		// TODO: need to also include alias RHS
		exportFact(pass, typ, members)
	}

	cfg := config{
		defaultSignifiesExhaustive: fDefaultSignifiesExhaustive,
		checkGeneratedFiles:        fCheckGeneratedFiles,
		ignoreEnumMembers:          fIgnoreEnumMembers.Get().(*regexp.Regexp),
	}
	checkSwitchStatements(pass, inspect, cfg)
	return nil, nil
}

// TODO(testing): add unit test
func assert(v bool, format string, args ...interface{}) {
	if !v {
		panic(fmt.Sprintf(format, args...))
	}
}

var (
	debug = log.New(os.Stderr, "", log.Lshortfile)
)
