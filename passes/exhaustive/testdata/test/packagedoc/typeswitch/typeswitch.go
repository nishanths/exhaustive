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

type i interface{ m() }

type t struct{}
func (*t) m() {}

type a = t

func f4(v i) {
	switch v.(type) {
	case *a: // or *t
	}

	switch v.(type) { // want "^missing cases in type switch: \\*t or \\*a$"
	}
}
