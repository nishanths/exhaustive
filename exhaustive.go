package exhaustive

import (
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Flag names used by the analyzer. They are exported for use by analyzer
// driver programs.
const (
	DefaultSignifiesExhaustiveFlag = "default-signifies-exhaustive"
	CheckGeneratedFlag             = "check-generated"
	IgnoreEnumMembersFlag          = "ignore-enum-members"

	IgnorePatternFlag    = "ignore-pattern"    // Deprecated: see IgnoreEnumMembersFlag instead.
	CheckingStrategyFlag = "checking-strategy" // Deprecated.
)

var (
	fDefaultSignifiesExhaustive bool
	fCheckGeneratedFiles        bool
	fIgnoreEnumMembers          regexpFlag

	fDeprecatedIgnorePattern    string // Deprecated: see fIgnoreEnumMembers instead.
	fDeprecatedCheckingStrategy string // Deprecated.
)

func init() {
	Analyzer.Flags.BoolVar(&fDefaultSignifiesExhaustive, DefaultSignifiesExhaustiveFlag, false, "presence of \"default\" case in switch statements satisfies exhaustiveness, even if all enum members are not listed")
	Analyzer.Flags.BoolVar(&fCheckGeneratedFiles, CheckGeneratedFlag, false, "check switch statements in generated files")
	Analyzer.Flags.Var(&fIgnoreEnumMembers, IgnoreEnumMembersFlag, "enum members matching `regex` do not have to be listed in switch statements to satisfy exhaustiveness")

	Analyzer.Flags.StringVar(&fDeprecatedIgnorePattern, IgnorePatternFlag, "", "no effect (deprecated); see -"+IgnoreEnumMembersFlag+" instead")
	Analyzer.Flags.StringVar(&fDeprecatedCheckingStrategy, CheckingStrategyFlag, "", "no effect (deprecated)")
}

// resetFlags resets the flag variables to their default values.
// Useful in tests.
func resetFlags() {
	fDefaultSignifiesExhaustive = false
	fCheckGeneratedFiles = false
	fIgnoreEnumMembers = regexpFlag{}

	fDeprecatedIgnorePattern = ""
	fDeprecatedCheckingStrategy = ""
}

var Analyzer = &analysis.Analyzer{
	Name:      "exhaustive",
	Doc:       "check exhaustiveness of enum switch statements",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&enumMembersFact{}},
}

func run(pass *analysis.Pass) (interface{}, error) {
	pkgScopeOnly := false // TODO: make this configurable by flag
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	for typ, members := range findEnums(pkgScopeOnly, pass.Pkg, inspect, pass.TypesInfo) {
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
