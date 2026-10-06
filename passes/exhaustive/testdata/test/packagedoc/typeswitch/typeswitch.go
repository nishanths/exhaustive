package foo

import "go/types"

func f3(t types.Type) bool {
	switch t.(type) { // want "^type switch not exhaustive: missing cases: \\*types.Basic, \\*types.TypeParam$"
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

type S struct{}
func (*S) m()
type A = S

func f4() {
	var x interface{ m() }
	switch x.(type) {
	case *A: // or *S
	}

	switch x.(type) { // want "^type switch not exhaustive: missing cases: \\*S or \\*A$"
	}
}
