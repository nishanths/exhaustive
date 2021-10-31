package exhaustive

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
)

func containsFuncCall(typesInfo *types.Info, e ast.Expr) bool {
	e = astutil.Unparen(e)
	c, ok := e.(*ast.CallExpr)
	if !ok {
		return false
	}
	if _, isFunc := typesInfo.TypeOf(c.Fun).Underlying().(*types.Signature); isFunc {
		return true
	}
	for _, a := range c.Args {
		if containsFuncCall(typesInfo, a) {
			return true
		}
	}
	return false
}

func hasImportWithPath(specs []*ast.ImportSpec, pathLiteral string) bool {
	in, err := strconv.Unquote(pathLiteral)
	if err != nil {
		panic("strconv.Unquote(" + pathLiteral + "): " + err.Error())
	}

	for _, spec := range specs {
		s, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			panic("strconv.Unquote(" + spec.Path.Value + "): " + err.Error())
		}
		if s == in {
			return true
		}
	}
	return false
}

func flattenImportSpec(iss [][]*ast.ImportSpec) []*ast.ImportSpec {
	var out []*ast.ImportSpec
	for _, is := range iss {
		for _, spec := range is {
			out = append(out, spec)
		}
	}
	return out
}

// copies a GenDecl in a manner such that appending to the returned GenDecl's Specs field
// doesn't affect the original GenDecl
func copyGenDecl(im *ast.GenDecl) *ast.GenDecl {
	imCopy := *im
	imCopy.Specs = make([]ast.Spec, len(im.Specs))
	for i := range im.Specs {
		imCopy.Specs[i] = im.Specs[i]
	}
	return &imCopy
}

func firstImportDecl(f *ast.File) *ast.GenDecl {
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok && genDecl.Tok == token.IMPORT {
			// first IMPORT GenDecl
			return genDecl
		}
	}
	return nil
}

// Returns a TextEdit that adds "fmt" import to the file.
func fmtImportTextEdit(fset *token.FileSet, f *ast.File) analysis.TextEdit {
	firstDecl := firstImportDecl(f)

	if firstDecl == nil {
		// file has no import declarations
		// insert "fmt" import spec after package statement
		return analysis.TextEdit{
			Pos: f.Name.End() + 1, // end of package name + 1
			End: f.Name.End() + 1,
			NewText: []byte(`import (
	"fmt"
)`),
		}
	}

	// copy because we'll be mutating its Specs field
	firstDeclCopy := copyGenDecl(firstDecl)

	// find insertion index for "fmt" import spec
	var i int
	for ; i < len(firstDeclCopy.Specs); i++ {
		im := firstDeclCopy.Specs[i].(*ast.ImportSpec)
		if v, _ := strconv.Unquote(im.Path.Value); v > "fmt" {
			break
		}
	}

	// insert "fmt" import spec at the index
	fmtSpec := &ast.ImportSpec{
		Path: &ast.BasicLit{
			// NOTE: Pos field doesn't seem to be required for our
			// purposes here.
			Kind:  token.STRING,
			Value: `"fmt"`,
		},
	}
	s := firstDeclCopy.Specs // local var for easier comprehension of next line
	s = append(s[:i], append([]ast.Spec{fmtSpec}, s[i:]...)...)
	firstDeclCopy.Specs = s

	// create the text edit
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, firstDeclCopy)
	return analysis.TextEdit{
		Pos:     firstDecl.Pos(),
		End:     firstDecl.End(),
		NewText: buf.Bytes(),
	}
}
