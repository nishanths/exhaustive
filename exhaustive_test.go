// Integration-style tests using analysistest.

package exhaustive

import (
	"regexp"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// Tests for enum discovery.
func TestEnum(t *testing.T) {
	resetFlags()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "enum")
}

// Switch statements associated with the ignore directive comment should not
// have diagnostics.
func TestIgnoreComment(t *testing.T) {
	resetFlags()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "ignorecomment")
}

// For an enum switch to be exhaustive, it is sufficient for each unique enum
// value to be listed, not each unique member by name.
func TestDuplicateEnumValue(t *testing.T) {
	resetFlags()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "duplicateenumvalue")
}

// No diagnostics for missing enum members that match the supplied regular expression.
func TestIgnoreEnumMember(t *testing.T) {
	resetFlags()
	fIgnorePattern = regexpFlag{regexp.MustCompile("_UNSPECIFIED$|^general/y.Echinodermata$")}
	analysistest.Run(t, analysistest.TestData(), Analyzer, "ignoreenummember")
}

// Generated files should not have diagnostics.
func TestGenerated(t *testing.T) {
	resetFlags()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "generated")
}

// General exhaustiveness tests.
func TestGeneral(t *testing.T) {
	resetFlags()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "general/x", "general/y")
}

// Tests for '-fix' option.
func TestFix(t *testing.T) {
	resetFlags()
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), Analyzer, "fix")
}
