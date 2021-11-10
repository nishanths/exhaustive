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
		tn := types.NewTypeName(50, types.NewPackage("example.org/enumpkg-go", "enumpkg"), "Biome", nil)
		got := diagnosticEnumTypeName(tn, true)
		want := "Biome"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})

	t.Run("different package", func(t *testing.T) {
		tn := types.NewTypeName(50, types.NewPackage("example.org/enumpkg-go", "enumpkg"), "Biome", nil)
		got := diagnosticEnumTypeName(tn, false)
		want := "enumpkg.Biome"
		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})
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
	tn := types.NewTypeName(50, types.NewPackage("example.org/enumpkg", "enumpkg"), "Biome", nil)
	enumTyp := enumType{tn}
	allMembers := enumMembers{
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
	m := map[string]constantValue{
		"Tundra":  "1",
		"Savanna": "2",
		"Desert":  "3",
	}

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

	assertFoundNames := func(t *testing.T, sw *ast.SwitchStmt, info *types.Info, want []constantValue, wantDefaultExists bool) {
		t.Helper()
		tagType := info.Types[sw.Tag].Type.(*types.Named)

		var got []constantValue
		gotDefaultExists := analyzeSwitchClauses(sw, tagType.Obj().Pkg(), m, info, func(val constantValue) {
			got = append(got, val)
		})

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
		if wantDefaultExists != gotDefaultExists {
			t.Errorf("want %v, got %v", wantDefaultExists, gotDefaultExists)
		}
	}

	type testSpec struct {
		// decl index of function
		declIdx int

		// which package and file to look in
		pkg  *packages.Package
		file *ast.File

		// what to expect at the function
		funcName      string
		vals          []constantValue
		defaultExists bool
	}

	cases := []testSpec{
		{1, switchtest, switchtestGoFile, "switchWithDefault", []constantValue{m["Tundra"], m["Desert"]}, true},
		{2, switchtest, switchtestGoFile, "switchWithoutDefault", []constantValue{m["Tundra"], m["Desert"]}, false},
		{3, switchtest, switchtestGoFile, "switchParen", []constantValue{m["Tundra"], m["Desert"]}, false},
		{4, switchtest, switchtestGoFile, "switchNotIdent", []constantValue{m["Savanna"]}, false},

		{1, otherpkg, otherpkgGoFile, "switchParen", []constantValue{m["Tundra"], m["Desert"]}, false},
		{2, otherpkg, otherpkgGoFile, "switchNotSelExpr", []constantValue{m["Tundra"]}, false},
		{4, otherpkg, otherpkgGoFile, "switchNotExpectedSelExpr", []constantValue{m["Desert"]}, false},
	}

	for _, tt := range cases {
		t.Run(tt.pkg.Name+"#"+tt.funcName, func(t *testing.T) {
			fn := tt.file.Decls[tt.declIdx]
			if getFuncName(fn) != tt.funcName {
				t.Errorf("want func name %q, got %q", tt.funcName, getFuncName(fn))
				return
			}
			sw := getSwitchStatement(fn)
			assertFoundNames(t, sw, tt.pkg.TypesInfo, tt.vals, tt.defaultExists)
		})
	}
}
