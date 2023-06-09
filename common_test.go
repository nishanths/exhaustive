package exhaustive

import (
	"go/token"
	"go/types"
	"reflect"
	"regexp"
	"testing"
)

func TestChecklist(t *testing.T) {
	et := enumType{types.NewTypeName(50, types.NewPackage("github.com/example/bar-go", "bar"), "T", nil)}
	em := enumMembers{
		Names: []string{"A", "B", "C", "D", "E", "F", "G"},
		NameToPos: map[string]token.Pos{
			"A": 0,
			"B": 0,
			"C": 0,
			"D": 0,
			"E": 0,
			"F": 0,
			"G": 0,
		},
		NameToValue: map[string]constantValue{
			"A": "1",
			"B": "2",
			"C": "5",
			"D": "2",
			"E": "3",
			"F": "2",
			"G": "4",
		},
		ValueToNames: map[constantValue][]string{
			"1": {"A"},
			"2": {"B", "D", "F"},
			"3": {"E"},
			"4": {"G"},
			"5": {"C"},
		},
	}
	checkEnumMembersLiteral("TestChecklist", em)

	checkRemaining := func(t *testing.T, h checklist, want map[string]struct{}) {
		t.Helper()
		rem := make(map[string]struct{})
		for k := range h.remaining() {
			rem[k.name] = struct{}{}
		}
		if !reflect.DeepEqual(rem, want) {
			t.Errorf("got %+v, want %+v", rem, want)
		}
	}

	t.Run("main operations", func(t *testing.T) {
		var c checklist
		c.add(et, em, false)
		checkRemaining(t, c, map[string]struct{}{
			"A": {},
			"B": {},
			"C": {},
			"D": {},
			"E": {},
			"F": {},
			"G": {},
		})

		c.found(`1`)
		checkRemaining(t, c, map[string]struct{}{
			"B": {},
			"C": {},
			"D": {},
			"E": {},
			"F": {},
			"G": {},
		})

		c.found(`2`)
		checkRemaining(t, c, map[string]struct{}{
			"C": {},
			"E": {},
			"G": {},
		})

		// repeated call should be a no-op.
		c.found(`2`)
		checkRemaining(t, c, map[string]struct{}{
			"C": {},
			"E": {},
			"G": {},
		})

		c.found(`2`)
		checkRemaining(t, c, map[string]struct{}{
			"C": {},
			"E": {},
			"G": {},
		})

		c.found(`5`)
		checkRemaining(t, c, map[string]struct{}{
			"E": {},
			"G": {},
		})

		// unknown value
		c.found(`100000`)
		checkRemaining(t, c, map[string]struct{}{
			"E": {},
			"G": {},
		})

		c.found(`3`)
		checkRemaining(t, c, map[string]struct{}{
			"G": {},
		})
	})

	t.Run("ignore pattern", func(t *testing.T) {
		t.Run("none", func(t *testing.T) {
			var c checklist
			c.add(et, em, false)
			checkRemaining(t, c, map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
				"D": {},
				"E": {},
				"F": {},
				"G": {},
			})
		})

		t.Run("constant", func(t *testing.T) {
			t.Run("basic", func(t *testing.T) {
				var c checklist
				c.ignoreConstant(regexp.MustCompile(`^github\.com/example/bar-go\.G$`))
				c.add(et, em, false)
				checkRemaining(t, c, map[string]struct{}{
					"A": {},
					"B": {},
					"C": {},
					"D": {},
					"E": {},
					"F": {},
				})
			})

			t.Run("matches multiple", func(t *testing.T) {
				var c checklist
				c.ignoreConstant(regexp.MustCompile(`^github\.com/example/bar-go`))
				c.add(et, em, false)
				checkRemaining(t, c, map[string]struct{}{})
			})

			t.Run("uses package path, not package name", func(t *testing.T) {
				var c checklist
				c.ignoreConstant(regexp.MustCompile(`bar\.G`)) // this should not cause anything to be ignored.
				c.add(et, em, false)
				checkRemaining(t, c, map[string]struct{}{
					"A": {},
					"B": {},
					"C": {},
					"D": {},
					"E": {},
					"F": {},
					"G": {},
				})
			})
		})

		t.Run("type", func(t *testing.T) {
			t.Run("basic", func(t *testing.T) {
				var c checklist
				c.ignoreType(regexp.MustCompile(`^github\.com/example/bar-go\.T$`))
				c.add(et, em, false)
				checkRemaining(t, c, map[string]struct{}{})
			})

			t.Run("uses package path, not package name", func(t *testing.T) {
				var c checklist
				c.ignoreType(regexp.MustCompile(`bar\.T`)) // this should not cause anything to be ignored.
				c.add(et, em, false)
				checkRemaining(t, c, map[string]struct{}{
					"A": {},
					"B": {},
					"C": {},
					"D": {},
					"E": {},
					"F": {},
					"G": {},
				})
			})
		})
	})

	t.Run("blank identifier", func(t *testing.T) {
		em := enumMembers{
			Names: []string{"A", "B", "C", "D", "E", "F", "G", "_"},
			NameToPos: map[string]token.Pos{
				"A": 0,
				"B": 0,
				"C": 0,
				"D": 0,
				"E": 0,
				"F": 0,
				"G": 0,
				"_": 0,
			},
			NameToValue: map[string]constantValue{
				"A": "1",
				"B": "2",
				"C": "5",
				"D": "2",
				"E": "3",
				"F": "2",
				"G": "4",
				"_": "0",
			},
			ValueToNames: map[constantValue][]string{
				"0": {"_"},
				"1": {"A"},
				"2": {"B", "D", "F"},
				"3": {"E"},
				"4": {"G"},
				"5": {"C"},
			},
		}
		checkEnumMembersLiteral("TestChecklist blank identifier", em)

		var c checklist
		c.add(et, em, true)
		checkRemaining(t, c, map[string]struct{}{
			"A": {},
			"B": {},
			"C": {},
			"D": {},
			"E": {},
			"F": {},
			"G": {},
		})
	})

	t.Run("unexported", func(t *testing.T) {
		em := enumMembers{
			Names: []string{"A", "B", "C", "D", "E", "F", "G", "lowercase"},
			NameToPos: map[string]token.Pos{
				"A":         0,
				"B":         0,
				"C":         0,
				"D":         0,
				"E":         0,
				"F":         0,
				"G":         0,
				"lowercase": 0,
			},
			NameToValue: map[string]constantValue{
				"A":         "1",
				"B":         "2",
				"C":         "5",
				"D":         "2",
				"E":         "3",
				"F":         "2",
				"G":         "4",
				"lowercase": "42",
			},
			ValueToNames: map[constantValue][]string{
				"1":  {"A"},
				"2":  {"B", "D", "F"},
				"3":  {"E"},
				"4":  {"G"},
				"5":  {"C"},
				"42": {"lowercase"},
			},
		}
		checkEnumMembersLiteral("TestChecklist lowercase", em)

		t.Run("include", func(t *testing.T) {
			var c checklist
			c.add(et, em, true)
			checkRemaining(t, c, map[string]struct{}{
				"A":         {},
				"B":         {},
				"C":         {},
				"D":         {},
				"E":         {},
				"F":         {},
				"G":         {},
				"lowercase": {},
			})
		})

		t.Run("don't include", func(t *testing.T) {
			var c checklist
			c.add(et, em, false)
			checkRemaining(t, c, map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
				"D": {},
				"E": {},
				"F": {},
				"G": {},
			})
		})
	})
}

