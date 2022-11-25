//go:build go1.18
// +build go1.18

package exhaustive

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func composingEnumTypes(pass *analysis.Pass, t types.Type) ([]typeAndMembers, bool) {
	switch t := t.(type) {
	case *types.Named:
		return composingEnumTypesNamed(pass, t)

	case *types.TypeParam:
		// TODO(nishanths) handle type param

		// log.Printf("%T", tagType.Constraint())
		// log.Printf("%#v", tagType.Constraint())
		// log.Printf("%#v", tagType.Obj())
		// log.Printf("%#v", tagType.Constraint().Underlying().(*types.Interface))
		// log.Printf("%#v", tagType.Obj().TypeParams())
		// log.Printf("%#v", tagType.Obj().TypeArgs())
		return nil, false

	default:
		return nil, false
	}
}
