package directive

type t struct{}

func (*t) m()

type u struct{}

func (*u) m()

func f1() {
	var x interface{ m() }

	//exhaustive:badname
	switch x.(type) { // want "^error parsing comment directives: invalid name \"badname\"$"
	}

	//exhaustive:enforce
	//exhaustive:ignore
	switch x.(type) { // want "^error parsing comment directives: conflicting directives$"
	}

	switch x.(type) { // want "^missing cases in type switch: \\*u$"
	case *t:
	}

	//exhaustive:ignore
	switch x.(type) {
	case *t:
	}

	//exhaustive:enforce
	switch x.(type) { // want "^missing cases in type switch: \\*u$"
	case *t:
	}
}

func e1() {
	var x interface{}

	switch x.(type) {
	}

	//exhaustive:ignore
	switch x.(type) {
	}

	//exhaustive:enforce
	switch x.(type) { // want "^enforce directive present and type switch not checked: unsupported: interface type of switch expression is the empty interface$"
	}
}

type myError struct{}

func (*myError) Error() string { return "" }

func e2() {
	var err1 error

	switch err1.(type) {
	}

	//exhaustive:ignore
	switch err1.(type) {
	}

	//exhaustive:enforce
	switch err1.(type) { // want "^missing cases in type switch: \\*myError$"
	}
}
