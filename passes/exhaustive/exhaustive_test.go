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
		return analysistest.Run(t, filepath.Join(testdata, "analysis"), Analyzer, patterns...)
	}

	t.Run("exprswitch", func(t *testing.T) {
		// Note: Many testdata files under the 'exprswitch' directory
		// include test cases for both expression switch
		// statements and map literals.
		//
		run(t, func() { fCheck = "switch,mapliteral" }, "test/analysis/exprswitch/general")
		run(t, func() { fCheck = "switch,mapliteral" }, "test/analysis/exprswitch/directive")
		run(t, func() { fCheck = "switch,mapliteral"; fCheckEnforceOnly = true }, "test/analysis/exprswitch/directive/enforce")
		run(t, func() { fCheck = "switch"; fDefaultEx = true }, "test/analysis/exprswitch/def")
		run(t, func() { fCheck = "switch"; fRequireDefaultCase = true }, "test/analysis/exprswitch/defrequire")
		run(t, func() { fCheck = "switch" }, "test/analysis/exprswitch/defrequire/directive1")
		run(t, func() { fCheck = "switch,mapliteral"; fCheckGenerated = true }, "test/analysis/exprswitch/generated")
		run(t, func() {
			fCheck = "switch,mapliteral"
			fExcludeType = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
				regexp.MustCompile("t1"),
				regexp.MustCompile("^test/analysis/exprswitch/pattern\\.P3$"), // no effect: not package-level declaration
				regexp.MustCompile("^test/analysis/exprswitch/typ\\.M1$"),
				regexp.MustCompile("^test/analysis/exprswitch/pattern\\.a2$"),
			}}
			fExcludeConst = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
				regexp.MustCompile("^test/analysis/exprswitch/pattern\\.x2+$"),
				regexp.MustCompile("^test/analysis/exprswitch/pattern\\.Y0$"), // no effect: not package-level declaration
				regexp.MustCompile("^test/analysis/exprswitch/typnew\\.Z3$"),
			}}
		}, "test/analysis/exprswitch/pattern")
		run(t, func() {
			fCheck = "switch,mapliteral"
			fIncludeType = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
				regexp.MustCompile("t1"),
				regexp.MustCompile("^test/analysis/exprswitch/pattern/includetype\\.P3$"), // no effect: not package-level declaration
				regexp.MustCompile("^test/analysis/exprswitch/typ\\.M1$"),
				regexp.MustCompile("^test/analysis/exprswitch/pattern/includetype\\.a2$"),
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

		}, "test/analysis/exprswitch/pattern/includetype")
		run(t, func() { fCheck = "switch" }, "test/analysis/exprswitch/unsupported")

		if !testing.Short() {
			run(t, func() { fCheck = "switch" }, "test/analysis/exprswitch/gostd")
		}
	})

	t.Run("mapliteral", func(t *testing.T) {
		run(t, func() { fCheck = "mapliteral"; fCheckEnforceOnly = true }, "test/analysis/mapliteral/commentassoc")
	})

	t.Run("doc", func(t *testing.T) {
		if !testing.Short() {
			run(t, func() { fCheck = "switch,typeswitch"; fCheckImpossible = true }, "test/analysis/packagedoc/...")
			run(t, func() { fCheck = "switch,typeswitch" }, "test/analysis/readmedoc/...")
		}
	})

	t.Run("typeswitch", func(t *testing.T) {
		run(t, func() { fCheck = "typeswitch" },
			"test/analysis/typeswitch/general",
			"test/analysis/typeswitch/general/embedded",
			"test/analysis/typeswitch/general/alias",
			"test/analysis/typeswitch/general/generic")
		run(t, func() { fCheck = "typeswitch"; ; fCheckImpossible = true }, "test/analysis/typeswitch/general/caseinterface")
		run(t, func() { fCheck = "typeswitch"; fRequireCaseNil = true }, "test/analysis/typeswitch/casenil")
		run(t, func() { fCheck = "typeswitch"; fDefaultEx = true }, "test/analysis/typeswitch/def")
		run(t, func() {
			fCheck = "typeswitch"
			fExcludeType = repeatFlag[*regexp.Regexp]{vals: []*regexp.Regexp{
				regexp.MustCompile("^test/analysis/typeswitch/pattern\\.i$"),
				regexp.MustCompile("^test/analysis/typeswitch/pattern\\.j$"), // no effect: not package-level declaration
				regexp.MustCompile("^test/analysis/typeswitch/pattern\\.a2$"),
			}}
		}, "test/analysis/typeswitch/pattern")
		run(t, func() { fCheck = "typeswitch"; fCheckImpossible = true }, "test/analysis/typeswitch/impossible")

		if !testing.Short() {
			run(t, func() { fCheck = "typeswitch" }, "test/analysis/typeswitch/gostd")
		}
	})
}

// This test does not check that the analysis reports expected
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

	run(t, func() { fCheck = "switch,mapliteral,typeswitch"; fCheckImpossible = true }, "net")

	if !testing.Short() {
		// Note: This could specify the name "std" instead of
		// specifying individual package names, but in low memory
		// environments where tests might run the analysis will
		// commonly run out of memory for the former.
		run(t, func() { fCheck = "switch,mapliteral,typeswitch"; fCheckImpossible = true },
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
