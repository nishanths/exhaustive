package quux

import (
	"typealias/bar"
	"typealias/foo"
)

func x() {
	var v foo.T1 = foo.ReturnsT1()
	switch v { // want "^missing cases in switch of type bar.T2: bar.D, bar.E$"
	case foo.A:
	case bar.B:
	case foo.C:
	case foo.D:
	case foo.F:
	case foo.H:
	}

	var w bar.T2 = foo.ReturnsT1()
	switch w { // want "^missing cases in switch of type bar.T2: bar.D, bar.E$"
	case foo.A:
	case bar.B:
	case foo.C:
	case foo.D:
	case foo.F:
	case foo.H:
	}

	_ = map[foo.T1]int{ // want "^missing keys in map of key type bar.T2: bar.D, bar.E$"
		foo.A: 1,
		bar.B: 2,
		foo.C: 3,
		foo.D: 4,
		foo.F: 5,
		foo.H: 6,
	}

	_ = map[bar.T2]int{ // want "^missing keys in map of key type bar.T2: bar.D, bar.E$"
		foo.A: 1,
		bar.B: 2,
		foo.C: 3,
		foo.D: 4,
		foo.F: 5,
		foo.H: 6,
	}
}
