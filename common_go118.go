//go:build go1.18
// +build go1.18

package exhaustive

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func fromNamed(pass *analysis.Pass, t *types.Named, typeparam bool) (result []typeAndMembers, all bool) {
	if tpkg := t.Obj().Pkg(); tpkg == nil {
		// go/types documentation says: nil for labels and
		// objects in the Universe scope. This happens for the built-in
		// error type for example.
		return nil, true
	}

	et := enumType{t.Obj()}
	if em, ok := importFact(pass, et); ok {
		return []typeAndMembers{{et, em}}, true
	}

	if typeparam {
		if intf, ok := t.Underlying().(*types.Interface); ok {
			return fromInterface(pass, intf, typeparam)
		}
	}

	return nil, true
}

func fromInterface(pass *analysis.Pass, intf *types.Interface, typeparam bool) (result []typeAndMembers, all bool) {
	all = true

	for i := 0; i < intf.NumEmbeddeds(); i++ {
		embed := intf.EmbeddedType(i)
		switch embed.(type) {
		case *types.Union:
			u := embed.(*types.Union)
			for i := 0; i < u.Len(); i++ {
				r, a := fromType(pass, u.Term(i).Type(), typeparam)
				result = append(result, r...)
				all = all && a
			}
		case *types.Named:
			r, a := fromNamed(pass, embed.(*types.Named), typeparam)
			result = append(result, r...)
			all = all && a
		default:
			// don't care about these.
			// e.g. basic type
		}
	}

	return
}

func fromType(pass *analysis.Pass, t types.Type, typeparam bool) (result []typeAndMembers, all bool) {
	switch t := t.(type) {
	case *types.Named:
		return fromNamed(pass, t, typeparam)

	case *types.TypeParam:
		// does not appear to be explicitly documented, but based on
		// spec (see section Type constraints) and source code, we can
		// expect constraints to have underlying type *types.Interface.
		intf := t.Constraint().Underlying().(*types.Interface)
		return fromInterface(pass, intf, typeparam)

	default:
		// ignore these.
		return nil, true
	}
}

func composingEnumTypes(pass *analysis.Pass, t types.Type) (result []typeAndMembers, all bool) {
	_, typeparam := t.(*types.TypeParam)
	return fromType(pass, t, typeparam)
}
