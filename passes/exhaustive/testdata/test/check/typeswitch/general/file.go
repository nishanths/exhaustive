package general

func f1() {
	var t types.Type
	
	switch t.(type) {
	case *types.Alias:
	case *types.Named:
	}
}