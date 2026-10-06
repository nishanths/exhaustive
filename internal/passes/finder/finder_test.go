package finder

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	run := func(t *testing.T, patterns ...string) []*analysistest.Result {
		t.Helper()
		return analysistest.Run(t, filepath.Join(testdata, "test"), Analyzer, patterns...)
	}
	run(t, "test/general")
	run(t, "test/general/sub")
}
