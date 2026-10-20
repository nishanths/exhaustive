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
	case interface{ m() }:
	case *t: // want "^unreachable type switch case \\*t \\(matched by earlier case interface\\{ m\\(\\) \\}\\)$"
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

func f2[P m]() {
	var x m

	switch x.(type) {
	case P:
	case *t: // want "^unreachable type switch case \\*t \\(matched by earlier case P\\)$"
	}
}
