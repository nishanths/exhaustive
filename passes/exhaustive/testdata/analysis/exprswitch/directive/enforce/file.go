package enforce

type t1 int

const (
	x0 t1 = iota
	x1
)

func f1() {
	var t t1

	switch t {
	case x0:
	}

	//exhaustive:ignore
	switch t {
	case x0:
	}

	//exhaustive:enforce
	switch t { // want "^missing cases in expression switch: x1$"
	case x0:
	}
}

func g1() {
	_ = map[t1]bool{x1: true}

	//exhaustive:ignore
	_ = map[t1]bool{x1: true}

	//exhaustive:enforce
	_ = map[t1]bool{x0: true} // want "^missing keys in map literal: x1$"

	// For map literals, the comments of the ancestor
	// assignment statement are considered. Therefore in the
	// following test case expect both the outer map and the
	// inner map to be analyzed.
	//

	//exhaustive:enforce
	_ = map[t1]map[t1]bool{ // want "^missing keys in map literal: x1$"
		x0: { // want "^missing keys in map literal: x0$"
			x1: true,
		},
	}
}
