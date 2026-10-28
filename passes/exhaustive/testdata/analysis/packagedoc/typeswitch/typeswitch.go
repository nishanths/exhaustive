package foo

import "go/ast"
import "go/types"

func f3(t types.Type) bool {
	switch t.(type) { // want "^missing cases in type switch: \\*types.Basic, \\*types.TypeParam$"
	case *types.Array,
		*types.Chan,
		*types.Interface,
		*types.Map,
		*types.Pointer,
		*types.Signature,
		*types.Slice,
		*types.Struct,
		*types.Tuple,
		*types.Union:
		return false
	case *types.Alias, *types.Named:
		return true
	default:
		panic("internal error: unhandled type")
	}
}

func count(file *ast.File) {
	var decls, fields, funcs int

	for n := range ast.Preorder(file) {
		switch n.(type) { // want "^missing cases in type switch: \\*ast\\..+$"
		case ast.Decl:
			decls += 1
		case *ast.Field:
			fields += 1
		case *ast.FuncDecl: // want "^ineffectual type \\*ast.FuncDecl due to case order: dynamic type that would be matched by type \\*ast.FuncDecl is matched in earlier case by type ast.Decl$"
			funcs += 1
		}
	}

	/* ... */
}

type I interface{ f6example() }
type T struct{}

func (*T) f6example() {}

type A = T

func f6(x I) {
	switch x.(type) {
	case *T:
	}

	switch x.(type) {
	case *A:
	}
}
