package exhaustive

import (
	"regexp"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestExhaustive(t *testing.T) {
	run := func(t *testing.T, pattern string, setup ...func()) {
		t.Helper()
		t.Run(pattern, func(t *testing.T) {
			resetFlags()
			// default to checking switch and map for test.
			fCheck = stringsFlag{
				[]string{
					string(elementSwitch),
					string(elementMap),
				},
				nil,
			}
			for _, f := range setup {
				f()
			}
			analysistest.Run(t, analysistest.TestData(), Analyzer, pattern)
		})
	}

	// Enum discovery, enum types.
	run(t, "enum/...")

	// Tests for the -check-generated flag.
	run(t, "generated-file/check-generated-off/...")
	run(t, "generated-file/check-generated-on/...", func() { fCheckGenerated = true })

	// Tests for the -default-signifies-exhaustive flag.
	// (For tests with this flag off, see other testdata packages
	// such as "general/...".)
	run(t, "default-signifies-exhaustive/default-absent/...", func() { fDefaultSignifiesExhaustive = true })
	run(t, "default-signifies-exhaustive/default-present/...", func() { fDefaultSignifiesExhaustive = true })

	// Tests for the -ignore-enum-member flag.
	run(t, "ignore-enum-member/...", func() {
		re := regexp.MustCompile(`_UNSPECIFIED$|^general/y\.Echinodermata$|^ignore-enum-member.User$`)
		fIgnoreEnumMembers = regexpFlag{re}
	})

	// Tests for -package-scope-only flag.
	run(t, "scope/allscope/...")
	run(t, "scope/pkgscope/...", func() { fPackageScopeOnly = true })

	// Program elements with ignore comment should not be
	// checked during implicitly exhaustive mode.
	run(t, "ignore-comment/...")

	// Program elements without enforce comment should not be
	// checked in explicitly exhaustive mode.
	run(t, "enforce-comment/...", func() {
		fExplicitExhaustiveSwitch = true
		fExplicitExhaustiveMap = true
	})

	// To satisfy exhaustiveness, it is sufficient for each unique constant
	// value of the members to be listed, not each member by name.
	run(t, "duplicate-enum-value/...")

	run(t, "typealias/...")
	run(t, "typeparam/...")

	// mixture of general tests.
	run(t, "general/...")
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("want nil error, got %s", err)
	}
}
