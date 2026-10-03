package finder

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	run := func(t *testing.T, dir string, setup func(), patterns ...string) []*analysistest.Result {
		t.Helper()
		setup()
		return analysistest.Run(t, filepath.Join(testdata, "test"), Analyzer, patterns...)
	}

	// run(t, "all", func() {}, "all/...")
	// XXX: add tests
	_ = run
}
