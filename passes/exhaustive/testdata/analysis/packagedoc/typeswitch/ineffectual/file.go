package ineffectual

import "go/ast"

func count(file *ast.File) {
	var decls, fields, funcs int
	for n := range ast.Preorder(file) {
		switch n.(type) {
		case ast.Decl:
			decls += 1
		case *ast.Field:
			fields += 1
		case *ast.FuncDecl: // want "^ineffectual: \\*ast.FuncDecl -> ast.Decl$"
			funcs += 1
		default:
			// do not care about other types
		}
	}
	/* ... */
}
