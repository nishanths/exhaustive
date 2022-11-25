package exhaustive

import (
	"go/types"
	"reflect"
	"regexp"
	"testing"
)

func TestRegexpFlag(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		var v regexpFlag
		if got := v.regexp(); got != nil {
			t.Errorf("want nil, got %+v", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set(""); err != nil {
			t.Errorf("error unexpectedly non-nil: %v", err)
		}
		if got := v.regexp(); got != nil {
			t.Errorf("want nil, got %+v", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("bad input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set("("); err == nil {
			t.Errorf("error unexpectedly nil")
		}
		if got := v.regexp(); got != nil {
			t.Errorf("want nil, got %+v", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("good input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set("^foo$"); err != nil {
			t.Errorf("error unexpectedly non-nil: %v", err)
		}
		if v.regexp() == nil {
			t.Errorf("unexpectedly nil")
		}
		if !v.regexp().MatchString("foo") {
			t.Errorf("did not match")
		}
		if got, want := v.String(), regexp.MustCompile("^foo$").String(); got != want {
			t.Errorf("want %q, got %q", got, want)
		}
	})

	// The flag.Value interface doc says: "The flag package may call the
	// String method with a zero-valued receiver, such as a nil pointer."
	t.Run("String() nil receiver", func(t *testing.T) {
		var v *regexpFlag
		// expect no panic, and ...
		if got := v.String(); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})
}

func TestChecklist(t *testing.T) {
	enumPkg := types.NewPackage("github.com/example/bar-go", "bar")

	em := enumMembers{
		Names: []string{"A", "B", "C", "D", "E", "F", "G"},
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

	checkRemaining := func(t *testing.T, h *checklist, want map[string]struct{}) {
		t.Helper()
		rem := h.remaining()
		if !reflect.DeepEqual(want, rem) {
			t.Errorf("want %+v, got %+v", want, rem)
		}
	}

	t.Run("main operations", func(t *testing.T) {
		checklist := makeChecklist(em, enumPkg, false, nil)
		checkRemaining(t, checklist, map[string]struct{}{
			"A": {},
			"B": {},
			"C": {},
			"D": {},
			"E": {},
			"F": {},
			"G": {},
		})

		checklist.found(`1`)
		checkRemaining(t, checklist, map[string]struct{}{
			"B": {},
			"C": {},
			"D": {},
			"E": {},
			"F": {},
			"G": {},
		})

		checklist.found(`2`)
		checkRemaining(t, checklist, map[string]struct{}{
			"C": {},
			"E": {},
			"G": {},
		})

		// repeated call should be a no-op.
		checklist.found(`2`)
		checkRemaining(t, checklist, map[string]struct{}{
			"C": {},
			"E": {},
			"G": {},
		})

		checklist.found(`2`)
		checkRemaining(t, checklist, map[string]struct{}{
			"C": {},
			"E": {},
			"G": {},
		})

		checklist.found(`5`)
		checkRemaining(t, checklist, map[string]struct{}{
			"E": {},
			"G": {},
		})

		// unknown value
		checklist.found(`100000`)
		checkRemaining(t, checklist, map[string]struct{}{
			"E": {},
			"G": {},
		})

		checklist.found(`3`)
		checkRemaining(t, checklist, map[string]struct{}{
			"G": {},
		})
	})

	t.Run("ignore regexp", func(t *testing.T) {
		t.Run("nil means no filtering", func(t *testing.T) {
			checklist := makeChecklist(em, enumPkg, false, nil)
			checkRemaining(t, checklist, map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
				"D": {},
				"E": {},
				"F": {},
				"G": {},
			})
		})

		t.Run("basic", func(t *testing.T) {
			checklist := makeChecklist(em, enumPkg, false, regexp.MustCompile(`^github.com/example/bar-go.G$`))
			checkRemaining(t, checklist, map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
				"D": {},
				"E": {},
				"F": {},
			})
		})

		t.Run("matches multiple", func(t *testing.T) {
			checklist := makeChecklist(em, enumPkg, false, regexp.MustCompile(`^github.com/example/bar-go`))
			checkRemaining(t, checklist, map[string]struct{}{})
		})

		t.Run("uses package path, not package name", func(t *testing.T) {
			checklist := makeChecklist(em, enumPkg, false, regexp.MustCompile(`bar.G`))
			checkRemaining(t, checklist, map[string]struct{}{
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

	t.Run("blank identifier", func(t *testing.T) {
		em := enumMembers{
			Names: []string{"A", "B", "C", "D", "E", "F", "G", "_"},
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

		checklist := makeChecklist(em, enumPkg, true, nil)
		checkRemaining(t, checklist, map[string]struct{}{
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
			checklist := makeChecklist(em, enumPkg, true, nil)
			checkRemaining(t, checklist, map[string]struct{}{
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
			checklist := makeChecklist(em, enumPkg, false, nil)
			checkRemaining(t, checklist, map[string]struct{}{
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
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestDiagnosticMissingMembers(t *testing.T) {
	em := enumMembers{
		Names: []string{"Ganga", "Yamuna", "Kaveri", "Unspecified"},
		NameToValue: map[string]constantValue{
			"Unspecified": "0",
			"Ganga":       "0",
			"Kaveri":      "1",
			"Yamuna":      "2",
		},
		ValueToNames: map[constantValue][]string{
			"0": {"Unspecified", "Ganga"},
			"1": {"Kaveri"},
			"2": {"Yamuna"},
		},
	}
	checkEnumMembersLiteral("River", em)

	t.Run("missing some: same-valued", func(t *testing.T) {
		got := diagnosticMissingMembers(map[string]struct{}{"Ganga": {}, "Unspecified": {}, "Kaveri": {}}, em)
		want := []string{"Ganga|Unspecified", "Kaveri"}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("missing some: unique or unknown values", func(t *testing.T) {
		got := diagnosticMissingMembers(map[string]struct{}{"Yamuna": {}, "Kaveri": {}}, em)
		want := []string{"Yamuna", "Kaveri"}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("missing none", func(t *testing.T) {
		got := diagnosticMissingMembers(nil, em)
		if len(got) != 0 {
			t.Errorf("want zero elements, got %d", len(got))
		}
	})

	t.Run("missing all", func(t *testing.T) {
		got := diagnosticMissingMembers(map[string]struct{}{"Ganga": {}, "Kaveri": {}, "Yamuna": {}, "Unspecified": {}}, em)
		want := []string{"Ganga|Unspecified", "Yamuna", "Kaveri"}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	em = enumMembers{
		Names: []string{"X", "A", "Unspecified"},
		NameToValue: map[string]constantValue{
			"Unspecified": "0",
			"X":           "0",
			"A":           "1",
		},
		ValueToNames: map[constantValue][]string{
			"0": {"Unspecified", "X"},
			"1": {"A"},
		},
	}
	checkEnumMembersLiteral("whatever", em)

	t.Run("AST order", func(t *testing.T) {
		got := diagnosticMissingMembers(map[string]struct{}{"Unspecified": {}, "X": {}, "A": {}}, em)
		want := []string{"X|Unspecified", "A"}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}
