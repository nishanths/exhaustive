//go:build !go1.18
// +build !go1.18

package exhaustive

func fromNamed(pass *analysis.Pass, t *types.Named) (typeAndMembers, bool) {
	if tpkg := t.Obj().Pkg(); tpkg == nil {
		return typeAndMembers{}, false
	}

	et := enumType{t.Obj()}
	em, ok := importFact(pass, et)
	if !ok {
		return typeAndMembers{}, false
	}

	return typeAndMembers{et, em}, true
}

func composingEnumTypes(pass *analysis.Pass, t types.Type) ([]typeAndMembers, bool) {
	switch t := t.(type) {
	case *types.Named:
		e, ok := fromNamed(pass, t)
		if !ok {
			return nil, false
		}
		return []typeAndMembers{e}, true
	default:
		return nil, false
	}
}
