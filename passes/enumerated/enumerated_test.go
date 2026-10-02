package enumerated

import (
	"cmp"
	"go/types"
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
	run := func(t *testing.T, setup func(), patterns ...string) []*analysistest.Result {
		t.Helper()
		resetFlags()
		setup()
		return analysistest.Run(t, analysistest.TestData(), Analyzer, patterns...)
	}

	run(t, func() {}, "enum/typ")
	run(t, func() {}, "enum/scope", "enum/scope/sub")
	run(t, func() { fIota = true }, "enum/iotavalue")
	run(t, func() { fNoBitwise = true }, "enum/nobitwise")
	run(t, func() { fIota = true; fNoBitwise = true }, "enum/iotavalue_nobitwise")
	run(t, func() { fPkgLevel = true }, "enum/pkglevel")
	run(t, func() {}, "depchain/...")
	run(t, func() {}, "packagedoc/...")
}

func TestResult(t *testing.T) {
	wantresult := map[string][]string{
		"depchain/d.D1": {"const depchain/d.DX0 depchain/d.D1", "const depchain/d.DX1 depchain/d.D1"},
		"depchain/d.D2": {"const depchain/d.DZ0 depchain/d.D2", "const depchain/d.DZ1 depchain/d.D2", "const depchain/d.DZ2 depchain/d.D2"},
		"depchain/c.C2": {"const depchain/c.CY0 depchain/c.C2", "const depchain/c.CY1 depchain/c.C2"},
		"depchain/b.B1": {"const depchain/b.BY0 depchain/b.B1", "const depchain/b.BY1 depchain/b.B1"},
		"depchain/a.A1": {"const depchain/a.AY0 depchain/a.A1", "const depchain/a.AY1 depchain/a.A1"},
	}
	r := analysistest.Run(t, analysistest.TestData(), Analyzer, "depchain/a")
	if len(r) != 1 {
		t.Errorf("len(r): got: %d, want: 1", len(r))
	} else if got := transform(r[0].Action.Result.(Result)); !reflect.DeepEqual(got, wantresult) {
		t.Errorf("result not equal:\ngot:  %s\nwant: %s", got, wantresult)
	}
}

func transform(m map[*types.Named][]*types.Const) map[string][]string {
	sortconst := func(a, b *types.Const) int { return cmp.Compare(a.Pos(), b.Pos()) }
	strnamed := (*types.Named).String
	strconst := (*types.Const).String

	out := make(map[string][]string)
	for t, cs := range m {
		cs := slices.Clone(cs)
		slices.SortFunc(cs, sortconst)
		out[strnamed(t)] = slicemap(cs, strconst)
	}
	return out
}
