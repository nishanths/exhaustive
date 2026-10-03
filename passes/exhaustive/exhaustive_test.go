package exhaustive

import (
	"errors"
	"path/filepath"
	"regexp"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/packages"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	run := func(t *testing.T, dir string, setup func(), patterns ...string) []*analysistest.Result {
		t.Helper()
		resetFlags()
		setup()
		return analysistest.Run(t, filepath.Join(testdata, "src", dir), Analyzer, patterns...)
	}

	// Note: Many testdata files under the 'swtch' directory
	// include test cases for both expression switch
	// statements and map literals.
	//
	run(t, "check", func() { fCheck = "switch,mapliteral" }, "check/swtch/general", "check/swtch/directive")
	run(t, "check", func() { fCheck = "switch,mapliteral"; fNeedEnforceDirective = true }, "check/swtch/directive/enforce")
	run(t, "check", func() { fCheck = "switch"; fDefaultEx = true }, "check/swtch/def")
	run(t, "check", func() { fCheck = "switch"; fDefaultRequired = true }, "check/swtch/defrequire")
	run(t, "check", func() { fCheck = "switch" }, "check/swtch/defrequire/directive1")
	run(t, "check", func() { fCheck = "switch,mapliteral"; fCheckGenerated = true }, "check/swtch/generated")
	run(t, "check", func() {
		fCheck = "switch,mapliteral"
		fExcludeType = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
			regexp.MustCompile("t1"),
			regexp.MustCompile("^check/swtch/pattern\\.P3$"), // no effect: not package-level declaration
			regexp.MustCompile("^check/swtch/typ\\.M1$"),
		}}
		fExcludeConst = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
			regexp.MustCompile("^check/swtch/pattern\\.x2+$"),
			regexp.MustCompile("^check/swtch/pattern\\.Y0$"), // no effect: not package-level declaration
			regexp.MustCompile("^check/swtch/typnew\\.Z3$"),
		}}
	}, "check/swtch/pattern")
	run(t, "check", func() {
		fCheck = "switch,mapliteral"
		fIncludeType = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
			regexp.MustCompile("t1"),
			regexp.MustCompile("^check/swtch/pattern\\.P3$"), // no effect: not package-level declaration
			regexp.MustCompile("^check/swtch/typ\\.M1$"),
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

	}, "check/swtch/pattern/includetype")
	run(t, "check", func() { fCheck = "switch" }, "check/swtch/unsupported")
	run(t, "check", func() { fCheck = "mapliteral"; fNeedEnforceDirective = true }, "check/mapliteral/commentassoc")
	run(t, "packagedoc", func() { fCheck = "switch" }, "packagedoc/...")
	run(t, "readmeexample", func() { fCheck = "switch" }, "readmeexample/...")
}

// This test does not check that the analysis produces expected
// diagnostics, facts, results, etc. It only checks that the
// analysis runs without errors and without panics on real
// packages, such as those in the standard library.
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

	run(t, func() { fCheck = "switch,mapliteral" }, "net")

	if testing.Short() {
		t.Skip("skipping extra tests in short mode")
		return
	}

	// Note: This could specify the name "std" instead of
	// specifying individual package names, but in low memory
	// environments where tests might run the analysis will
	// commonly run out of memory for the former.
	run(t, func() { fCheck = "switch,mapliteral" },
		"fmt",
		"go/...",
		"html/template",
		"io/...",
		"net/...",
		"os/...",
		"reflect",
		"regexp/...",
		"runtime/...",
		"sync/...",
		"unsafe")
}
