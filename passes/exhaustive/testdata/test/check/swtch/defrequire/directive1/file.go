package defrequire_directive1

type t1 int

const (
	x0 t1 = iota
	x1
)

func f1() {
	var t t1

	switch t { // want "^expression switch not exhaustive: missing cases: x1$"
	case x0:
	}

	switch t {
	case x0:
	case x1:
	}

	switch t {
	case x0:
	case x1:
	default:
	}
}

func f2() {
	var t t1

	//exhaustive:defrequire=1
	switch t { // want "^missing default case$" "^expression switch not exhaustive: missing cases: x1$"
	case x0:
	}

	//exhaustive:defrequire=1
	switch t { // want "^missing default case$"
	case x0:
	case x1:
	}

	//exhaustive:defrequire=1
	switch t {
	case x0:
	case x1:
	default:
	}

	// Deprecated spellings of the same comment directives as above:

	//exhaustive:enforce-default-case-required
	switch t { // want "^missing default case$" "^expression switch not exhaustive: missing cases: x1$"
	case x0:
	}

	//exhaustive:enforce-default-case-required
	switch t { // want "^missing default case$"
	case x0:
	case x1:
	}

	//exhaustive:enforce-default-case-required
	switch t {
	case x0:
	case x1:
	default:
	}
}
