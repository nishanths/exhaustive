package exhaustive

import (
	"fmt"
	"go/token"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

// checks that an enumMembers literal is correctly defined in tests.
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

		if want, got := []string{"foo", "z", "bar", "y", "x"}, v.Names; !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}

		if want, got := map[string]token.Pos{
			"foo": 10,
			"z":   20,
			"bar": 30,
			"y":   40,
			"x":   50,
		}, v.NameToPos; !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}

		if want, got := map[string]constantValue{
			"foo": "\"A\"",
			"z":   "X",
			"bar": "\"B\"",
			"y":   "Y",
			"x":   "X",
		}, v.NameToValue; !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}

		if want, got := map[constantValue][]string{
			"\"A\"": {"foo"},
			"\"B\"": {"bar"},
			"X":     {"z", "x"},
			"Y":     {"y"},
		}, v.ValueToNames; !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	// TODO(testing): add tests for iota, repeated values, ...
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
		t.Run("package scopes only "+strconv.FormatBool(pkgOnly), func(t *testing.T) {
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

func equalCheckEnum(t *testing.T, want, got checkEnum) {
	t.Helper()
	if want.typeName != got.typeName {
		t.Errorf("type name: want %s, got %s", want.typeName, got.typeName)
	}
	if !reflect.DeepEqual(want.members, got.members) {
		t.Errorf("type name %s: members: want %+v, got %+v", want.typeName, want.members, got.members)
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

	if len(want) != len(got) {
		var wantNames, gotNames []string
		for _, c := range want {
			wantNames = append(wantNames, c.typeName)
		}
		for _, c := range got {
			gotNames = append(gotNames, c.typeName)
		}
		t.Errorf("unequal lengths: %d != %d; want %v, got %v", len(want), len(got), wantNames, gotNames)
		return
	}

	for i := range want {
		// don't bother with checking ast positions.
		// zero out the values.
		for k := range got[i].members.NameToPos {
			got[i].members.NameToPos[k] = 0
		}
		equalCheckEnum(t, want[i], got[i])
	}
}
