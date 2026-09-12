package exhaustive

import (
	"regexp"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	run := func(t *testing.T, setup func(), pattern ...string) {
		t.Helper()
		resetFlags()
		setup()
		analysistest.Run(t, analysistest.TestData(), Analyzer, pattern...)
	}

	// Note: Some 'swtch' testdata files include test cases for
	// both expression switch statements and map literals.
	//
	run(t, func() { fCheck = "switch,mapliteral" }, "swtch/general", "swtch/directive")
	run(t, func() { fCheck = "switch,mapliteral"; fNeedEnforceDirective = true }, "swtch/directive/enforce")
	run(t, func() { fCheck = "switch"; fDefaultEx = true }, "swtch/def")
	run(t, func() { fCheck = "switch"; fDefaultRequired = true }, "swtch/defrequire")
	run(t, func() { fCheck = "switch" }, "swtch/defrequire/directivetrue")
	run(t, func() { fCheck = "switch,mapliteral"; fCheckGenerated = true }, "swtch/generated")
	run(t, func() {
		fCheck = "switch"
		fExcludeType = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
			regexp.MustCompile("t1"),
			regexp.MustCompile("^swtch/pattern\\.P3$"), // no effect: not package-level declaration
			regexp.MustCompile("^swtch/typ\\.M1$"),
		}}
		fExcludeConst = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
			regexp.MustCompile("^swtch/pattern\\.x2+$"),
			regexp.MustCompile("^swtch/pattern\\.Y0$"), // no effect: not package-level declaration
			regexp.MustCompile("^swtch/typnew\\.Z3$"),
		}}
	}, "swtch/pattern")
	run(t, func() {
		fCheck = "switch"
		fIncludeType = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
			regexp.MustCompile("t1"),
			regexp.MustCompile("^swtch/pattern\\.P3$"), // no effect: not package-level declaration
			regexp.MustCompile("^swtch/typ\\.M1$"),
		}}
		fExcludeConst = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
			// Note: Exclusion of constants is by
			// name, not by value.
			// Therefore, though V0 == VV0 by value, the value
			// corresponding to VV0 is still required
			// in the switch statement cases.
			regexp.MustCompile("typ\\.V0$"),
			regexp.MustCompile("typ\\.V4$"),
		}}

	}, "swtch/pattern/includetype")
	run(t, func() { fCheck = "switch" }, "swtch/unsupported")
	run(t, func() { fCheck = "mapliteral"; fNeedEnforceDirective = true }, "mapliteral/commentassoc")
	run(t, func() { fCheck = "switch" }, "packagedoc/...")
	run(t, func() { fCheck = "switch" }, "readmeexample/...")
}
