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

type m interface{ m() }

type T struct{}
func (*T) m() {}

type A = T

type U struct{}
func (*U) m() {}

func f4(v m) {
	switch v.(type) {
	case *A:
	case *U:
	}

	switch v.(type) { // want "^missing cases in type switch: \\*T or \\*A, \\*U$"
	}
}

type Node interface{ node() }
type Statement interface{ Node; statement() }

type Literal struct{}
func (*Literal) node() {}

type Loop struct{}
func (*Loop) node() {}
func (*Loop) statement() {}

func f5(v Node) {
	switch v.(type) {
	case Statement:
	case *Literal:
	case *Loop: // want "^unreachable type switch case \\*Loop \\(matched by earlier case Statement\\)$"
	}
}
