package exhaustive

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestEnum(t *testing.T) {
	resetFlags()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "enumvariant")
}

func TestSwitch(t *testing.T) {
	resetFlags()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "switch/x", "switch/y")
}

func TestSwitch_ignorePattern(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		resetFlags()
		if err := fIgnorePattern.Set("_UNSPECIFIED$|^switch/y.Echinodermata$"); !checkNoError(t, err) {
			return
		}
		analysistest.Run(t, analysistest.TestData(), Analyzer, "switch/ignorepattern")
	})
}

func TestSwitchFix(t *testing.T) {
	resetFlags()
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), Analyzer, "switchfix")
}
