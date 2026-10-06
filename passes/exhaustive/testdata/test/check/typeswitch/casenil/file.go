package casenil

type t struct{}
func (t)  m()

type u struct{}
func (*u) m()

func f1() {
	var x interface{ m() }

	switch x.(type) { // want "^type switch not exhaustive: missing cases: nil$"
	case t:
	case *t:
	case *u:
	}

	switch x.(type) { // want "^type switch not exhaustive: missing cases: nil, t, \\*t, \\*u$"
	}

	switch x.(type) { // want "^type switch not exhaustive: missing cases: nil, \\*t, \\*u$"
	case t:
	}
}

func f2() {
	var x interface{ m() }

	switch x.(type) {
	case t:
	case *t:
	case *u:
	case nil:
	}

	switch x.(type) { // want "^type switch not exhaustive: missing cases: t, \\*t, \\*u$"
	case nil:
	}

	switch x.(type) { // want "^type switch not exhaustive: missing cases: \\*t, \\*u$"
	case t:
	case nil:
	}
}
