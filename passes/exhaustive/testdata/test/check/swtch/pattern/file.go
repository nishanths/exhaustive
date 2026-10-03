package pattern

import "test/check/swtch/typ"
import "test/check/swtch/typnew"

type t1 int

const (
	w0 t1 = iota
	w1
)

type P1 int

const (
	x0 P1 = iota
	x1
	x2
	x3
)

func f1() {
	type P2 int
	const (
		Y0 P2 = iota
		Y1
	)

	type P3 int
	const (
		Z0 P3 = iota
		Z1
	)

	switch t1(0) {
	case w1:
	}

	switch P1(0) { // want "^switch not exhaustive: missing cases: x0, x3$"
	case x1:
	}

	switch P2(0) { // want "^switch not exhaustive: missing cases: Y0$"
	case Y1:
	}

	switch P3(0) { // want "^switch not exhaustive: missing cases: Z0$"
	case Z1:
	}

	switch typ.M1(0) {
	case typ.V1:
	}

	switch typnew.S1(0) { // want "^switch not exhaustive: missing cases: typnew.Z0, typnew.Z2$"
	case typnew.Z1:
	case typnew.Z4:
	}
}

func g1() {
	_ = map[t1]bool{
		w1: true,
	}

	_ = map[P1]bool{ // want "^map literal not exhaustive: missing keys: x0, x3$"
		x1: true,
	}
}
