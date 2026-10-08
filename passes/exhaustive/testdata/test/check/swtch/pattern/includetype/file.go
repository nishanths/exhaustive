package includetype

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

type p4 int

const (
	u0 p4 = iota
	u1
)

type a1 = p4
type a2 = p4

func f1() {
	type P3 int
	const (
		Z0 P3 = iota
		Z1
	)

	switch t1(0) { // want "^expression switch not exhaustive: missing cases: w0$"
	case w1:
	}

	switch P1(0) {
	case x1:
	}

	switch P3(0) {
	case Z1:
	}

	switch typ.M1(0) { // want "^expression switch not exhaustive: missing cases: typ.VV0, typ.V2, typ.V3, typ.V5$"
	case typ.V1:
	}

	switch typnew.S1(0) {
	case typnew.Z1:
	case typnew.Z4:
	}

	switch p4(0) {
	case u1:
	}

	switch a1(0) {
	case u1:
	}

	switch a2(0) { // want "^expression switch not exhaustive: missing cases: u0$"
	case u1:
	}
}

func g1() {
	_ = map[t1]bool{ // want "^map literal not exhaustive: missing keys: w0$"
		w1: true,
	}

	_ = map[P1]bool{
		x1: true,
	}

	_ = map[p4]bool{
		u1: true,
	}

	_ = map[a1]bool{
		u1: true,
	}

	_ = map[a2]bool{ // want "^map literal not exhaustive: missing keys: u0$"
		u1: true,
	}
}
