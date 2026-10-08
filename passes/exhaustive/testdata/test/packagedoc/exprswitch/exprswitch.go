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

	var v T

	switch v {
	case X0:
	case X1:
	case X2:
	}

	switch v { // want "^missing cases in expression switch: X2$"
	case X0:
	case X1:
	}
}

func f2() {
	var v a.A

	switch v {
	case a.X0: // or b.X0
	case a.X1: // or b.X1
	}

	switch v {
	case b.X0:
	case b.X1:
	}

	switch v { // want "^missing cases in expression switch: b.X1$"
	case a.X0:
	}

	switch v { // want "^missing cases in expression switch: b.X1$"
	case b.X0:
	}
}
