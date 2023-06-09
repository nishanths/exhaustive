package exhaustive

import (
	"go/ast"
	"go/types"
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// TODO: write tests that assert on the "result" returned by
// switchStmtChecker.

// This test mainly exists to ensure stability of the diagnostic message
// format.
func TestMakeSwitchDiagnostic(t *testing.T) {
	sw := &ast.SwitchStmt{
		Switch: 1,
		Body: &ast.BlockStmt{
			Rbrace: 10,
		},
		// other fields shouldn't matter
	}
	tn := types.NewTypeName(50, types.NewPackage("example.org/enumpkg", "enumpkg"), "Biome", nil)
	et := enumType{tn}
	missing := map[member]struct{}{
		{102, et, "Savanna", "2"}: {},
		{109, et, "Desert", "3"}:  {},
	}

	got := makeSwitchDiagnostic(sw, []enumType{et}, missing)
	want := analysis.Diagnostic{
		Pos:     1,
		End:     11,
		Message: "missing cases in switch of type enumpkg.Biome: enumpkg.Savanna, enumpkg.Desert",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
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

		var got []constantValue
		gotDefaultExists := analyzeSwitchClauses(sw, info, func(val constantValue) {
			got = append(got, val)
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
		if gotDefaultExists != wantDefaultExists {
			t.Errorf("got %v, want %v", gotDefaultExists, wantDefaultExists)
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
				t.Errorf("func name: got %q, want %q", getFuncName(fn), tt.funcName)
				return
			}
			sw := getSwitchStatement(fn)
			assertFoundNames(t, sw, tt.pkg.TypesInfo, tt.vals, tt.defaultExists)
		})
	}
}
