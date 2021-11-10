package exhaustive

import (
	"flag"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var _ flag.Value = (*regexpFlag)(nil)

// regexpFlag implements the flag.Value interface for parsing
// regular expression flag values.
type regexpFlag struct{ r *regexp.Regexp }

func (v *regexpFlag) String() string {
	if v == nil || v.r == nil {
		return ""
	}
	return v.r.String()
}

func (v *regexpFlag) Set(expr string) error {
	if expr == "" {
		v.r = nil
		return nil
	}

	r, err := regexp.Compile(expr)
	if err != nil {
		return err
	}

	v.r = r
	return nil
}

func (v *regexpFlag) value() *regexp.Regexp { return v.r }

func init() {
	var unused string

	// Public flags.
	Analyzer.Flags.BoolVar(&fCheckGeneratedFiles, CheckGeneratedFlag, false, "check switch statements in generated files")
	Analyzer.Flags.BoolVar(&fDefaultSignifiesExhaustive, DefaultSignifiesExhaustiveFlag, false, "presence of \"default\" case in switch statements satisfies exhaustiveness, even if all enum members are not listed")
	Analyzer.Flags.Var(&fIgnoreEnumMembers, IgnoreEnumMembersFlag, "enum members matching `regex` do not have to be listed in switch statements to satisfy exhaustiveness")
	Analyzer.Flags.BoolVar(&fPackageScopeOnly, PackageScopeOnlyFlag, false, "consider enums only in package scopes, not in inner scopes")

	// Deprecated flags.
	Analyzer.Flags.StringVar(&unused, IgnorePatternFlag, "", "no effect (deprecated); see -"+IgnoreEnumMembersFlag+" instead")
	Analyzer.Flags.StringVar(&unused, CheckingStrategyFlag, "", "no effect (deprecated)")
}

// Flag names used by the analyzer. They are exported for use by analyzer
// driver programs.
const (
	CheckGeneratedFlag             = "check-generated"
	DefaultSignifiesExhaustiveFlag = "default-signifies-exhaustive"
	IgnoreEnumMembersFlag          = "ignore-enum-members"
	PackageScopeOnlyFlag           = "package-scope-only"

	IgnorePatternFlag    = "ignore-pattern"    // Deprecated: see IgnoreEnumMembersFlag instead.
	CheckingStrategyFlag = "checking-strategy" // Deprecated.
)

var (
	// Public flags.
	fCheckGeneratedFiles        bool
	fDefaultSignifiesExhaustive bool
	fIgnoreEnumMembers          regexpFlag
	fPackageScopeOnly           bool
)

// resetFlags resets the flag variables to their default values.
// Useful in tests.
func resetFlags() {
	// Public flags.
	fCheckGeneratedFiles = false
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

	for typ, members := range findEnums(fPackageScopeOnly, pass.Pkg, inspect, pass.TypesInfo) {
		exportFact(pass, typ, members)
	}

	cfg := config{
		defaultSignifiesExhaustive: fDefaultSignifiesExhaustive,
		checkGeneratedFiles:        fCheckGeneratedFiles,
		ignoreEnumMembers:          fIgnoreEnumMembers.value(),
	}
	checkSwitchStatements(pass, inspect, cfg)
	return nil, nil
}

// NOTE: Fact types must remain gob-coding compatible.
// See TestFactsGob.
var _ analysis.Fact = (*enumMembersFact)(nil)

type enumMembersFact struct{ Members enumMembers }

func (f *enumMembersFact) AFact()         {}
func (f *enumMembersFact) String() string { return f.Members.factString() }

// exportFact exports the enum members for the given enum type.
func exportFact(pass *analysis.Pass, enumTyp enumType, members enumMembers) {
	pass.ExportObjectFact(enumTyp.factObject(), &enumMembersFact{members})
}

// importFact imports the enum members for the given possible enum type.
// An (_, false) return indicates that the enum type is not a known one.
func importFact(pass *analysis.Pass, possibleEnumType enumType) (enumMembers, bool) {
	var f enumMembersFact
	ok := pass.ImportObjectFact(possibleEnumType.factObject(), &f)
	if !ok {
		return enumMembers{}, false
	}
	return f.Members, true
}
