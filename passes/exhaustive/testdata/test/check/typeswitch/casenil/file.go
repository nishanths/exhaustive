package casenil

type t struct{}
func (t)  m()

type u struct{}
func (*u) m()

func f1() {
	var x interface{ m() }

	switch x.(type) { // want "^missing cases in type switch: nil$"
	case t:
	case *t:
	case *u:
	}

	switch x.(type) { // want "^missing cases in type switch: nil, t, \\*t, \\*u$"
	}

	switch x.(type) { // want "^missing cases in type switch: nil, \\*t, \\*u$"
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

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	case nil:
	}

	switch x.(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
	case nil:
	}
}
