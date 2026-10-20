package foo

import "test/analysis/packagedoc/a"
import "test/analysis/packagedoc/b"

type T int

const (
	X0 T = iota
	X1
	X2
)

func f1(v T) {
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

func f2(v a.A) {
	switch v {
	case a.X0:
	case a.X1:
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
