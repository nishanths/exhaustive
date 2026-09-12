package foo

import "packagedoc/oldpkg"
import "packagedoc/newpkg"

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

	switch t { // want "^switch not exhaustive: missing cases: X1, X2$"
	case X0:
	}
}

func f2() {
	var t oldpkg.T

	switch t {
	case oldpkg.X0: // equivalently newpkg.X0
	case oldpkg.X1: // equivalently newpkg.X1
	}

	switch t {
	case newpkg.X0:
	case newpkg.X1:
	}

	switch t { // want "^switch not exhaustive: missing cases: newpkg.X1$"
	case oldpkg.X0:
	}

	switch t { // want "^switch not exhaustive: missing cases: newpkg.X1$"
	case newpkg.X0:
	}
}
