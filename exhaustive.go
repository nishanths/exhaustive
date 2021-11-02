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
	IgnorePatternFlag              = "ignore-pattern" // Deprecated. See IgnoreEnumMembersFlag instead.
	IgnoreEnumMembersFlag          = "ignore-enum-members"
	ByNameFlag                     = "by-name"
)

var (
	fDefaultSignifiesExhaustive bool
	fCheckGeneratedFiles        bool
	fDeprecatedIgnorePattern    string // Deprecated.
	fIgnoreEnumMembers          regexpFlag
	fByName                     bool
)

func init() {
	Analyzer.Flags.BoolVar(&fDefaultSignifiesExhaustive, DefaultSignifiesExhaustiveFlag, false, "presence of 'default' case in a switch statement satisfies exhaustiveness, even if all enum members aren't listed")
	Analyzer.Flags.BoolVar(&fCheckGeneratedFiles, CheckGeneratedFlag, false, "check switch statements in generated files")
	Analyzer.Flags.StringVar(&fDeprecatedIgnorePattern, IgnorePatternFlag, "", "no effect (deprecated); see -"+IgnoreEnumMembersFlag+" instead")
	Analyzer.Flags.Var(&fIgnoreEnumMembers, IgnoreEnumMembersFlag, "enum members matching `regex` do not have to be listed in a switch statement to satisfy exhaustiveness")
	Analyzer.Flags.BoolVar(&fByName, ByNameFlag, false, "require every enum member name to be listed in a switch statement to satisfy exhaustiveness, as opposed to every enum member value (which is the default behavior)")
}

// resetFlags resets the flag variables to their default values.
// Useful in tests.
func resetFlags() {
	fDefaultSignifiesExhaustive = false
	fCheckGeneratedFiles = false
	fDeprecatedIgnorePattern = ""
	fIgnoreEnumMembers = regexpFlag{}
	fByName = false
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

	strategy := byValue
	if fByName {
		strategy = byName
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	cfg := config{
		defaultSignifiesExhaustive: fDefaultSignifiesExhaustive,
		checkGeneratedFiles:        fCheckGeneratedFiles,
		ignoreEnumMembers:          fIgnoreEnumMembers.Get().(*regexp.Regexp),
		checkingStrategy:           strategy,
	}

	err := checkSwitchStatements(pass, inspect, cfg)
	return nil, err
}
