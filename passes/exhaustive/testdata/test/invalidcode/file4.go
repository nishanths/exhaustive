package invalidcode

func f6() {
	type n interface{}
	type s struct{ *n } // invalid: embedded field type cannot be a pointer to an interface
}

func f7() {
	type im interface{ m() }
	type in interface{ n() }
	var x im

	switch x.(type) {
	case im: // valid
	case *im: // invalid: x (variable of interface type im) cannot have dynamic type *im (type *im is pointer to interface, not interface)
	case **im: // invalid: x (variable of interface type im) cannot have dynamic type **im (missing method m)
	case in: // valid
	case *in: // invalid: x (variable of interface type im) cannot have dynamic type *in (type *im is pointer to interface, not interface)
	case **in: // invalid: x (variable of interface type im) cannot have dynamic type **in  (missing method m)
	}

	var y interface{}
	switch y.(type) {
	case im: // valid
	case *im: // valid
	case **im: // valid
	}
}