func TestDiagnosticEnumType(t *testing.T) {
	tn := types.NewTypeName(50, types.NewPackage("example.org/enumpkg-go", "enumpkg"), "Biome", nil)
	got := diagnosticEnumType(tn)
	want := "enumpkg.Biome"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGroupMissing(t *testing.T) {
	groupStrings := func(groups []group) [][]string {
		var out [][]string
		for i := range groups {
			var x []string
			for j := range groups[i] {
				x = append(x, diagnosticMember(groups[i][j]))
			}
			out = append(out, x)
		}
		return out
	}

	// f adapts groupMissing for easy use in the test.
	f := func(missing map[member]struct{}, types []enumType) [][]string {
		return groupStrings(groupify(missing, types))
	}

	tn := types.NewTypeName(50, types.NewPackage("example.org/enumpkg-go", "enumpkg"), "River", nil)
	et := enumType{tn}

	members := []member{
		0: {10, et, "Ganga", "0"},
		1: {20, et, "Yamuna", "2"},
		2: {30, et, "Kaveri", "1"},
		3: {60, et, "Unspecified", "0"},
	}

	t.Run("missing some: same-valued", func(t *testing.T) {
		got := f(map[member]struct{}{
			members[0]: {},
			members[3]: {},
			members[2]: {},
		}, []enumType{et})
		want := [][]string{{"enumpkg.Ganga", "enumpkg.Unspecified"}, {"enumpkg.Kaveri"}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("missing some: unique or unknown values", func(t *testing.T) {
		got := f(map[member]struct{}{
			members[1]: {},
			members[2]: {},
		}, []enumType{et})
		want := [][]string{{"enumpkg.Yamuna"}, {"enumpkg.Kaveri"}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("missing none", func(t *testing.T) {
		got := f(nil, []enumType{et})
		if len(got) != 0 {
			t.Errorf("got %d elements, want 0 elements", len(got))
		}
	})

	t.Run("missing all", func(t *testing.T) {
		got := f(map[member]struct{}{
			members[0]: {},
			members[2]: {},
			members[1]: {},
			members[3]: {},
		}, []enumType{et})
		want := [][]string{{"enumpkg.Ganga", "enumpkg.Unspecified"}, {"enumpkg.Yamuna"}, {"enumpkg.Kaveri"}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	tn = types.NewTypeName(50, types.NewPackage("example.org/xkcd-go", "xkcd"), "T", nil)
	et = enumType{tn}
	members = []member{
		0: {12, et, "X", "0"},
		1: {13, et, "A", "1"},
		2: {14, et, "Unspecified", "0"},
	}

	t.Run("AST order", func(t *testing.T) {
		got := f(map[member]struct{}{
			members[2]: {},
			members[0]: {},
			members[1]: {},
		}, []enumType{et})
		want := [][]string{{"xkcd.X", "xkcd.Unspecified"}, {"xkcd.A"}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
