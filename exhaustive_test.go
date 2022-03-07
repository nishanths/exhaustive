package exhaustive

import (
	"regexp"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRegexpFlag(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		var v regexpFlag
		if got := v.value(); got != nil {
			t.Errorf("want nil, got %+v", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set(""); err != nil {
			t.Errorf("error unexpectedly non-nil: %v", err)
		}
		if got := v.value(); got != nil {
			t.Errorf("want nil, got %+v", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("bad input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set("("); err == nil {
			t.Errorf("error unexpectedly nil")
		}
		if got := v.value(); got != nil {
			t.Errorf("want nil, got %+v", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("good input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set("^foo$"); err != nil {
			t.Errorf("error unexpectedly non-nil: %v", err)
		}
		if v.value() == nil {
			t.Errorf("unexpectedly nil")
		}
		if !v.value().MatchString("foo") {
			t.Errorf("did not match")
		}
		if got, want := v.String(), regexp.MustCompile("^foo$").String(); got != want {
			t.Errorf("want %q, got %q", got, want)
		}
	})

	// The flag.Value interface doc says: "The flag package may call the
	// String method with a zero-valued receiver, such as a nil pointer."
	t.Run("String() nil receiver", func(t *testing.T) {
		var v *regexpFlag
		// expect no panic, and ...
		if got := v.String(); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})
}

func TestExhaustive(t *testing.T) {
	run := func(t *testing.T, pattern string, setup ...func()) {
		t.Helper()
		t.Run(pattern, func(t *testing.T) {
			resetFlags()
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

	// Switch statements with ignore directive comment should not be checked during implicitly exhaustive switch
	// mode
	run(t, "ignore-comment/...")

	// Switch statements without enforce directive comment should not be checked during explicitly exhaustive
	// switch mode
	run(t, "enforce-comment/...", func() { fExplicitExhaustiveSwitch = true })

	// To satisfy exhaustiveness, it is sufficient for each unique constant
	// value of the members to be listed, not each member by name.
	run(t, "duplicate-enum-value/...")

	// Type alias switch statements.
	run(t, "typealias/...")

	// General tests (a mixture).
	run(t, "general/...")
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("want nil error, got %s", err)
	}
}
