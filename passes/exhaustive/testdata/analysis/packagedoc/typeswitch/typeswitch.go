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

type U struct{}
func (*U) m() {}

type A = T

func f4(v I) {
	switch v.(type) {
	case *A: // note that *A is identical to *T
	case *U:
	}

	switch v.(type) { // want "^missing cases in type switch: \\*T or \\*A, \\*U$"
	}
}

type node interface{ node() }
type stmt interface{ node; stmt() }

type loop struct{}
func (*loop) node() {}
func (*loop) stmt() {}

type literal struct{}
func (*literal) node() {}

func f5(v node) {
	switch v.(type) {
	case stmt:
	case *loop: // want "^unreachable type switch case \\*loop \\(matched by earlier case stmt\\)$"
	case *literal:
	}
}
