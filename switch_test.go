package exhaustive

import (
	"go/ast"
	"go/types"
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// TODO(testing): write tests that assert on the "result" returned by switchStmtChecker.

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
	em := &enumMembers{
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
	checkEnumMembersLiteral(t, "River", em)

	t.Run("missing some: same-valued", func(t *testing.T) {
		got := diagnosticMissingMembers([]string{"Ganga", "Unspecified", "Kaveri"}, em)
		want := []string{"Ganga|Unspecified", "Kaveri"}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("missing some: unique or unknown values", func(t *testing.T) {
		got := diagnosticMissingMembers([]string{"Yamuna", "Kaveri"}, em)
		want := []string{"Kaveri", "Yamuna"}
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
		got := diagnosticMissingMembers([]string{"Ganga", "Kaveri", "Yamuna", "Unspecified"}, em)
		want := []string{"Ganga|Unspecified", "Kaveri", "Yamuna"}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
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
	named := types.NewNamed(
		types.NewTypeName(50, types.NewPackage("example.org/enumpkg", "enumpkg"), "Biome", nil),
		nil, /* underlying type should not matter */
		nil,
	)
	enumTyp := enumType{named}
	allMembers := &enumMembers{
		Names: []string{"Tundra", "Savanna", "Desert"},
		NameToValue: map[string]constantValue{
			"Tundra":  "1",
			"Savanna": "2",
			"Desert":  "3",
		},
		ValueToNames: map[constantValue][]string{
			"1": {"Tundra"},
			"2": {"Savanna"},
			"3": {"Desert"},
		},
	}
	checkEnumMembersLiteral(t, "Biome", allMembers)
	missingMembers := []string{"Savanna", "Desert"}

	got := makeDiagnostic(sw, samePkg, enumTyp, allMembers, missingMembers)
	want := analysis.Diagnostic{
		Pos:     1,
		End:     11,
		Message: "missing cases in switch of type enumpkg.Biome: Desert, Savanna",
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %+v, got %+v", want, got)
	}
}

func TestAnalyzeSwitchClauses(t *testing.T) {
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedTypesInfo | packages.NeedTypes | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, "./testdata/switchtest/...")
	assertNoError(t, err)

	switchtest, otherpkg := pkgs[0], pkgs[1]
	switchtestGoFile, otherpkgGoFile := switchtest.Syntax[1], otherpkg.Syntax[0]

	getFuncName := func(fn ast.Decl) string {
		funcDecl := fn.(*ast.FuncDecl)
		return funcDecl.Name.Name
	}

	getSwitchStatement := func(fn ast.Decl) *ast.SwitchStmt {
		// in this testdata, the switch statement is always the first statement
		// in the function body.
		funcDecl := fn.(*ast.FuncDecl)
		return funcDecl.Body.List[0].(*ast.SwitchStmt)
	}

	assertFoundNames := func(t *testing.T, sw *ast.SwitchStmt, typesInfo *types.Info, samePkg bool, wantNames []string, wantDefaultExists bool) {
		t.Helper()
		var gotNames []string
		gotDefaultExists := analyzeSwitchClauses(sw, typesInfo, samePkg, func(name string) {
			gotNames = append(gotNames, name)
		})
		if !reflect.DeepEqual(wantNames, gotNames) {
			t.Errorf("want %v, got %v", wantNames, gotNames)
		}
		if wantDefaultExists != gotDefaultExists {
			t.Errorf("want %v, got %v", wantDefaultExists, gotDefaultExists)
		}
	}

	type testSpec struct {
		declIdx int // func decl index

		samePkg bool
		pkg     *packages.Package
		file    *ast.File

		// what to expect at the decl index:
		funcName      string
		memberNames   []string
		defaultExists bool
	}

	cases := []testSpec{
		{1, true, switchtest, switchtestGoFile, "switchWithDefault", []string{"Tundra", "Desert"}, true},
		{2, true, switchtest, switchtestGoFile, "switchWithoutDefault", []string{"Tundra", "Desert"}, false},
		{3, true, switchtest, switchtestGoFile, "switchParen", []string{"Tundra", "Desert"}, false},
		{4, true, switchtest, switchtestGoFile, "switchNotIdent", []string{"Savanna"}, false},

		{1, false, otherpkg, otherpkgGoFile, "switchParen", []string{"Tundra", "Desert"}, false},
		{2, false, otherpkg, otherpkgGoFile, "switchNotSelExpr", []string{"Tundra"}, false},
		{4, false, otherpkg, otherpkgGoFile, "switchNotExpectedSelExpr", []string{"Desert"}, false},
	}

	for _, tt := range cases {
		t.Run(tt.pkg.Name+"#"+tt.funcName, func(t *testing.T) {
			fn := tt.file.Decls[tt.declIdx]
			if getFuncName(fn) != tt.funcName {
				t.Errorf("want func name %q, got %q", tt.funcName, getFuncName(fn))
				return
			}
			sw := getSwitchStatement(fn)
			assertFoundNames(t, sw, tt.pkg.TypesInfo, tt.samePkg, tt.memberNames, tt.defaultExists)
		})
	}
}
