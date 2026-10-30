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

type e1[E any] interface {
	f() E
}

type e2[E any] interface {
	~int
	f() E
}

func f12[P any, Q e1[P], R e2[P]]() {
	var x any
	switch x.(type) {
	case P: // valid
	case Q: // valid
	case R: // valid
	case e1[P]: // valid (e1 is a basic interface)
	case e2[P]: // error: cannot use type e2[P] outside a type constraint: interface contains type constraints
	}
}
