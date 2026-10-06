package alias

import "test/check/typeswitch/alias/sub"

type t struct{}

func (t) m()

type u struct{}

func (*u) m()

type a = t
type aa = *t
type aa2 = aa
type c = u
type cc = *u

func f4() {
	var x interface{ m() }

	switch x.(type) {
	case a:
	case aa:
	case cc:
	}

	switch x.(type) {
	case a:
	case *t:
	case cc:
	}

	switch x.(type) { // want "^type switch not exhaustive: missing cases: t or a, \\*u or \\*c or cc$"
	case *t:
	}

	switch x.(type) { // want "^type switch not exhaustive: missing cases: t or a, \\*u or \\*c or cc$"
	case *a:
	}

	switch x.(type) { // want "^type switch not exhaustive: missing cases: \\*t or \\*a or aa or aa2$"
	case a:
	case *u:
	}

	switch x.(type) { // want "^type switch not exhaustive: missing cases: t or a, \\*t or \\*a or aa or aa2, \\*u or \\*c or cc$"
	}
}

func f5() {
	var x interface{ N() }

	switch x.(type) { // want "^type switch not exhaustive: missing cases: sub.V, \\*sub.V or sub.VV$"
	}

	switch x.(type) { // want "^type switch not exhaustive: missing cases: \\*sub.V or sub.VV$"
	case sub.V:
	}

	switch x.(type) { // want "^type switch not exhaustive: missing cases: sub.V$"
	case *sub.V:
	}

	switch x.(type) { // want "^type switch not exhaustive: missing cases: sub.V$"
	case sub.VV:
	}

	switch x.(type) {
	case sub.V:
	case *sub.V:
	}

	switch x.(type) {
	case sub.V:
	case sub.VV:
	}
}
