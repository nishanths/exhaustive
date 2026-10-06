package def

type t1 int

const (
	x0 t1 = iota
	x1
)

func f1() {
	var t t1

	switch t {
	case x0:
	case x1:
	}

	switch t {
	case x0:
	case x1:
	default:
	}

	switch t { // want "^expression switch not exhaustive: missing cases: x0$"
	case x1:
	}

	switch t {
	case x1:
	default:
	}

	switch t { // want "^expression switch not exhaustive: missing cases: x0, x1$"
	}

	switch t {
	default:
	}
}
