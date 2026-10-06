package pattern

type t struct{}
func (t)  m()

type u struct{}
func (*u) m()

type i interface{ m() }

func f1() {
	type j interface{ m() }

	var x interface{ m() }
	var y i
	var z j

	switch x.(type) { // want "^type switch not exhaustive: missing cases: \\*t, \\*u$"
	case t:
	}

	switch y.(type) {
	case t:
	}

	switch z.(type) { // want "^type switch not exhaustive: missing cases: \\*t, \\*u$"
	case t:
	}
}
