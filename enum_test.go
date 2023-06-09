package exhaustive

import (
	"fmt"
	"go/token"
	"reflect"
	"sort"
	"testing"

	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

// checkEnumMembersLiteral checks that an enumMembers literal is correctly
// defined in tests.
func checkEnumMembersLiteral(id string, v enumMembers) {
	var count int
	for _, names := range v.ValueToNames {
		count += len(names)
	}
	if len(v.Names) != len(v.NameToPos) {
		panic(fmt.Sprintf("%s: wrong lengths: %d != %d (test definition bug)", id, len(v.Names), len(v.NameToPos)))
	}
	if len(v.Names) != len(v.NameToValue) {
		panic(fmt.Sprintf("%s: wrong lengths: %d != %d (test definition bug)", id, len(v.Names), len(v.NameToValue)))
	}
	if len(v.Names) != count {
		panic(fmt.Sprintf("%s: wrong lengths: %d != %d (test definition bug)", id, len(v.Names), count))
	}
}

func TestEnumMembers_add(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var v enumMembers
		v.add("foo", "\"A\"", 10)
		v.add("z", "X", 20)
		v.add("bar", "\"B\"", 30)
		v.add("y", "Y", 40)
		v.add("x", "X", 50)

		if got, want := v.Names, []string{"foo", "z", "bar", "y", "x"}; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}

		if got, want := v.NameToPos, map[string]token.Pos{
			"foo": 10,
			"z":   20,
			"bar": 30,
			"y":   40,
			"x":   50,
		}; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}

		if got, want := v.NameToValue, map[string]constantValue{
			"foo": "\"A\"",
			"z":   "X",
			"bar": "\"B\"",
			"y":   "Y",
			"x":   "X",
		}; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}

		if got, want := v.ValueToNames, map[constantValue][]string{
			"\"A\"": {"foo"},
			"\"B\"": {"bar"},
			"X":     {"z", "x"},
			"Y":     {"y"},
		}; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	// TODO: add tests for iota, repeated values, ...
}

func TestFindEnums(t *testing.T) {
	transform := func(in map[enumType]enumMembers) []checkEnum {
		var out []checkEnum
		for typ, members := range in {
			out = append(out, checkEnum{typ.TypeName.Name(), members})
		}
		return out
	}

	testdataEnumPkg := func() *packages.Package {
		cfg := &packages.Config{Mode: packages.NeedTypesInfo | packages.NeedTypes | packages.NeedSyntax}
		pkgs, err := packages.Load(cfg, "./testdata/src/enum")
		if err != nil {
			panic(err)
		}
		return pkgs[0]
	}()
	inspect := inspector.New(testdataEnumPkg.Syntax)

	for _, pkgOnly := range [...]bool{false, true} {
		t.Run(fmt.Sprint("pkgOnly", pkgOnly), func(t *testing.T) {
			result := findEnums(pkgOnly, testdataEnumPkg.Types, inspect, testdataEnumPkg.TypesInfo)
			checkEnums(t, transform(result), pkgOnly)
		})
	}
}

// see func checkEnums.
type checkEnum struct {
	typeName string
	members  enumMembers
}

func equalCheckEnum(t *testing.T, got, want checkEnum) {
	t.Helper()
	if got.typeName != want.typeName {
		t.Errorf("type name: got %s, want %s", got.typeName, want.typeName)
	}
	if !reflect.DeepEqual(got.members, want.members) {
		t.Errorf("type name %s: members: got %+v, want %+v", want.typeName, got.members, want.members)
	}
}

type byNameAndMembers []checkEnum

func (c byNameAndMembers) Len() int      { return len(c) }
func (c byNameAndMembers) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c byNameAndMembers) Less(i, j int) bool {
	if c[i].typeName != c[j].typeName {
		return c[i].typeName < c[j].typeName
	}
	return len(c[i].members.Names) < len(c[j].members.Names)
}

