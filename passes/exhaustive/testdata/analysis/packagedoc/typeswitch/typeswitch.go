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
