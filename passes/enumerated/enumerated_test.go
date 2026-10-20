package enumerated

import (
	"cmp"
	"go/types"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMain(m *testing.M) {
	testingExtraFacts = true
	Analyzer.FactTypes = append(Analyzer.FactTypes, new(enumeratedTypeFact), new(permittedTypeFact), new(enumeratedConstantFact))
	m.Run()
}

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	run := func(t *testing.T, setup func(), patterns ...string) []*analysistest.Result {
		t.Helper()
		resetFlags()
		setup()
		return analysistest.Run(t, filepath.Join(testdata, "analysis"), Analyzer, patterns...)
	}

	run(t, func() {}, "test/analysis/enum/typ")
	run(t, func() {}, "test/analysis/enum/scope", "test/analysis/enum/scope/sub")
	run(t, func() { fIota = true }, "test/analysis/enum/iotavalue")
	run(t, func() { fNoBitwise = true }, "test/analysis/enum/nobitwise")
	run(t, func() { fIota = true; fNoBitwise = true }, "test/analysis/enum/iotavalue_nobitwise")
	run(t, func() { fPkgLevel = true }, "test/analysis/enum/pkglevel")
	run(t, func() {}, "test/analysis/depchain/...")
	run(t, func() {}, "test/analysis/packagedoc/...")
}

func TestResult(t *testing.T) {
	wantresult := map[string][]string{
		"test/analysis/depchain/d.D1": {"const test/analysis/depchain/d.DX0 test/analysis/depchain/d.D1", "const test/analysis/depchain/d.DX1 test/analysis/depchain/d.D1"},
		"test/analysis/depchain/d.D2": {"const test/analysis/depchain/d.DZ0 test/analysis/depchain/d.D2", "const test/analysis/depchain/d.DZ1 test/analysis/depchain/d.D2", "const test/analysis/depchain/d.DZ2 test/analysis/depchain/d.D2"},
		"test/analysis/depchain/c.C2": {"const test/analysis/depchain/c.CY0 test/analysis/depchain/c.C2", "const test/analysis/depchain/c.CY1 test/analysis/depchain/c.C2"},
		"test/analysis/depchain/b.B1": {"const test/analysis/depchain/b.BY0 test/analysis/depchain/b.B1", "const test/analysis/depchain/b.BY1 test/analysis/depchain/b.B1"},
		"test/analysis/depchain/a.A1": {"const test/analysis/depchain/a.AX0 test/analysis/depchain/a.A1", "const test/analysis/depchain/a.AX1 test/analysis/depchain/a.A1", "const test/analysis/depchain/a.AX2 test/analysis/depchain/a.A1"},
	}
	r := analysistest.Run(t, filepath.Join(analysistest.TestData(), "analysis"), Analyzer, "test/analysis/depchain/a")
	if len(r) != 1 {
		t.Errorf("len(r): got: %d, want: 1", len(r))
	} else if got := transform(r[0].Action.Result.(Result)); !reflect.DeepEqual(got, wantresult) {
		t.Errorf("result not equal:\ngot:  %s\nwant: %s", got, wantresult)
	}
}

func transform(m map[*types.Named][]*types.Const) map[string][]string {
	cmpconst := func(a, b *types.Const) int { return cmp.Compare(a.Pos(), b.Pos()) }
	strnamed := (*types.Named).String
	strconst := (*types.Const).String

	out := make(map[string][]string)
	for t, cs := range m {
		cs := slices.Clone(cs)
		slices.SortFunc(cs, cmpconst)
		out[strnamed(t)] = slicemap(cs, strconst)
	}
	return out
}
