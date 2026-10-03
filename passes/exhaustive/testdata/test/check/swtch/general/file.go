package general

import "test/check/swtch/typ"
import "test/check/swtch/typnew"

type T1 int

const (
	X0 T1 = iota
	X1
	X2
	X3
	X4
)

// switch statement is exhaustive
func f1() {
	var t T1

	switch t {
	case X0:
	case X1:
	case X2:
	case X3:
	case X4:
	}

	switch t {
	case X3:
	case X0:
	case X1:
	case X4:
	case X2:
	}

	switch t {
	case X0, X1, X4:
	case X3, X2:
	}

	_ = map[T1]bool{
		X3: true,
		X0: true,
		X1: true,
		X4: true,
		X2: true,
	}
	_ = map[T1]bool{X3: true, X0: true, X1: true, X4: true, X2: true}
}

// simple non-exhaustive switch statement
func f2() {
	var t T1

	switch t { // want "^switch not exhaustive: missing cases: X1, X3, X4$"
	case X0:
	case X2:
	}

	switch t { // want "^switch not exhaustive: missing cases: X0, X1, X2, X3, X4$"
	}

	_ = map[T1]bool{ // want "^map literal not exhaustive: missing keys: X1, X3, X4$"
		X0: true,
		X2: true,
	}
	_ = map[T1]bool{X0: true, X2: true} // want "^map literal not exhaustive: missing keys: X1, X3, X4$"
	_ = map[T1]bool{}                   // want "^map literal not exhaustive: missing keys: X0, X1, X2, X3, X4$"
	_ = make(map[T1]bool)               // Note: use 'make' instead of writing a map literal with zero elements.
}

// multiple names for same value
func f3() {
	type m1 int

	const (
		v0, vv0 m1 = iota, iota
		v1      m1 = iota
		v2, vv2 m1 = iota, iota
		v3, vv3
		v4 m1 = iota
		v5
	)

	var t m1
	switch t { // want "^switch not exhaustive: missing cases: v0\\|vv0, v1, v2\\|vv2, v4$"
	case vv3:
	case v5:
	}

	_ = map[m1]bool{vv3: true, v5: true} // want "map literal not exhaustive: missing keys: v0\\|vv0, v1, v2\\|vv2, v4$"
}

// exported and unexported names.
// also: multiple names for same value.
func f4() {
	var t typ.M1
	switch t { // want "^switch not exhaustive: missing cases: typ.V0\\|typ.VV0, typ.V1, typ.V2, typ.V4$"
	case typ.V3:
	case typ.V5:
	}

	_ = map[typ.M1]bool{typ.V3: true, typ.V5: true} // want "map literal not exhaustive: missing keys: typ.V0\\|typ.VV0, typ.V1, typ.V2, typ.V4$"
}

// nested switch statements
func f6() {
	var outer, inner T1

	switch outer { // want "^switch not exhaustive: missing cases: X2, X3$"
	case X0:
	case X1:
		switch inner { // want "^switch not exhaustive: missing cases: X0, X2$"
		case X1:
		case X3:
		case X4:
		}
	case X4:
	}

	_ = map[T1]map[T1]bool{ // want "^map literal not exhaustive: missing keys: X2, X3$"
		X0: {}, // want "^map literal not exhaustive: missing keys: X0, X1, X2, X3, X4$"
		X1: { // want "^map literal not exhaustive: missing keys: X0, X2$"
			X1: true,
			X3: true,
			X4: true,
		},
		X4: {
			X0: true,
			X1: true,
			X3: true,
			X2: true,
			X4: true,
		},
	}

	_ = map[T1]*map[T1]bool{ // want "^map literal not exhaustive: missing keys: X2, X3$"
		X0: {}, // want "^map literal not exhaustive: missing keys: X0, X1, X2, X3, X4$"
		X1: { // want "^map literal not exhaustive: missing keys: X0, X2$"
			X1: true,
			X3: true,
			X4: true,
		},
		X4: {
			X0: true,
			X1: true,
			X3: true,
			X2: true,
			X4: true,
		},
	}
}

// switch statement with default case
func f7() {
	var t T1

	switch t {
	case X0:
	case X1:
	case X2:
	case X3:
	case X4:
	default:
	}

	switch t { // want "^switch not exhaustive: missing cases: X0, X3$"
	case X1:
	case X2:
	case X4:
	default:
	}

	switch t { // want "^switch not exhaustive: missing cases: X0, X1, X2, X3, X4$"
	default:
	}
}

