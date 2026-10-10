package embedded

type t struct{}

func (t) m()
func (*t) n()
func (*t) o()

type u struct{}

func (*u) m()
func (u) n()

type w struct{}

func (*w) n()

type none struct{}

type et1 struct{ t }
type et2 struct{ *t }
type eu1 struct{ u }
type eu2 struct{ *u }
type ew1 struct{ w }
type ew2 struct{ *w }

func (*ew1) o() {}

type eTW struct {
	t
	ew1
}

func f1() {
	// embedded interface

	type i interface {
		m()
	}
	type e interface {
		i
		n()
	}
	var y e

	switch y.(type) { // want "^missing cases in type switch: \\*eTW, \\*et1, et2, \\*et2, \\*eu1, eu2, \\*eu2, \\*t, \\*u$"
	}

	switch y.(type) {
	case *eTW:
	case *et1:
	case et2:
	case *et2:
	case *eu1:
	case eu2:
	case *eu2:
	case *t:
	case *u:
	}
}

func f2() {
	// type contains embedded field

	var x interface{ o() }

	switch x.(type) { // want "^missing cases in type switch: \\*et1, et2, \\*et2, \\*ew1, \\*t$"
	}

	switch x.(type) {
	case *et1:
	case et2:
	case *et2:
	case *ew1:
	case *t:
	}

	var y interface{
		m()
		n()
		o()
	}

	switch y.(type) { // want "^missing cases in type switch: \\*et1, et2, \\*et2, \\*t$"
	}

	switch y.(type) {
	case *et1:
	case et2:
	case *et2:
	case *t:
	}
}
