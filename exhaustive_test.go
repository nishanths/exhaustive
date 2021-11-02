package exhaustive

import (
	"regexp"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// Integration-style tests using the analysistest package.
func TestAnalyzer(t *testing.T) {
	// Enum discovery.
	t.Run("enum", func(t *testing.T) {
		resetFlags()
		analysistest.Run(t, analysistest.TestData(), Analyzer, "enum")
	})

	// Switch statements associated with the ignore directive comment should not
	// have diagnostics.
	t.Run("ignore directive comment", func(t *testing.T) {
		resetFlags()
		analysistest.Run(t, analysistest.TestData(), Analyzer, "ignorecomment")
	})

	// For an enum switch to be exhaustive, it is sufficient for each unique enum
	// value to be listed, not each unique member by name.
	t.Run("duplicate enum value", func(t *testing.T) {
		resetFlags()
		analysistest.Run(t, analysistest.TestData(), Analyzer, "duplicateenumvalue")
	})

	// No diagnostics for missing enum members that match the supplied regular expression.
	t.Run("ignore enum member", func(t *testing.T) {
		resetFlags()
		fIgnoreEnumMembers = regexpFlag{regexp.MustCompile("_UNSPECIFIED$|^general/y.Echinodermata$")}
		analysistest.Run(t, analysistest.TestData(), Analyzer, "ignoreenummember")
	})

	// Generated files should not have diagnostics.
	t.Run("generated file", func(t *testing.T) {
		resetFlags()
		analysistest.Run(t, analysistest.TestData(), Analyzer, "generated")
	})

	// General tests (a mixture).
	t.Run("general", func(t *testing.T) {
		resetFlags()
		analysistest.Run(t, analysistest.TestData(), Analyzer, "general/...")
	})
}

// func TestCheckAndAdjustFlags(t *testing.T) {
// 	t.Run("setting both "+IgnorePatternFlag+" and "+IgnoreEnumMembersFlag+"is an error", func(t *testing.T) {
// 		resetFlags()
// 		fIgnorePattern = regexpFlag{regexp.MustCompile("a")}
// 		fIgnoreEnumMembers = regexpFlag{regexp.MustCompile("b")}

// 		err := checkAndAdjustFlags()

// 		assertError(t, err)
// 		got := err.Error()
// 		want := "cannot specify both -ignore-pattern and -ignore-enum-members"
// 		if got != want {
// 			t.Errorf("want %q, got %q", want, got)
// 		}
// 	})

// 	t.Run(IgnorePatternFlag+" value is copied to "+IgnoreEnumMembersFlag, func(t *testing.T) {
// 		resetFlags()
// 		fIgnorePattern = regexpFlag{regexp.MustCompile("a")}

// 		err := checkAndAdjustFlags()

// 		assertNoError(t, err)
// 		if fIgnorePattern.r.String() != fIgnoreEnumMembers.r.String() {
// 			t.Errorf("wrong flag value")
// 		}
// 	})
// }
