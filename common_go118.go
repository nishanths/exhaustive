//go:build go1.18
// +build go1.18

package exhaustive

import (
	"go/types"
	"log"

	"golang.org/x/tools/go/analysis"
)

func fromNamed(pass *analysis.Pass, t *types.Named, typeparam bool) ([]typeAndMembers, bool) {
	if tpkg := t.Obj().Pkg(); tpkg == nil {
		// The Go documentation says: nil for labels and objects in
		// the Universe scope. This happens for the built-in error
		// type for example.
		return nil, false
	}

	et := enumType{t.Obj()}
	em, ok := importFact(pass, et)
	if !ok {
		// type is not a known enum type.
		return nil, false
	}

	return []typeAndMembers{{et, em}}, true
}

func composingEnumTypes(pass *analysis.Pass, t types.Type) ([]typeAndMembers, bool) {
	switch t := t.(type) {
	case *types.Named:
		return fromNamed(pass, t, false)

	case *types.TypeParam:
		intf, ok := t.Constraint().Underlying().(*types.Interface)
		if !ok {
			panic("expect constraint to have underlying type *types.Interface")
		}
		log.Printf("---")
		log.Printf("%#v", intf)
		for i := 0; i < intf.NumEmbeddeds(); i++ {
			log.Printf("%#v", intf.EmbeddedType(i))
		}

		// log.Printf("%T", t.Constraint())
		// log.Printf("%#v", t.Constraint())
		// log.Printf("%#v", t.Obj())
		// log.Printf("%#v", t.Constraint().Underlying().(*types.Interface))
		// log.Printf("%#v", tagType.Obj().TypeParams())
		// log.Printf("%#v", tagType.Obj().TypeArgs())
		return nil, false

	default:
		return nil, false
	}
}
