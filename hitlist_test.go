package exhaustive

import (
	"fmt"
	"go/types"
	"reflect"
	"regexp"
	"testing"
)

func TestHitlist(t *testing.T) {
	enumPkg := types.NewPackage("github.com/example/bar-go", "bar")

	em := &enumMembers{
		Names: []string{"A", "B", "C", "D", "E", "F", "G"},
		NameToValue: map[string]string{
			"A": "1",
			"B": "2",
			// C has no AST value
			"D": "2",
			"E": "3",
			"F": "2",
			"G": "4",
		},
		ValueToNames: map[string][]string{
			"1": {"A"},
			"2": {"B", "D", "F"},
			"3": {"E"},
			"4": {"G"},
		},
	}

	checkRemaining := func(t *testing.T, h *hitlist, want map[string]struct{}) {
		t.Helper()
		rem := h.remaining()
		if !reflect.DeepEqual(want, rem) {
			t.Errorf("want %+v, got %+v", want, rem)
		}
	}

	t.Run("panics on unknown strategy", func(t *testing.T) {
		hitlist := makeHitlist(em, enumPkg, false, nil)
		f := func() {
			hitlist.found("A", hitlistStrategy(8238))
		}
		assertPanic(t, f, fmt.Sprintf("unknown strategy %v", hitlistStrategy(8238)))
	})

	t.Run("main operations", func(t *testing.T) {
		hitlist := makeHitlist(em, enumPkg, false, nil)
		checkRemaining(t, hitlist, map[string]struct{}{
			"A": {},
			"B": {},
			"C": {},
			"D": {},
			"E": {},
			"F": {},
			"G": {},
		})

		hitlist.found("A", byValue)
		checkRemaining(t, hitlist, map[string]struct{}{
			"B": {},
			"C": {},
			"D": {},
			"E": {},
			"F": {},
			"G": {},
		})

		hitlist.found("B", byName)
		checkRemaining(t, hitlist, map[string]struct{}{
			"C": {},
			"D": {},
			"E": {},
			"F": {},
			"G": {},
		})

		// repeated call should be a no-op.
		hitlist.found("B", byName)
		checkRemaining(t, hitlist, map[string]struct{}{
			"C": {},
			"D": {},
			"E": {},
			"F": {},
			"G": {},
		})

		hitlist.found("F", byValue)
		checkRemaining(t, hitlist, map[string]struct{}{
			"C": {},
			"E": {},
			"G": {},
		})

		hitlist.found("C", byValue)
		checkRemaining(t, hitlist, map[string]struct{}{
			"E": {},
			"G": {},
		})

		hitlist.found("E", byName)
		checkRemaining(t, hitlist, map[string]struct{}{
			"G": {},
		})
	})

	t.Run("ignore regexp", func(t *testing.T) {
		t.Run("nil means no filtering", func(t *testing.T) {
			hitlist := makeHitlist(em, enumPkg, false, nil)
			checkRemaining(t, hitlist, map[string]struct{}{
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
			hitlist := makeHitlist(em, enumPkg, false, regexp.MustCompile(`^github.com/example/bar-go.G$`))
			checkRemaining(t, hitlist, map[string]struct{}{
				"A": {},
				"B": {},
				"C": {},
				"D": {},
				"E": {},
				"F": {},
			})
		})

		t.Run("uses package path, not package name", func(t *testing.T) {
			hitlist := makeHitlist(em, enumPkg, false, regexp.MustCompile(`bar.G`))
			checkRemaining(t, hitlist, map[string]struct{}{
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
		em := *em
		em.Names = append([]string{}, em.Names...)
		em.Names = append(em.Names, "_")

		hitlist := makeHitlist(&em, enumPkg, true, nil)
		checkRemaining(t, hitlist, map[string]struct{}{
			"A": {},
			"B": {},
			"C": {},
			"D": {},
			"E": {},
			"F": {},
			"G": {},
		})
	})

	t.Run("include unexported", func(t *testing.T) {
		em := *em
		em.Names = append([]string{}, em.Names...)
		em.Names = append(em.Names, "lowercase")

		t.Run("include", func(t *testing.T) {
			hitlist := makeHitlist(&em, enumPkg, true, nil)
			checkRemaining(t, hitlist, map[string]struct{}{
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
			hitlist := makeHitlist(&em, enumPkg, false, nil)
			checkRemaining(t, hitlist, map[string]struct{}{
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
