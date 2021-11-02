package exhaustive

import (
	"go/ast"
	"go/types"
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// TODO: write tests that assert on the "result" returned by nodeVisitor.

func TestDiagnosticEnumTypeName(t *testing.T) {}

func TestDiagnosticMissingMembers(t *testing.T) {}

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
		types.NewTypeName(1, types.NewPackage("example.org/enumpkg", "enumpkg"), "Biome", nil),
		nil, /* underlying type should not matter */
		nil,
	)
	allMembers := &enumMembers{Names: []string{"Tundra", "Savanna", "Desert"}}
	missingMembers := []string{"Savanna", "Desert"}
	strategy := byValue

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
