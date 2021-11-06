package exhaustive

import (
	"go/types"
	"reflect"
	"testing"

	"golang.org/x/tools/go/packages"
)

// Asserts that the enumMembers literal is correctly defined.
func checkEnumMembersLiteral(t *testing.T, id string, v *enumMembers) {
	t.Helper()

	if len(v.Names) != len(v.NameToValue) {
		t.Fatalf("%s: wrong lengths: %d != %d  (test definition bug)", id, len(v.Names), len(v.NameToValue))
	}

	var count int
	for _, names := range v.ValueToNames {
		count += len(names)
	}
	if len(v.Names) != count {
		t.Fatalf("%s: wrong lengths: %d != %d (test definition bug)", id, len(v.Names), count)
	}
}

func TestEnumMembers_add(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var v enumMembers
		v.add("foo", "\"A\"")
		v.add("z", "X")
		v.add("bar", "\"B\"")
		v.add("y", "Y")
		v.add("x", "X")

		if want, got := []string{"foo", "z", "bar", "y", "x"}, v.Names; !reflect.DeepEqual(want, got) {
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

var testdataEnumPkg = func() *packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypesInfo | packages.NeedTypes | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, "./testdata/src/enum")
	if err != nil {
		panic(err)
	}
	return pkgs[0]
}()

func TestFindPossibleEnumTypes(t *testing.T) {
	var got []string
	findPossibleEnumTypes(testdataEnumPkg.Syntax, testdataEnumPkg.TypesInfo, func(named *types.Named) {
		got = append(got, named.Obj().Name())
	})
	want := []string{
		"VarMember",
		"VarConstMixed",
		"IotaEnum",
		"MemberlessEnum",
		"RepeatedValue",
		"AcrossBlocksDeclsFiles",
		"UnexportedMembers",
		"NonTopLevel",
		"ParenVal",
		"T",
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
	possibleEnumTypes := make(map[*types.Named]struct{})
	findPossibleEnumTypes(testdataEnumPkg.Syntax, testdataEnumPkg.TypesInfo, func(named *types.Named) {
		possibleEnumTypes[named] = struct{}{}
	})

	got := make(map[string]*enumMembers)
	findEnumMembers(testdataEnumPkg.Syntax, testdataEnumPkg.TypesInfo, possibleEnumTypes, func(memberName string, enumTyp enumType, val constantValue) {
		if _, ok := got[enumTyp.name()]; !ok {
			got[enumTyp.name()] = &enumMembers{}
		}
		got[enumTyp.name()].add(memberName, val)
	})

	checkEnums(t, got)
}

func TestFindEnums(t *testing.T) {
	result := findEnums(testdataEnumPkg.Syntax, testdataEnumPkg.TypesInfo)

	transformForChecking := func(in map[enumType]*enumMembers) map[string]*enumMembers {
		out := make(map[string]*enumMembers)
		for typ, mem := range in {
			out[typ.name()] = mem
		}
		return out
	}

	checkEnums(t, transformForChecking(result))
}

// shared utility for TestFindEnumMembers and TestFindEnums.
func checkEnums(t *testing.T, got map[string]*enumMembers) {
	t.Helper()

	want := map[string]*enumMembers{
		"VarConstMixed": {
			[]string{"VCMixedB"},
			map[string]constantValue{
				"VCMixedB": `1`,
			},
			map[constantValue][]string{
				`1`: {"VCMixedB"},
			},
		},
		"IotaEnum": {
			[]string{"IotaA", "IotaB"},
			map[string]constantValue{
				"IotaA": `0`,
				"IotaB": `2`,
			},
			map[constantValue][]string{
				`0`: {"IotaA"},
				`2`: {"IotaB"},
			},
		},
		"RepeatedValue": {
			[]string{"RepeatedValueA", "RepeatedValueB"},
			map[string]constantValue{
				"RepeatedValueA": `1`,
				"RepeatedValueB": `1`,
			},
			map[constantValue][]string{
				`1`: {"RepeatedValueA", "RepeatedValueB"},
			},
		},
		"AcrossBlocksDeclsFiles": {
			[]string{"Here", "Separate", "There"},
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
		},
		"UnexportedMembers": {
			[]string{"unexportedMembersA", "unexportedMembersB"},
			map[string]constantValue{
				"unexportedMembersA": `1`,
				"unexportedMembersB": `2`,
			},
			map[constantValue][]string{
				`1`: {"unexportedMembersA"},
				`2`: {"unexportedMembersB"},
			},
		},
		"ParenVal": {
			[]string{"ParenVal0", "ParenVal1"},
			map[string]constantValue{
				"ParenVal0": `0`,
				"ParenVal1": `1`,
			},
			map[constantValue][]string{
				`0`: {"ParenVal0"},
				`1`: {"ParenVal1"},
			},
		},
		"T": {
			[]string{"A", "B"},
			map[string]constantValue{
				"A": `0`,
				"B": `1`,
			},
			map[constantValue][]string{
				`0`: {"A"},
				`1`: {"B"},
			},
		},
		"UIntEnum": {
			[]string{"UIntA", "UIntB"},
			map[string]constantValue{
				"UIntA": "0",
				"UIntB": "1",
			},
			map[constantValue][]string{
				"0": {"UIntA"},
				"1": {"UIntB"},
			},
		},
		"StringEnum": {
			[]string{"StringA", "StringB", "StringC"},
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
		},
		"RuneEnum": {
			[]string{"RuneA"},
			map[string]constantValue{
				"RuneA": `97`,
			},
			map[constantValue][]string{
				`97`: {"RuneA"},
			},
		},
		"ByteEnum": {
			[]string{"ByteA"},
			map[string]constantValue{
				"ByteA": `97`,
			},
			map[constantValue][]string{
				`97`: {"ByteA"},
			},
		},
		"Int32Enum": {
			[]string{"Int32A", "Int32B"},
			map[string]constantValue{
				"Int32A": "0",
				"Int32B": "1",
			},
			map[constantValue][]string{
				"0": {"Int32A"},
				"1": {"Int32B"},
			},
		},
		"Float64Enum": {
			[]string{"Float64A", "Float64B"},
			map[string]constantValue{
				"Float64A": `0`,
				"Float64B": `1`,
			},
			map[constantValue][]string{
				`0`: {"Float64A"},
				`1`: {"Float64B"},
			},
		},
	}

	// check the `want` declaration for programmer error.
	for k, v := range want {
		checkEnumMembersLiteral(t, k, v)
	}

	if len(want) != len(got) {
		t.Errorf("unequal lengths: want %d, got %d", len(want), len(got))
		return
	}

	// check members for each type.
	for k := range want {
		g, ok := got[k]
		if !ok {
			t.Errorf("missing %q in got", k)
			return
		}
		if !reflect.DeepEqual(want[k], g) {
			t.Errorf("%s: want %v, got %v", k, *want[k], *g)
		}
	}
}