func checkEnums(t *testing.T, got []checkEnum, pkgOnly bool) {
	t.Helper()

	wantPkg := []checkEnum{
		{"VarConstMixed", enumMembers{
			[]string{"VCMixedB"},
			map[string]token.Pos{
				"VCMixedB": 0,
			},
			map[string]constantValue{
				"VCMixedB": `1`,
			},
			map[constantValue][]string{
				`1`: {"VCMixedB"},
			},
		}},
		{"IotaEnum", enumMembers{
			[]string{"IotaA", "IotaB"},
			map[string]token.Pos{
				"IotaA": 0,
				"IotaB": 0,
			},
			map[string]constantValue{
				"IotaA": `0`,
				"IotaB": `2`,
			},
			map[constantValue][]string{
				`0`: {"IotaA"},
				`2`: {"IotaB"},
			},
		}},
		{"RepeatedValue", enumMembers{
			[]string{"RepeatedValueA", "RepeatedValueB"},
			map[string]token.Pos{
				"RepeatedValueA": 0,
				"RepeatedValueB": 0,
			},
			map[string]constantValue{
				"RepeatedValueA": `1`,
				"RepeatedValueB": `1`,
			},
			map[constantValue][]string{
				`1`: {"RepeatedValueA", "RepeatedValueB"},
			},
		}},
		{"AcrossBlocksDeclsFiles", enumMembers{
			[]string{"Here", "Separate", "There"},
			map[string]token.Pos{
				"Here":     0,
				"Separate": 0,
				"There":    0,
			},
			map[string]constantValue{
				"Here":     `0`,
				"Separate": `1`,
				"There":    `2`,
			},
			map[constantValue][]string{
				`0`: {"Here"},
				`1`: {"Separate"},
				`2`: {"There"},
			},
		}},
		{"UnexportedMembers", enumMembers{
			[]string{"unexportedMembersA", "unexportedMembersB"},
			map[string]token.Pos{
				"unexportedMembersA": 0,
				"unexportedMembersB": 0,
			},
			map[string]constantValue{
				"unexportedMembersA": `1`,
				"unexportedMembersB": `2`,
			},
			map[constantValue][]string{
				`1`: {"unexportedMembersA"},
				`2`: {"unexportedMembersB"},
			},
		}},
		{"ParenVal", enumMembers{
			[]string{"ParenVal0", "ParenVal1"},
			map[string]token.Pos{
				"ParenVal0": 0,
				"ParenVal1": 0,
			},
			map[string]constantValue{
				"ParenVal0": `0`,
				"ParenVal1": `1`,
			},
			map[constantValue][]string{
				`0`: {"ParenVal0"},
				`1`: {"ParenVal1"},
			},
		}},
		{"EnumRHS", enumMembers{
			[]string{"EnumRHS_A", "EnumRHS_B"},
			map[string]token.Pos{
				"EnumRHS_A": 0,
				"EnumRHS_B": 0,
			},
			map[string]constantValue{
				"EnumRHS_A": `0`,
				"EnumRHS_B": `1`,
			},
			map[constantValue][]string{
				`0`: {"EnumRHS_A"},
				`1`: {"EnumRHS_B"},
			},
		}},
		{"WithMethod", enumMembers{
			[]string{"WithMethodA", "WithMethodB"},
			map[string]token.Pos{
				"WithMethodA": 0,
				"WithMethodB": 0,
			},
			map[string]constantValue{
				"WithMethodA": `1`,
				"WithMethodB": `2`,
			},
			map[constantValue][]string{
				`1`: {"WithMethodA"},
				`2`: {"WithMethodB"},
			},
		}},
		{"T", enumMembers{
			[]string{"A", "B"},
			map[string]token.Pos{
				"A": 0,
				"B": 0,
			},
			map[string]constantValue{
				"A": `0`,
				"B": `1`,
			},
			map[constantValue][]string{
				`0`: {"A"},
				`1`: {"B"},
			},
		}},
		{"PkgRequireSameLevel", enumMembers{
			[]string{"PA"},
			map[string]token.Pos{
				"PA": 0,
			},
			map[string]constantValue{
				"PA": `200`,
			},
			map[constantValue][]string{
				`200`: {"PA"},
			},
		}},
		{"UIntEnum", enumMembers{
			[]string{"UIntA", "UIntB"},
			map[string]token.Pos{
				"UIntA": 0,
				"UIntB": 0,
			},
			map[string]constantValue{
				"UIntA": "0",
				"UIntB": "1",
			},
			map[constantValue][]string{
				"0": {"UIntA"},
				"1": {"UIntB"},
			},
		}},
		{"StringEnum", enumMembers{
			[]string{"StringA", "StringB", "StringC"},
			map[string]token.Pos{
				"StringA": 0,
				"StringB": 0,
				"StringC": 0,
			},
			map[string]constantValue{
				"StringA": `"stringa"`,
				"StringB": `"stringb"`,
				"StringC": `"stringc"`,
			},
			map[constantValue][]string{
				`"stringa"`: {"StringA"},
				`"stringb"`: {"StringB"},
				`"stringc"`: {"StringC"},
			},
		}},
		{"RuneEnum", enumMembers{
			[]string{"RuneA"},
			map[string]token.Pos{
				"RuneA": 0,
			},
			map[string]constantValue{
				"RuneA": `97`,
			},
			map[constantValue][]string{
				`97`: {"RuneA"},
			},
		}},
		{"ByteEnum", enumMembers{
			[]string{"ByteA"},
			map[string]token.Pos{
				"ByteA": 0,
			},
			map[string]constantValue{
				"ByteA": `97`,
			},
			map[constantValue][]string{
				`97`: {"ByteA"},
			},
		}},
		{"Int32Enum", enumMembers{
			[]string{"Int32A", "Int32B"},
			map[string]token.Pos{
				"Int32A": 0,
				"Int32B": 0,
			},
			map[string]constantValue{
				"Int32A": "0",
				"Int32B": "1",
			},
			map[constantValue][]string{
				"0": {"Int32A"},
				"1": {"Int32B"},
			},
		}},
		{"Float64Enum", enumMembers{
			[]string{"Float64A", "Float64B"},
			map[string]token.Pos{
				"Float64A": 0,
				"Float64B": 0,
			},
			map[string]constantValue{
				"Float64A": `0`,
				"Float64B": `1`,
			},
			map[constantValue][]string{
				`0`: {"Float64A"},
				`1`: {"Float64B"},
			},
		}},
	}

	for _, c := range wantPkg {
		checkEnumMembersLiteral(c.typeName, c.members)
	}

	wantInner := []checkEnum{
		{"InnerRequireSameLevel", enumMembers{
			[]string{"IX", "IY"},
			map[string]token.Pos{
				"IX": 0,
				"IY": 0,
			},
			map[string]constantValue{
				"IX": `200`,
				"IY": `200`,
			},
			map[constantValue][]string{
				`200`: {"IX", "IY"},
			},
		}},
		{"T", enumMembers{
			[]string{"C", "D", "E", "F"},
			map[string]token.Pos{
				"C": 0,
				"D": 0,
				"E": 0,
				"F": 0,
			},
			map[string]constantValue{
				"C": `0`,
				"D": `1`,
				"E": `42`,
				"F": `43`,
			},
			map[constantValue][]string{
				`0`:  {"C"},
				`1`:  {"D"},
				`42`: {"E"},
				`43`: {"F"},
			},
		}},
		{"T", enumMembers{
			[]string{"A", "B"},
			map[string]token.Pos{
				"A": 0,
				"B": 0,
			},
			map[string]constantValue{
				"A": `0`,
				"B": `1`,
			},
			map[constantValue][]string{
				`0`: {"A"},
				`1`: {"B"},
			},
		}},
	}

	for _, c := range wantInner {
		checkEnumMembersLiteral(c.typeName, c.members)
	}

	want := append([]checkEnum{}, wantPkg...)
	if !pkgOnly {
		want = append(want, wantInner...)
	}

	sort.Sort(byNameAndMembers(want))
	sort.Sort(byNameAndMembers(got))

	if len(got) != len(want) {
		var gotNames, wantNames []string
		for _, c := range got {
			gotNames = append(gotNames, c.typeName)
		}
		for _, c := range want {
			wantNames = append(wantNames, c.typeName)
		}
		t.Errorf("unequal lengths: %d != %d; got %v, got %v", len(got), len(want), gotNames, wantNames)
		return
	}

	for i := range want {
		// don't bother with checking ast positions.
		// zero out these values.
		for k := range got[i].members.NameToPos {
			got[i].members.NameToPos[k] = 0
		}
		equalCheckEnum(t, got[i], want[i])
	}
}
