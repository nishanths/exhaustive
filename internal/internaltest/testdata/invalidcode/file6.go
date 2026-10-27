package invalidcode

func f11() {
	type i interface{ i() }

	var x0 any
	switch x0.(type) {
	case *i: // valid
	}

	var x1 i
	switch x1.(type) {
	case *i: // error: x1 (variable of interface type i) cannot have dynamic type *i (type *i is pointer to interface, not interface)
	}

	type j interface{ j() }
	var x2 j
	switch x2.(type) {
	case *i: // error: x2 (variable of interface type j) cannot have dynamic type *i (type *i is pointer to interface, not interface)
	}
}
