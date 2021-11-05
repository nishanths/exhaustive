package exhaustive

import (
	"reflect"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestEnumMembers_add(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var v enumMembers
		v.add("foo", nil)
		v.add("z", ptrString("X"))
		v.add("bar", nil)
		v.add("y", ptrString("Y"))
		v.add("x", ptrString("X"))

		if want, got := []string{"foo", "z", "bar", "y", "x"}, v.Names; !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
		if want, got := map[string]string{
			"z": "X",
			"y": "Y",
			"x": "X",
		}, v.NameToValue; !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}

		if want, got := map[string][]string{
			"X": {"z", "x"},
			"Y": {"y"},
		}, v.ValueToNames; !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	// TODO: add tests for iota, repeated values, ...
}

var enumpkg = func() *packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypesInfo | packages.NeedTypes | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, "./testdata/src/enum")
	if err != nil {
		panic(err)
	}
	return pkgs[0]
}()

func TestFindPossibleEnumTypes(t *testing.T) {
	var got []string
	findPossibleEnumTypes(enumpkg.Syntax, enumpkg.TypesInfo, func(name string) {
		got = append(got, name)
	})
	want := []string{
		"VarMembers",
		"IotaEnum",
		"MemberlessEnum",
		"RepeatedValue",
		"AcrossBlocksDeclsFiles",
		"UnexportedMembers",
		"NonTopLevel",
		"ParenVal",
		"UIntEnum",
		"StringEnum",
		"RuneEnum",
		"ByteEnum",
		"Int32Enum",
		"Float64Enum",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\nwant %v\ngot  %v", want, got)
		return
	}
}

func TestFindEnumMembers(t *testing.T) {
	knownEnumTypes := make(map[string]struct{})
	findPossibleEnumTypes(enumpkg.Syntax, enumpkg.TypesInfo, func(name string) {
		knownEnumTypes[name] = struct{}{}
	})

	got := make(map[string]*enumMembers)
	findEnumMembers(enumpkg.Syntax, enumpkg.TypesInfo, knownEnumTypes, func(memberName, typeName string, constVal *string) {
		if _, ok := got[typeName]; !ok {
			got[typeName] = &enumMembers{}
		}
		got[typeName].add(memberName, constVal)
	})

	checkEnums(t, got)
}

func TestFindEnums(t *testing.T) {
	result := findEnums(enumpkg.Syntax, enumpkg.TypesInfo)
	checkEnums(t, result)
}

// shared utility for TestFindEnumMembers and TestFindEnums.
func checkEnums(t *testing.T, got map[string]*enumMembers) {
	t.Helper()

	want := enums{
		"VarMembers": {
			[]string{"VarMemberA"},
			nil,
			nil,
		},
		"IotaEnum": {
			[]string{"IotaA", "IotaB"},
			map[string]string{
				"IotaA": `2`,
			},
			map[string][]string{
				`2`: {"IotaA"},
			},
		},
		"RepeatedValue": {
			[]string{"RepeatedValueA", "RepeatedValueB"},
			map[string]string{
				"RepeatedValueA": `1`,
				"RepeatedValueB": `1`,
			},
			map[string][]string{
				`1`: {"RepeatedValueA", "RepeatedValueB"},
			},
		},
		"AcrossBlocksDeclsFiles": {
			[]string{"Here", "Separate", "There"},
			map[string]string{
				"Here":     `0`,
				"Separate": `1`,
				"There":    `2`,
			},
			map[string][]string{
				`0`: {"Here"},
				`1`: {"Separate"},
				`2`: {"There"},
			},
		},
		"UnexportedMembers": {
			[]string{"unexportedMembersA", "unexportedMembersB"},
			map[string]string{
				"unexportedMembersA": `1`,
				"unexportedMembersB": `2`,
			},
			map[string][]string{
				`1`: {"unexportedMembersA"},
				`2`: {"unexportedMembersB"},
			},
		},
		"ParenVal": {
			[]string{"ParenVal0", "ParenVal1"},
			map[string]string{
				"ParenVal0": `0`,
				"ParenVal1": `1`,
			},
			map[string][]string{
				`0`: {"ParenVal0"},
				`1`: {"ParenVal1"},
			},
		},
		"UIntEnum": {
			[]string{"UIntA", "UIntB"},
			map[string]string{
				"UIntA": "0",
				"UIntB": "1",
			},
			map[string][]string{
				"0": {"UIntA"},
				"1": {"UIntB"},
			},
		},
		"StringEnum": {
			[]string{"StringA", "StringB", "StringC"},
			map[string]string{
				"StringA": `"stringa"`,
				"StringB": `"stringb"`,
				"StringC": `"stringc"`,
			},
			map[string][]string{
				`"stringa"`: {"StringA"},
				`"stringb"`: {"StringB"},
				`"stringc"`: {"StringC"},
			},
		},
		"RuneEnum": {
			[]string{"RuneA"},
			map[string]string{
				"RuneA": `97`,
			},
			map[string][]string{
				`97`: {"RuneA"},
			},
		},
		"ByteEnum": {
			[]string{"ByteA"},
			map[string]string{
				"ByteA": `97`,
			},
			map[string][]string{
				`97`: {"ByteA"},
			},
		},
		"Int32Enum": {
			[]string{"Int32A", "Int32B"},
			map[string]string{
				"Int32A": "0",
				"Int32B": "1",
			},
			map[string][]string{
				"0": {"Int32A"},
				"1": {"Int32B"},
			},
		},
		"Float64Enum": {
			[]string{"Float64A", "Float64B"},
			map[string]string{
				"Float64A": `1`,
			},
			map[string][]string{
				`1`: {"Float64A"},
			},
		},
	}

	if len(want) != len(got) {
		t.Errorf("unequal lengths: want %d, got %d", len(want), len(got))
		return
	}

	// check members for each type.
	for k := range want {
		if !reflect.DeepEqual(want[k], got[k]) {
			t.Errorf("%s: want %v, got %v", k, *want[k], *got[k])
		}
	}
}
