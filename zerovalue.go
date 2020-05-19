package exhaustive

import (
	"go/constant"
	"go/types"
)

func zeroValue(b *types.Basic) constant.Value {
	switch i := b.Info(); {
	case i&types.IsInteger != 0:
		return constant.MakeInt64(0)
	case i&types.IsFloat != 0:
		return constant.MakeFloat64(0.0)
	case i&types.IsString != 0:
		return constant.MakeString("")
	case i&types.IsBoolean != 0:
		return constant.MakeBool(false)
	default:
		panic("unhandled TODO")
	}
}
