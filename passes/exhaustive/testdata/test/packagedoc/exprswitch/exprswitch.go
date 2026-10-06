package foo

import "test/packagedoc/a"
import "test/packagedoc/b"

func f1() {
	type T int

	const (
		X0 T = iota
		X1
		X2
	)

	var t T

	switch t {
	case X0:
	case X1:
	case X2:
	}

	switch t { // want "^expression switch not exhaustive: missing cases: X1$"
	case X0:
	case X2:
	}
}

func f2() {
	var t a.T
	switch t {
	case a.X0: // or b.X0
	case a.X1: // or b.X1
	}

	switch t {
	case b.X0:
	case b.X1:
	}

	switch t { // want "^expression switch not exhaustive: missing cases: b.X1$"
	case a.X0:
	}

	switch t { // want "^expression switch not exhaustive: missing cases: b.X1$"
	case b.X0:
	}
}
