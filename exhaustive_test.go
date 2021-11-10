package exhaustive

import (
	"regexp"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestExhaustive(t *testing.T) {
	// Enum discovery.
	t.Run("enum", func(t *testing.T) {
		resetFlags()
		analysistest.Run(t, analysistest.TestData(), Analyzer, "enum/...")
	})

	// Switch statements with ignore directive comment should not
	// have diagnostics.
	t.Run("ignore directive comment", func(t *testing.T) {
		resetFlags()
		analysistest.Run(t, analysistest.TestData(), Analyzer, "ignorecomment/...")
	})

	// For an enum switch to be exhaustive, it is sufficient for each unique
	// constant value of the members to be listed, not each member by name.
	t.Run("duplicate enum value", func(t *testing.T) {
		resetFlags()
		analysistest.Run(t, analysistest.TestData(), Analyzer, "duplicateenumvalue/...")
	})

	// Tests for the -default-signifies-exhaustive flag.
	t.Run("default signifies exhaustive", func(t *testing.T) {
		resetFlags()
		fDefaultSignifiesExhaustive = true

		t.Run("default case absent", func(t *testing.T) {
			analysistest.Run(t, analysistest.TestData(), Analyzer, "defaultsignifiesexhaustive/defaultabsent/...")
		})

		t.Run("default case present", func(t *testing.T) {
			analysistest.Run(t, analysistest.TestData(), Analyzer, "defaultsignifiesexhaustive/defaultpresent/...")
		})
	})

	// There should be no diagnostics for missing enum members that match the
	// supplied regular expression.
	t.Run("ignore enum member", func(t *testing.T) {
		resetFlags()
		fIgnoreEnumMembers = regexpFlag{regexp.MustCompile(`_UNSPECIFIED$|^general/y\.Echinodermata$|^ignoreenummember.User$`)}
		analysistest.Run(t, analysistest.TestData(), Analyzer, "ignoreenummember/...")
	})

	// Generated files should not have diagnostics.
	t.Run("generated file", func(t *testing.T) {
		resetFlags()
		analysistest.Run(t, analysistest.TestData(), Analyzer, "generated/...")
	})

	// Switch statements using package-scoped and inner-scoped enums.
	t.Run("scope", func(t *testing.T) {
		t.Run("all scopes", func(t *testing.T) {
			resetFlags()
			fPackageScopeOnly = false
			analysistest.Run(t, analysistest.TestData(), Analyzer, "scope/allscope/...")
		})

		t.Run("package scope only", func(t *testing.T) {
			resetFlags()
			fPackageScopeOnly = true
			analysistest.Run(t, analysistest.TestData(), Analyzer, "scope/pkgscope/...")
		})
	})

	// Type alias switch statements.
	t.Run("type alias", func(t *testing.T) {
		resetFlags()
		analysistest.Run(t, analysistest.TestData(), Analyzer, "typealias/...")
	})

	// General tests (a mixture).
	t.Run("general", func(t *testing.T) {
		resetFlags()
		analysistest.Run(t, analysistest.TestData(), Analyzer, "general/...")
	})
}
