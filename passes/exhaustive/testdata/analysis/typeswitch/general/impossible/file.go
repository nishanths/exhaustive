package impossible

type m interface{ m() }
type n interface{ n() }
type mn interface {
	m()
	n()
}
type mn2 interface { // same as mn
	m()
	n()
}
type e interface{ e() }

type t struct{}

func (t) m()
func (t) n()

type u struct{}

func (*u) m()

func f0() {
	var x m

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	}

	switch x.(type) {
	case t:
	case *t:
	case *u:
	}

	switch x.(type) {
	case m:
	}

	switch x.(type) {
	case m:
	case *t: // want "^ineffectual: \\*t -> m$"
	}

	switch x.(type) {
	case interface {
		m()
	}:
	case *t: // want "^ineffectual: \\*t -> interface\\{m\\(\\)\\}$"
	}

	switch x.(type) {
	case m:
	case t: // want "^ineffectual: t -> m$"
	case *t: // want "^ineffectual: \\*t -> m$"
	case *u: // want "^ineffectual: \\*u -> m$"
	}

	switch x.(type) {
	case t:
	case *t:
	case *u:
	case m:
	}

	switch x.(type) {
	case *u:
	case m, *t:
	}

	switch x.(type) {
	case *u:
	case *t, m:
	}

	switch x.(type) {
	case n:
	case t: // want "^ineffectual: t -> n$"
	case *t: // want "^ineffectual: \\*t -> n$"
	case *u:
	}

	switch x.(type) { // want "^missing cases in type switch: \\*u$"
	case mn:
	case *t: // want "^ineffectual: \\*t -> mn$"
	}

	switch x.(type) {
	case m:
	case mn: // want "^ineffectual: mn -> m$"
	}

	switch x.(type) {
	case mn:
	case m:
	}

	switch x.(type) {
	case t:
	case n, m:
	case *t, *u: // want "^ineffectual: \\*t -> n$" "^ineffectual: \\*t -> m$" "^ineffectual: \\*u -> m$"
	}

	switch x.(type) {
	case m, mn, n, t, *t, *u:
	}

	switch x.(type) {
	case e:
	case n:
	case *u:
	}

	switch x.(type) {
	case any:
	case *t: // want "^ineffectual: \\*t -> any$"
	case mn: // want "^ineffectual: mn -> any$"
	}

	switch x.(type) {
	case *t:
	case m:
	case n:
	case any:
	}

	// Note: e.g. "*t -> mn" is true
	// but should not be reported.
	// See comment in func checkTypeSwitchIneff.
	switch x.(type) {
	case m:
	case mn: // want "^ineffectual: mn -> m$"
	case n:
	case mn2: // want "^ineffectual: mn2 -> m$" "^ineffectual: mn2 -> n$"
	case *t: // want "^ineffectual: \\*t -> m$" "^ineffectual: \\*t -> n$"
	}
}

func f1() {
	var x interface{}

	switch x.(type) {
	}

	switch x.(type) {
	case t:
	case *t:
	case *u:
	}

	switch x.(type) {
	case m:
	}

	switch x.(type) {
	case m:
	case *t: // want "^ineffectual: \\*t -> m$"
	}
}

type o interface{ o() }

type t1 struct{}

func (t1) o()

type t2 struct{ f []byte } // note that type t2 is not comparable

func (t2) o()

func f2[P any, Q o, R1 comparable, R2 interface {
	comparable
	o
}, R3 t1, R4 t2, R5 t1 | float64, R6 *t1, R7 interface {
	comparable
	t1
}, R8 interface {
	comparable
	t2
}, R9 interface {
	o
	m
}, R10 interface {
	t1
	o
}]() {
	// These scenarios are not supported.
	// See the comment in func checkTypeSwitchIneff.

	var v o

	switch v.(type) {
	case P:
	case Q:
	case R1:
	case R2:
	case R3:
	case R4:
	case R5:
	case R6:
	case R7:
	case R8:
	case o:
	}

	switch v.(type) {
	case P:
	case Q:
	case R1:
	case R2:
	case R3:
	case R4:
	case R5:
	case R6:
	case R7:
	case R8:
	case t1:
	case t2:
	case *t1:
	case *t2:
	}

	switch v.(type) {
	case o:
	case P:
	case Q: // want "^ineffectual: Q -> o$"
	case R1:
	case R2: // want "^ineffectual: R2 -> o$"
	case R3:
	case R7:
	case R9: // want "^ineffectual: R9 -> o$"
	case R10: // want "^ineffectual: R10 -> o$"
	}
}
