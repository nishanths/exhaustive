package exhaustive

import (
	"go/ast"
	"go/types"
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// TODO: write tests that assert on the "result" returned by nodeVisitor.

func TestDiagnosticEnumTypeName(t *testing.T) {
	t.Run("same package", func(t *testing.T) {
		enumType := types.NewNamed(
			types.NewTypeName(50, types.NewPackage("example.org/enumpkg-go", "enumpkg"), "Biome", nil),
			nil, /* underlying type should not matter */
			nil,
		)
		got := diagnosticEnumTypeName(enumType, true)
		want := "Biome"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})

	t.Run("different package", func(t *testing.T) {
		enumType := types.NewNamed(
			types.NewTypeName(50, types.NewPackage("example.org/enumpkg-go", "enumpkg"), "Biome", nil),
			nil, /* underlying type should not matter */
			nil,
		)
		got := diagnosticEnumTypeName(enumType, false)
		want := "enumpkg.Biome"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})
}

func TestDiagnosticMissingMembers(t *testing.T) {
	t.Run("strategy: value", func(t *testing.T) {
		strategy := strategyValue
		em := &enumMembers{
			Names: []string{"Ganga", "Yamuna", "Kaveri", "Unspecified"},
			NameToValue: map[string]string{
				"Unspecified": "0",
				"Ganga":       "0",
			},
			ValueToNames: map[string][]string{
				"0": {"Unspecified", "Ganga"},
			},
		}

		t.Run("missing some: same-valued", func(t *testing.T) {
			got := diagnosticMissingMembers([]string{"Ganga", "Unspecified", "Kaveri"}, em, strategy)
			want := []string{"Ganga|Unspecified", "Kaveri"}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v, got %v", want, got)
			}
		})

		t.Run("missing some: all unique values", func(t *testing.T) {
			got := diagnosticMissingMembers([]string{"Yamuna", "Kaveri"}, em, strategy)
			want := []string{"Kaveri", "Yamuna"}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v, got %v", want, got)
			}
		})

		t.Run("missing none", func(t *testing.T) {
			got := diagnosticMissingMembers(nil, em, strategy)
			if len(got) != 0 {
				t.Errorf("want zero elements, got %d", len(got))
			}
		})

		t.Run("missing all", func(t *testing.T) {
			got := diagnosticMissingMembers([]string{"Ganga", "Kaveri", "Yamuna", "Unspecified"}, em, strategy)
			want := []string{"Ganga|Unspecified", "Kaveri", "Yamuna"}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v, got %v", want, got)
			}
		})
	})

	t.Run("strategy: name", func(t *testing.T) {
		strategy := strategyName
		em := &enumMembers{
			Names: []string{"Ganga", "Yamuna", "Kaveri", "Unspecified"},
			NameToValue: map[string]string{
				"Unspecified": "0",
				"Ganga":       "0",
			},
			ValueToNames: map[string][]string{
				"0": {"Unspecified", "Ganga"},
			},
		}

		t.Run("missing some: same-valued", func(t *testing.T) {
			got := diagnosticMissingMembers([]string{"Ganga", "Unspecified", "Kaveri"}, em, strategy)
			want := []string{"Ganga", "Kaveri", "Unspecified"}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v, got %v", want, got)
			}
		})

		t.Run("missing some: all unique values", func(t *testing.T) {
			got := diagnosticMissingMembers([]string{"Yamuna", "Kaveri"}, em, strategy)
			want := []string{"Kaveri", "Yamuna"}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v, got %v", want, got)
			}
		})

		t.Run("missing none", func(t *testing.T) {
			got := diagnosticMissingMembers(nil, em, strategy)
			if len(got) != 0 {
				t.Errorf("want zero elements, got %d", len(got))
			}
		})

		t.Run("missing all", func(t *testing.T) {
			got := diagnosticMissingMembers([]string{"Ganga", "Kaveri", "Yamuna", "Unspecified"}, em, strategy)
			want := []string{"Ganga", "Kaveri", "Unspecified", "Yamuna"}
			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v, got %v", want, got)
			}
		})
	})
}

// This test mainly exists to ensure stability of the diagnostic message format.
func TestMakeDiagnostic(t *testing.T) {
	sw := &ast.SwitchStmt{
		Switch: 1,
		Body: &ast.BlockStmt{
			Rbrace: 10,
		},
		// other fields shouldn't matter
	}
	samePkg := false
	enumType := types.NewNamed(
		types.NewTypeName(50, types.NewPackage("example.org/enumpkg", "enumpkg"), "Biome", nil),
		nil, /* underlying type should not matter */
		nil,
	)
	allMembers := &enumMembers{Names: []string{"Tundra", "Savanna", "Desert"}}
	missingMembers := []string{"Savanna", "Desert"}
	strategy := strategyValue

	got := makeDiagnostic(sw, samePkg, enumType, allMembers, missingMembers, strategy)
	want := analysis.Diagnostic{
		Pos:     1,
		End:     11,
		Message: "missing cases in switch of type enumpkg.Biome: Desert, Savanna",
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %+v, got %+v", want, got)
	}
}

func TestAnalyzeCaseClauseExpr(t *testing.T) {}

func TestAnalyzeSwitchClauses(t *testing.T) {}
