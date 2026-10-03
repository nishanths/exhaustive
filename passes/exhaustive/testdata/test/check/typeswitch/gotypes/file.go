package gotypes

import (
	"go/types"
)

func f1() bool {
	var t types.Type

	switch t.(type) {
	case *types.Array,
		*types.Basic,
		*types.Chan,
		*types.Map,
		*types.Pointer,
		*types.Signature,
		*types.Slice,
		*types.Struct,
		*types.Tuple,
		*types.TypeParam,
		*types.Union,
		*types.Interface:
		return false
	case *types.Alias:
		return true
	case *types.Named:
		return true
	default:
		return false
	}
}

func f2() bool {
	var t types.Type

	switch t.(type) { // want "^type switch not exhaustive: missing cases: \\*types.Basic, \\*types.Pointer$"
	case *types.Array,
		*types.Chan,
		*types.Map,
		*types.Signature,
		*types.Slice,
		*types.Struct,
		*types.Tuple,
		*types.TypeParam,
		*types.Union,
		*types.Interface:
		return false
	case *types.Alias:
		return true
	case *types.Named:
		return true
	default:
		return false
	}
}
