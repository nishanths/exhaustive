package def

type t struct{}
func (t)  m()

type u struct{}
func (*u) m()

func f1() {
	var x interface{ m() }

	switch x.(type) {
	case t:
	case *t:
	case *u:
	}

	switch x.(type) {
	case t:
	case *t:
	case *u:
	default:
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	}

	switch x.(type) {
	default:
	}

	switch x.(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
	}

	switch x.(type) {
	case t:
	default:
	}
}
