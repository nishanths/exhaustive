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
	run := func(t *testing.T, setup func(), patterns ...string) []*analysistest.Result {
		t.Helper()
		resetFlags()
		setup()
		return analysistest.Run(t, filepath.Join(testdata, "test"), Analyzer, patterns...)
	}

	t.Run("switch", func(t *testing.T) {
		// Note: Many testdata files under the 'swtch' directory
		// include test cases for both expression switch
		// statements and map literals.
		//
		run(t, func() { fCheck = "switch,mapliteral" }, "test/check/swtch/general", "test/check/swtch/directive")
		run(t, func() { fCheck = "switch,mapliteral"; fCheckEnforceOnly = true }, "test/check/swtch/directive/enforce")
		run(t, func() { fCheck = "switch"; fDefaultEx = true }, "test/check/swtch/def")
		run(t, func() { fCheck = "switch"; fRequireDefaultCase = true }, "test/check/swtch/defrequire")
		run(t, func() { fCheck = "switch" }, "test/check/swtch/defrequire/directive1")
		run(t, func() { fCheck = "switch,mapliteral"; fCheckGenerated = true }, "test/check/swtch/generated")
		run(t, func() {
			fCheck = "switch,mapliteral"
			fExcludeType = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
				regexp.MustCompile("t1"),
				regexp.MustCompile("^test/check/swtch/pattern\\.P3$"), // no effect: not package-level declaration
				regexp.MustCompile("^test/check/swtch/typ\\.M1$"),
				regexp.MustCompile("^test/check/swtch/pattern\\.a2$"),
			}}
			fExcludeConst = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
				regexp.MustCompile("^test/check/swtch/pattern\\.x2+$"),
				regexp.MustCompile("^test/check/swtch/pattern\\.Y0$"), // no effect: not package-level declaration
				regexp.MustCompile("^test/check/swtch/typnew\\.Z3$"),
			}}
		}, "test/check/swtch/pattern")
		run(t, func() {
			fCheck = "switch,mapliteral"
			fIncludeType = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
				regexp.MustCompile("t1"),
				regexp.MustCompile("^test/check/swtch/pattern/includetype\\.P3$"), // no effect: not package-level declaration
				regexp.MustCompile("^test/check/swtch/typ\\.M1$"),
				regexp.MustCompile("^test/check/swtch/pattern/includetype\\.a2$"),
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

		}, "test/check/swtch/pattern/includetype")
		run(t, func() { fCheck = "switch" }, "test/check/swtch/unsupported")
	})

	t.Run("mapliteral", func(t *testing.T) {
		run(t, func() { fCheck = "mapliteral"; fCheckEnforceOnly = true }, "test/check/mapliteral/commentassoc")
	})

	t.Run("doc", func(t *testing.T) {
		if !testing.Short() {
			run(t, func() { fCheck = "switch,typeswitch" }, "test/packagedoc/...")
			run(t, func() { fCheck = "switch,typeswitch" }, "test/readmedoc/...")
		}
	})

	t.Run("typeswitch", func(t *testing.T) {
		run(t, func() { fCheck = "typeswitch" }, "test/check/typeswitch/general")
		run(t, func() { fCheck = "typeswitch" }, "test/check/typeswitch/alias")
		run(t, func() { fCheck = "typeswitch"; fRequireCaseNil = true }, "test/check/typeswitch/casenil")
		run(t, func() { fCheck = "typeswitch"; fDefaultEx = true }, "test/check/typeswitch/def")
		run(t, func() {
			fCheck = "typeswitch"
			fExcludeType = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
				regexp.MustCompile("^test/check/typeswitch/pattern\\.i$"),
				regexp.MustCompile("^test/check/typeswitch/pattern\\.j$"), // no effect: not package-level declaration
				regexp.MustCompile("^test/check/typeswitch/pattern\\.a2$"),
			}}
		}, "test/check/typeswitch/pattern")

		if !testing.Short() {
			run(t, func() { fCheck = "typeswitch" }, "test/check/typeswitch/gostd")
		}
	})
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

	run(t, func() { fCheck = "switch,mapliteral,typeswitch" }, "net")

	if !testing.Short() {
		// Note: This could specify the name "std" instead of
		// specifying individual package names, but in low memory
		// environments where tests might run the analysis will
		// commonly run out of memory for the former.
		run(t, func() { fCheck = "switch,mapliteral,typeswitch" },
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
}
