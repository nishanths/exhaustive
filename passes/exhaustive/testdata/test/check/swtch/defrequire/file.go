package defrequire

type t1 int

const (
	x0 t1 = iota
	x1
)

func f1() {
	var t t1

	switch t { // want "^missing default case$" "^expression switch not exhaustive: missing cases: x1$"
	case x0:
	}

	switch t { // want "^missing default case$"
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

	//exhaustive:defrequire=0
	switch t {
	case x0:
	case x1:
	}

	//exhaustive:ignore-default-case-required
	switch t {
	case x0:
	case x1:
	}

	//exhaustive:defrequire=1
	switch t { // want "^missing default case$"
	case x0:
	case x1:
	}

	//exhaustive:enforce-default-case-required
	switch t { // want "^missing default case$"
	case x0:
	case x1:
	}
}

func f3() {
	var t t1

	//exhaustive:defrequire=0
	//exhaustive:defrequire=1
	switch t { // want "^error parsing comment directives: conflicting directives$"
	case x0:
	case x1:
	}

	//exhaustive:ignore-default-case-required
	//exhaustive:enforce-default-case-required
	switch t { // want "^error parsing comment directives: conflicting directives$"
	case x0:
	case x1:
	}
}

func f4() {
	var t t1

	//exhaustive:ignore
	//exhaustive:defrequire=1
	switch t {
	case x0:
	case x1:
	}

	//exhaustive:ignore
	//exhaustive:enforce-default-case-required
	switch t {
	case x0:
	case x1:
	}
}

func f5() {
	var v int

	switch {
	}

	switch v {
	}
}