// type conversions; and literal values.
func f8() {
	var t T1

	switch T1((int(T1((t))))) {
	case T1(0):
	case T1((int(T1((1))))):
	case X2:
	case X4:
	case T1(3):
	}

	switch T1((int(T1((t))))) { // want "^switch not exhaustive: missing cases: X2$"
	case T1(0):
	case T1((int(T1((1))))):
	case X4:
	case T1(3):
	}

	switch t {
	case X0:
	case T1(1000):
	case X1:
	case X2:
	case X3:
	case X4:
	}

	_ = map[T1]bool{
		T1(0):              true,
		T1((int(T1((1))))): true,
		X2:                 true,
		X4:                 true,
		T1(3):              true,
	}

	_ = map[T1]bool{ // want "^map literal not exhaustive: missing keys: X2$"
		T1(0):              true,
		T1((int(T1((1))))): true,
		X4:                 true,
		T1(3):              true,
	}

	_ = map[T1]bool{
		X0:       true,
		T1(1000): true,
		X1:       true,
		X2:       true,
		X3:       true,
		X4:       true,
	}
}

// non-constant case expression
func f9(t T1) {
	var x T1 = X1

	switch t { // want "^switch not exhaustive: missing cases: X1, X2, X4$"
	case X0:
	case x:
	case x + 1:
	case X3:
	}

	_ = map[T1]bool{ // want "^map literal not exhaustive: missing keys: X1, X2, X4$"
		X0:    true,
		x:     true,
		x + 1: true,
		X3:    true,
	}
}

// inner composite literal type omitted in the source.
func f10() {
	_ = []map[T1]bool{
		{X1: true},           // want "^map literal not exhaustive: missing keys: X0, X2, X3, X4$"
		{X0: true, X3: true}, // want "^map literal not exhaustive: missing keys: X1, X2, X4$"
	}
	_ = []*map[T1]bool{
		{X1: true},           // want "^map literal not exhaustive: missing keys: X0, X2, X3, X4$"
		{X0: true, X3: true}, // want "^map literal not exhaustive: missing keys: X1, X2, X4$"
	}
}

// composite literal with underlying type map
func f11() {
	type n map[T1]bool
	_ = map[T1]bool{X1: true}    // want "^map literal not exhaustive: missing keys: X0, X2, X3, X4$"
	_ = n{X1: true}              // want "^map literal not exhaustive: missing keys: X0, X2, X3, X4$"
	_ = map[T1]bool(n{X1: true}) // want "^map literal not exhaustive: missing keys: X0, X2, X3, X4$"
	_ = n(map[T1]bool{X1: true}) // want "^map literal not exhaustive: missing keys: X0, X2, X3, X4$"
}

// pointer to map type; address operators
func f12() {
	type n map[T1]bool
	_ = &map[T1]bool{X0: true, X3: true} // want "^map literal not exhaustive: missing keys: X1, X2, X4$"
	_ = &n{X0: true, X3: true}           // want "^map literal not exhaustive: missing keys: X1, X2, X4$"
}

// type of the switch expression is not an enumerated type
func f13() {
	var v int
	switch v {
	}

	switch 1 {
	}

	_ = map[int]bool{
		0: true,
		1: true,
		2: true,
	}
}

// no switch expression
func f14() {
	switch {
	}
}

// type aliases
func f15() {
	var a typ.A1

	switch a { // want "^switch not exhaustive: missing cases: typnew.Z0, typnew.Z2, typnew.Z4$"
	case typ.Za1:
	case typ.Za3:
	}

	switch a { // want "^switch not exhaustive: missing cases: typnew.Z0, typnew.Z2, typnew.Z4$"
	case typnew.Z1:
	case typnew.Z3:
	}

	switch a { // want "^switch not exhaustive: missing cases: typnew.Z0, typnew.Z4$"
	case typnew.Z1:
	case typnew.Z2:
	case typ.Za3:
	}

	var b typnew.S1

	switch b { // want "^switch not exhaustive: missing cases: typnew.Z0, typnew.Z4$"
	case typnew.Z1:
	case typnew.Z2:
	case typ.Za3:
	}

	_ = map[typ.A1]bool{ // want "^map literal not exhaustive: missing keys: typnew.Z0, typnew.Z2, typnew.Z4$"
		typ.Za1: true,
		typ.Za3: true,
	}

	_ = map[typ.A1]bool{ // want "^map literal not exhaustive: missing keys: typnew.Z0, typnew.Z2, typnew.Z4$"
		typnew.Z1: true,
		typnew.Z3: true,
	}

	_ = map[typ.A1]bool{ // want "^map literal not exhaustive: missing keys: typnew.Z0, typnew.Z4$"
		typnew.Z1: true,
		typnew.Z2: true,
		typ.Za3:   true,
	}

	_ = map[typnew.S1]bool{ // want "^map literal not exhaustive: missing keys: typnew.Z0, typnew.Z4$"
		typnew.Z1: true,
		typnew.Z2: true,
		typ.Za3:   true,
	}
}

// non-identifier expression in switch expression.
func f16() {
	var t T1
	f := func() T1 { return X0 }

	// Note: This is a weird test case, but it captures
	// the current behavior.
	switch t + 9999 { // want "^switch not exhaustive: missing cases: X1, X3, X4$"
	case X0:
	case X2:
	}

	switch f() { // want "^switch not exhaustive: missing cases: X1, X3, X4$"
	case X0:
	case X2:
	}
}
