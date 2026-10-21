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
	case *t: // want "^unreachable type switch case \\*t \\(matched by earlier case m\\)$"
	}

	switch x.(type) {
	case interface {
		m()
	}:
	case *t: // want "^unreachable type switch case \\*t \\(matched by earlier case interface\\{m\\(\\)\\}\\)$"
	}

	switch x.(type) {
	case m:
	case t: // want "^unreachable type switch case t \\(matched by earlier case m\\)$"
	case *t: // want "^unreachable type switch case \\*t \\(matched by earlier case m\\)$"
	case *u: // want "^unreachable type switch case \\*u \\(matched by earlier case m\\)$"
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
	case t: // want "^unreachable type switch case t \\(matched by earlier case n\\)$"
	case *t: // want "^unreachable type switch case \\*t \\(matched by earlier case n\\)$"
	case *u:
	}

	switch x.(type) { // want "^missing cases in type switch: \\*u$"
	case mn:
	case *t: // want "^unreachable type switch case \\*t \\(matched by earlier case mn\\)$"
	}

	switch x.(type) {
	case m:
	case mn: // want "^unreachable type switch case mn \\(matched by earlier case m\\)$"
	}

	switch x.(type) {
	case mn:
	case m:
	}

	switch x.(type) {
	case t:
	case n, m:
	case *t, *u: // want "^unreachable type switch case \\*t \\(matched by earlier case n\\)$" "^unreachable type switch case \\*t \\(matched by earlier case m\\)$" "^unreachable type switch case \\*u \\(matched by earlier case m\\)$"
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
	case *t: // want "^unreachable type switch case \\*t \\(matched by earlier case any\\)$"
	}
}

// See Go spec section "Satisfying a type constraint".

type o interface{ o() }

type t1 struct{}
func (t1) o()

type t2 struct{ f []byte }
func (t2) o()

func f2[P any](v o) {
	switch v.(type) {
	case P:
	case t1: // want "^unreachable type switch case t1 \\(matched by earlier case P\\)$"
	case t2: // want "^unreachable type switch case t2 \\(matched by earlier case P\\)$"
	case o: // want "^unreachable type switch case o \\(matched by earlier case P\\)$"
	case n: // want "^unreachable type switch case n \\(matched by earlier case P\\)$"
	}
}

func f3[P o](v o) {
	switch v.(type) {
	case P:
	case t1: // want "^unreachable type switch case t1 \\(matched by earlier case P\\)$"
	case t2: // want "^unreachable type switch case t2 \\(matched by earlier case P\\)$"
	case o: // want "^unreachable type switch case o \\(matched by earlier case P\\)$"
	case n:
	}
}

func f4[P comparable](v o) {
	switch v.(type) {
	case P:
	case t1: // want "^unreachable type switch case t1 \\(matched by earlier case P\\)$"
	case t2:
	case o: // want "^unreachable type switch case o \\(matched by earlier case P\\)$"
	case n: // want "^unreachable type switch case n \\(matched by earlier case P\\)$"
	}

	switch v.(type) {
	case P:
	case o: // want "^unreachable type switch case o \\(matched by earlier case P\\)$"
	case n: // want "^unreachable type switch case n \\(matched by earlier case P\\)$"
	case t1: // want "^unreachable type switch case t1 \\(matched by earlier case P\\)$" "^unreachable type switch case t1 \\(matched by earlier case o\\)$"
	case t2: // want "^unreachable type switch case t2 \\(matched by earlier case o\\)$"
	}
}

func f5[P interface {
	comparable
	o()
}](v o) {
	switch v.(type) {
	case P:
	case t1: // want "^unreachable type switch case t1 \\(matched by earlier case P\\)$"
	case t2:
	case o: // want "^unreachable type switch case o \\(matched by earlier case P\\)$"
	case n:
	}

	switch v.(type) {
	case P:
	case o: // want "^unreachable type switch case o \\(matched by earlier case P\\)$"
	case n:
	case t1: // want "^unreachable type switch case t1 \\(matched by earlier case P\\)$" "^unreachable type switch case t1 \\(matched by earlier case o\\)$"
	case t2: // want "^unreachable type switch case t2 \\(matched by earlier case o\\)$"
	}
}
