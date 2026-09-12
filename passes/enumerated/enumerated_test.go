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
	run := func(t *testing.T, setup func(), pattern ...string) []*analysistest.Result {
		t.Helper()
		resetFlags()
		setup()
		return analysistest.Run(t, analysistest.TestData(), Analyzer, pattern...)
	}

	run(t, func() {}, "enum/typ")
	run(t, func() {}, "enum/scope", "enum/scope/sub")
	run(t, func() { fIota = true }, "enum/iotavalue")
	run(t, func() { fNoBitwise = true }, "enum/nobitwise")
	run(t, func() { fIota = true; fNoBitwise = true }, "enum/iotavalue_nobitwise")
	run(t, func() { fPkgLevel = true }, "enum/pkglevel")
	run(t, func() {}, "enum/depchain/...")
	run(t, func() {}, "packagedoc/...")
}

func TestResult(t *testing.T) {
	// Note: c.C1 absent: not an enumerated type
	//       e.E1 absent: not a dependency of package a
	wantresult := map[string][]string{
		"enum/depchain/d.D1": {"const enum/depchain/d.X0 enum/depchain/d.D1", "const enum/depchain/d.X1 enum/depchain/d.D1"},
		"enum/depchain/d.D2": {"const enum/depchain/d.Z0 enum/depchain/d.D2", "const enum/depchain/d.Z1 enum/depchain/d.D2", "const enum/depchain/d.Z2 enum/depchain/d.D2"},
		"enum/depchain/c.C2": {"const enum/depchain/c.Y0 enum/depchain/c.C2", "const enum/depchain/c.Y1 enum/depchain/c.C2"},
		"enum/depchain/a.A1": {"const enum/depchain/a.Y0 enum/depchain/a.A1", "const enum/depchain/a.Y1 enum/depchain/a.A1"},
	}
	r := analysistest.Run(t, analysistest.TestData(), Analyzer, "enum/depchain/a")
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
