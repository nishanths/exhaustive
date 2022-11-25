//go:build !go1.18
// +build !go1.18

package exhaustive

func composingEnumTypes(pass *analysis.Pass, t types.Type) ([]typeAndMembers, bool) {
	switch t := t.(type) {
	case *types.Named:
		return composingEnumTypesNamed(pass, t)
	default:
		return nil, false
	}
}
