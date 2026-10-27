package impossible

type m interface{ m() }
type n interface{ n() }
type mn interface {
	m()
	n()
}
type e interface{ e() }

type t struct{}

func (t) m()
func (t) n()

type u struct{}

func (*u) m()

func f1() {
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
	case *t: // want "^ineffectual type \\*t due to case order: dynamic type that would be matched by type \\*t is matched in earlier case by type m$"
	}

	switch x.(type) {
	case interface {
		m()
	}:
	case *t: // want "^ineffectual type \\*t due to case order: dynamic type that would be matched by type \\*t is matched in earlier case by type interface\\{m\\(\\)\\}$"
	}

	switch x.(type) {
	case m:
	case t: // want "^ineffectual type t due to case order: dynamic type that would be matched by type t is matched in earlier case by type m$"
	case *t: // want "^ineffectual type \\*t due to case order: dynamic type that would be matched by type \\*t is matched in earlier case by type m$"
	case *u: // want "^ineffectual type \\*u due to case order: dynamic type that would be matched by type \\*u is matched in earlier case by type m$"
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
	case t: // want "^ineffectual type t due to case order: dynamic type that would be matched by type t is matched in earlier case by type n$"
	case *t: // want "^ineffectual type \\*t due to case order: dynamic type that would be matched by type \\*t is matched in earlier case by type n$"
	case *u:
	}

	switch x.(type) { // want "^missing cases in type switch: \\*u$"
	case mn:
	case *t: // want "^ineffectual type \\*t due to case order: dynamic type that would be matched by type \\*t\\ is matched in earlier case by type mn$"
	}

	switch x.(type) {
	case m:
	case mn: // want "^ineffectual type mn due to case order: dynamic type that would be matched by type mn\\ is matched in earlier case by type m$"
	}

	switch x.(type) {
	case mn:
	case m:
	}

	switch x.(type) {
	case t:
	case n, m:
	case *t, *u: // want "^ineffectual type \\*t due to case order: dynamic type that would be matched by type \\*t is matched in earlier case by type n$" "^ineffectual type \\*t due to case order: dynamic type that would be matched by type \\*t is matched in earlier case by type m$" "^ineffectual type \\*u due to case order: dynamic type that would be matched by type \\*u is matched in earlier case by type m$"
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
	case *t: // want "^ineffectual type \\*t due to case order: dynamic type that would be matched by type \\*t is matched in earlier case by type any$"
	case mn: // want "^ineffectual type mn due to case order: dynamic type that would be matched by type mn is matched in earlier case by type any$"
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
}]() {
	// Note: checking for impossible/unreachable cases due to
	// an earlier type constraint interface in a case is not
	// supported by this analysis.
	// See the comment in func checkImpossibleCases.

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
}
