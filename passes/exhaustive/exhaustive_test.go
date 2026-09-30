package exhaustive

import (
	"errors"
	"regexp"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/packages"
)

func TestAnalyzer(t *testing.T) {
	run := func(t *testing.T, setup func(), patterns ...string) {
		t.Helper()
		resetFlags()
		setup()
		analysistest.Run(t, analysistest.TestData(), Analyzer, patterns...)
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

// This test does not assert that the analyses produce expected
// diagnostics, facts, results, etc. It only checks that the
// analyses run without errors and without panics on real-world
// packages.
func TestRealPackages(t *testing.T) {
	analyze := func(args []string, as ...*analysis.Analyzer) error {
		cfg := packages.Config{
			Mode:  packages.LoadAllSyntax | packages.NeedModule,
			Tests: true,
		}
		pkgs, err := packages.Load(&cfg, args...)
		if err != nil {
			return err
		}
		if len(pkgs) == 0 {
			return errors.New("matched no packages")
		}
		_, err = checker.Analyze(as, pkgs, nil)
		if err != nil {
			return err
		}
		return nil
	}

	run := func(t *testing.T, setup func(), patterns ...string) {
		t.Helper()
		resetFlags()
		setup()
		if err := analyze(patterns, Analyzer); err != nil {
			t.Errorf("unexpected error analyzing %v: %s", patterns, err)
		}
	}

	run(t, func() { fCheck = "switch" }, "net")

	if testing.Short() {
		t.Skip("skipping extra tests in short mode")
		return
	}

	run(t, func() { fCheck = "switch,mapliteral" },
		"fmt/...",
		"go/...",
		"html/template/...",
		"io/...",
		"net/...",
		"os/...",
		"reflect/...",
		"regexp/...",
		"runtime/...",
		"sync/...",
		"unsafe/...")
}
