package exhaustive

import (
	"go/types"
	"reflect"
	"regexp"
	"testing"
)

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
	checkEnumMembersLiteral(t, "TestChecklist", em)

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
		checkEnumMembersLiteral(t, "TestChecklist blank identifier", em)

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
		checkEnumMembersLiteral(t, "TestChecklist lowercase", em)

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
