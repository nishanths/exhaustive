package exhaustive

import (
	"go/ast"
	"go/token"
	"go/types"

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
	for _, spec := range specs {
		if spec.Path.Value == pathLiteral {
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

// copies an GenDecl in a manner such that appending to the returned GenDecl's Specs field
// doesn't mutate the original GenDecl
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
