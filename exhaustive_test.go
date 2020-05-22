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

func TestMap(t *testing.T) {
	fCheckMaps = true
	analysistest.Run(t, analysistest.TestData(), Analyzer, "map/...")
}

func TestSwitchFix(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "switchfix/...")
}
