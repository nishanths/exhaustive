package foo

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

type I interface{ m() }

type T struct{}
func (*T) m() {}

type A = T

func f4() {
	var v I

	switch v.(type) {
	case *A: // or *T
	}

	switch v.(type) { // want "^missing cases in type switch: \\*T or \\*A$"
	}
}
