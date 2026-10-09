package pattern

type t struct{}
func (t)  m()

type u struct{}
func (*u) m()

type i interface{ m() }
type h interface{ m() }

type a1 = h
type a2 = h

func f1() {
	type j interface{ m() }

	var x interface{ m() }
	var y i
	var z j

	switch x.(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
	}

	switch y.(type) {
	case t:
	}

	switch z.(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
	}

	switch ((h)(nil)).(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
	}

	switch ((a1)(nil)).(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
	}

	switch ((a2)(nil)).(type) {
	case t:
	}
}
