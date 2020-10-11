package exhaustive

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestEnum(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "enumvariant/...")
}

func TestSwitch(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "switch/...")
}

func TestSwitchFix(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), Analyzer, "switchfix/...")
}
