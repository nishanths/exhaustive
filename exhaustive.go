package exhaustive

import (
	"fmt"
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
	IgnorePatternFlag              = "ignore-pattern" // Deprecated. See IgnoreEnumMembersFlag instead.
	IgnoreEnumMembersFlag          = "ignore-enum-members"
	CheckingStrategyFlag           = "checking-strategy"
)

var (
	fDefaultSignifiesExhaustive bool
	fCheckGeneratedFiles        bool
	fDeprecatedIgnorePattern    string // Deprecated.
	fIgnoreEnumMembers          regexpFlag
	fCheckingStrategy           string
)

func init() {
	Analyzer.Flags.BoolVar(&fDefaultSignifiesExhaustive, DefaultSignifiesExhaustiveFlag, false, "presence of \"default\" case in switch statements satisfies exhaustiveness, even if all enum members are not listed")
	Analyzer.Flags.BoolVar(&fCheckGeneratedFiles, CheckGeneratedFlag, false, "check switch statements in generated files")
	Analyzer.Flags.StringVar(&fDeprecatedIgnorePattern, IgnorePatternFlag, "", "no effect (deprecated); see -"+IgnoreEnumMembersFlag+" instead")
	Analyzer.Flags.Var(&fIgnoreEnumMembers, IgnoreEnumMembersFlag, "enum members matching `regex` do not have to be listed in switch statements to satisfy exhaustiveness")
	Analyzer.Flags.StringVar(&fCheckingStrategy, CheckingStrategyFlag, "value", "`strategy` to use when checking exhaustiveness of switch statements; one of: \"value\", \"name\"")
}

// resetFlags resets the flag variables to their default values.
// Useful in tests.
func resetFlags() {
	fDefaultSignifiesExhaustive = false
	fCheckGeneratedFiles = false
	fDeprecatedIgnorePattern = ""
	fIgnoreEnumMembers = regexpFlag{}
	fCheckingStrategy = "value"
}

var Analyzer = &analysis.Analyzer{
	Name:      "exhaustive",
	Doc:       "check exhaustiveness of enum switch statements",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&enumsFact{}},
}

func run(pass *analysis.Pass) (interface{}, error) {
	e := findEnums(pass.Files, pass.TypesInfo)
	if len(e) != 0 {
		pass.ExportPackageFact(&enumsFact{Enums: e})
	}

	strategy, err := determineCheckingStrategy()
	if err != nil {
		return nil, err
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	cfg := config{
		defaultSignifiesExhaustive: fDefaultSignifiesExhaustive,
		checkGeneratedFiles:        fCheckGeneratedFiles,
		ignoreEnumMembers:          fIgnoreEnumMembers.Get().(*regexp.Regexp),
		checkingStrategy:           strategy,
	}

	err = checkSwitchStatements(pass, inspect, cfg)
	return nil, err
}

func determineCheckingStrategy() (checkingStrategy, error) {
	switch fCheckingStrategy {
	case "value":
		return byValue, nil
	case "name":
		return byName, nil
	default:
		return 0, fmt.Errorf("bad value %q for flag -%s", fCheckingStrategy, CheckingStrategyFlag)
	}
}
