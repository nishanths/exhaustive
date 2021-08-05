package exhaustive

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestContainsFuncCall(t *testing.T) {
	cfg := &packages.Config{Mode: packages.NeedTypesInfo | packages.NeedTypes | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, "./testdata/funccall")
	if !checkNoError(t, err) {
		return
	}

	pkg := pkgs[0]
	syn := pkg.Syntax[0]
	decl := syn.Decls[len(syn.Decls)-1].(*ast.GenDecl)

	checkEqualf(t, 8, len(decl.Specs), "wrong number of specs (either test or testdata file needs update?)")

	for idx, spec := range decl.Specs {
		spec := spec.(*ast.ValueSpec)
		want, err := strconv.ParseBool(strings.Trim(spec.Comment.Text(), "\r\n"))
		if !checkNoError(t, err) {
			// testdata file doesn't have comment in right format?
			continue
		}
		got := containsFuncCall(pkg.TypesInfo, spec.Values[0])
		checkEqualf(t, want, got, "index %d", idx)
	}
}

func TestHasImportWithPath(t *testing.T) {
	t.Run("does not have", func(t *testing.T) {
		got := hasImportWithPath([]*ast.ImportSpec{
			{Path: &ast.BasicLit{Value: `"foo/bar"`}},
			{Path: &ast.BasicLit{Value: `"x/y"`}},
			{Path: &ast.BasicLit{Value: `"github.com/foo"`}},
		}, `"foo"`)
		checkEqual(t, false, got)
	})

	t.Run("has", func(t *testing.T) {
		got := hasImportWithPath([]*ast.ImportSpec{
			{Path: &ast.BasicLit{Value: `"foo/bar"`}},
			{Path: &ast.BasicLit{Value: `"x/y"`}},
			{Path: &ast.BasicLit{Value: `"github.com/foo"`}},
		}, `"x/y"`)
		checkEqual(t, true, got)
	})
}

func TestFlattenImportSpec(t *testing.T) {
	in := [][]*ast.ImportSpec{{
		{Path: &ast.BasicLit{Value: `"foo/bar"`}},
		{Path: &ast.BasicLit{Value: `"x/y"`}},
		{Path: &ast.BasicLit{Value: `"github.com/foo"`}},
	}, {
		{Path: &ast.BasicLit{Value: `"golang.org/x/net"`}},
	}, {
		{Path: &ast.BasicLit{Value: `"github.com/example"`}},
	}}

	got := flattenImportSpec(in)

	checkEqual(t, []*ast.ImportSpec{
		{Path: &ast.BasicLit{Value: `"foo/bar"`}},
		{Path: &ast.BasicLit{Value: `"x/y"`}},
		{Path: &ast.BasicLit{Value: `"github.com/foo"`}},
		{Path: &ast.BasicLit{Value: `"golang.org/x/net"`}},
		{Path: &ast.BasicLit{Value: `"github.com/example"`}},
	}, got)
}

func TestFirstImportDecl(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		const source = `
package foo

import "fmt"

import ( "bytes" )
`
		f, err := parser.ParseFile(token.NewFileSet(), "", source, parser.AllErrors)
		if !checkNoError(t, err) {
			return
		}
		decl := firstImportDecl(f)
		checkEqual(t, `"fmt"`, decl.Specs[0].(*ast.ImportSpec).Path.Value)
	})

	t.Run("none", func(t *testing.T) {
		const source = `package foo`
		f, err := parser.ParseFile(token.NewFileSet(), "", source, parser.AllErrors)
		if !checkNoError(t, err) {
			return
		}
		decl := firstImportDecl(f)
		checkNil(t, decl)
	})
}
