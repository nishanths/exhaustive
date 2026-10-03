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
		return analysistest.Run(t, filepath.Join(testdata, "test"), Analyzer, patterns...)
	}

	run(t, func() {}, "test/enum/typ")
	run(t, func() {}, "test/enum/scope", "test/enum/scope/sub")
	run(t, func() { fIota = true }, "test/enum/iotavalue")
	run(t, func() { fNoBitwise = true }, "test/enum/nobitwise")
	run(t, func() { fIota = true; fNoBitwise = true }, "test/enum/iotavalue_nobitwise")
	run(t, func() { fPkgLevel = true }, "test/enum/pkglevel")
	run(t, func() {}, "test/depchain/...")
	run(t, func() {}, "test/packagedoc/...")
}

func TestResult(t *testing.T) {
	wantresult := map[string][]string{
		"test/depchain/d.D1": {"const test/depchain/d.DX0 test/depchain/d.D1", "const test/depchain/d.DX1 test/depchain/d.D1"},
		"test/depchain/d.D2": {"const test/depchain/d.DZ0 test/depchain/d.D2", "const test/depchain/d.DZ1 test/depchain/d.D2", "const test/depchain/d.DZ2 test/depchain/d.D2"},
		"test/depchain/c.C2": {"const test/depchain/c.CY0 test/depchain/c.C2", "const test/depchain/c.CY1 test/depchain/c.C2"},
		"test/depchain/b.B1": {"const test/depchain/b.BY0 test/depchain/b.B1", "const test/depchain/b.BY1 test/depchain/b.B1"},
		"test/depchain/a.A1": {"const test/depchain/a.AX0 test/depchain/a.A1", "const test/depchain/a.AX1 test/depchain/a.A1", "const test/depchain/a.AX2 test/depchain/a.A1"},
	}
	r := analysistest.Run(t, filepath.Join(analysistest.TestData(), "test"), Analyzer, "test/depchain/a")
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
