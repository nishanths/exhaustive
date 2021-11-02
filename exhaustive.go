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
)

var (
	fDefaultSignifiesExhaustive bool
	fCheckGeneratedFiles        bool
	fIgnorePattern              regexpFlag // Deprecated.
	fIgnoreEnumMembers          regexpFlag
)

func init() {
	Analyzer.Flags.BoolVar(&fDefaultSignifiesExhaustive, DefaultSignifiesExhaustiveFlag, false, "switch statements are to be considered exhaustive if a 'default' case is present, even if all enum members aren't listed in the switch")
	Analyzer.Flags.BoolVar(&fCheckGeneratedFiles, CheckGeneratedFlag, false, "check switch statements in generated files")
	Analyzer.Flags.Var(&fIgnorePattern, IgnorePatternFlag, "deprecated: use -"+IgnoreEnumMembersFlag+" instead")
	Analyzer.Flags.Var(&fIgnoreEnumMembers, IgnoreEnumMembersFlag, "ignore enum members matching `regex` when checking for exhaustiveness")
}

// resetFlags resets the flag variables to their default values.
// Useful in tests.
func resetFlags() {
	fDefaultSignifiesExhaustive = false
	fCheckGeneratedFiles = false
	fIgnorePattern = regexpFlag{}
	fIgnoreEnumMembers = regexpFlag{}
}

var Analyzer = &analysis.Analyzer{
	Name:      "exhaustive",
	Doc:       "check exhaustiveness of enum switch statements",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&enumsFact{}},
}

func run(pass *analysis.Pass) (interface{}, error) {
	if err := checkAndAdjustFlags(); err != nil {
		return nil, err
	}

	e := findEnums(pass.Files, pass.TypesInfo)
	if len(e) != 0 {
		pass.ExportPackageFact(&enumsFact{Enums: e})
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	cfg := config{
		defaultSignifiesExhaustive: fDefaultSignifiesExhaustive,
		checkGeneratedFiles:        fCheckGeneratedFiles,
		ignoreEnumMembers:          fIgnoreEnumMembers.Get().(*regexp.Regexp),
		hitlistStrategy:            byValue, // TODO: support other hitlist strategies via a user-specified flag
	}

	err := checkSwitchStatements(pass, inspect, cfg)
	return nil, err
}

func checkAndAdjustFlags() error {
	if fIgnorePattern.Get().(*regexp.Regexp) != nil && fIgnoreEnumMembers.Get().(*regexp.Regexp) != nil {
		return fmt.Errorf("cannot specify both -%s and -%s", IgnorePatternFlag, IgnoreEnumMembersFlag)
	}
	if fIgnorePattern.Get().(*regexp.Regexp) != nil {
		fIgnoreEnumMembers = fIgnorePattern
	}
	return nil
}
