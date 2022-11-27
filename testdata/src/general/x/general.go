package x

import (
	"errors"
	bar "general/y"
	barpkg "general/y"
)

type Direction int // want Direction:"^N,E,S,W,directionInvalid$"

const (
	N                Direction = 1
	E                Direction = 2
	S                Direction = 3
	W                Direction = 4
	directionInvalid Direction = 5
)

func _a() {
	// Basic test of same package enum.
	//
	// Additionally: unexported members should be included in exhaustiveness
	// check since enum is in same package.

	var d Direction
	switch d { // want "^missing cases in switch of type x.Direction: x.E, x.directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}

	_ = map[Direction]int{ // want "^missing keys in map of key type x.Direction: x.E, x.directionInvalid$"
		N: 1,
		S: 2,
		W: 3,
	}
}

func _b() {
	// Basic test of external package enum.
	//
	// Additionally: unexported members should not be included in exhaustiveness
	// check since enum is in external package.

	var p bar.Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: bar.Mollusca$"
	case bar.Chordata:
	case bar.Echinodermata:
	}

	_ = map[bar.Phylum]int{ // want "^missing keys in map of key type bar.Phylum: bar.Mollusca$"
		bar.Chordata:      1,
		bar.Echinodermata: 2,
	}
}

func _j() {
	// Named imports still report real package name.

	var p barpkg.Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: bar.Mollusca$"
	case barpkg.Chordata:
	case barpkg.Echinodermata:
	}

	_ = map[barpkg.Phylum]int{ // want "^missing keys in map of key type bar.Phylum: bar.Mollusca$"
		barpkg.Chordata:      1,
		barpkg.Echinodermata: 2,
	}
}

func _f() {
	// Multiple values in single case.

	var d Direction
	switch d { // want "^missing cases in switch of type x.Direction: x.W$"
	case E, directionInvalid, S:
	default:
	case N:
	}
}

func _g() {
	// Switch isn't at top-level of function -- should still be checked.

	var d Direction
	if true {
		switch d { // want "^missing cases in switch of type x.Direction: x.S, x.directionInvalid$"
		case (N):
		case (E):
		case (W):
		}
	}

	switch d { // want "^missing cases in switch of type x.Direction: x.E, x.S, x.W, x.directionInvalid$"
	case N:
		switch d { // want "^missing cases in switch of type x.Direction: x.N, x.S, x.W$"
		case E, directionInvalid:
		}
	}
}

type SortDirection int // want SortDirection:"^Asc,Desc$"

const (
	_ SortDirection = iota // blank identifier need not be listed in switch statement to satisfy exhaustiveness
	Asc
	Desc
)

func _n() {
	var d SortDirection
	switch d {
	case Asc:
	case Desc:
	}
}

func _o() {
	// Selector isn't of the form "enumPkg.enumMember"

	type holdsPhylum struct {
		Mollusca bar.Phylum // can technically hold any Phylum value, but field is named Mollusca
	}

	var p bar.Phylum
	var h holdsPhylum

	switch p { // want "^missing cases in switch of type bar.Phylum: bar.Mollusca$"
	case bar.Chordata:
	case bar.Echinodermata:
	case h.Mollusca:
	}

	_ = map[bar.Phylum]int{ // want "^missing keys in map of key type bar.Phylum: bar.Mollusca$"
		bar.Chordata:      1,
		bar.Echinodermata: 2,
		h.Mollusca:        3,
	}
}

var ErrFoo = errors.New("foo")

func _p() {
	// Switch tag variable's type has nil package (lives in Universe scope).
	// Expect to not panic and to not fail unexpectedly.

	var err error

	switch err {
	case nil:
	case ErrFoo:
	}

	_ = map[error]int{
		nil:    1,
		ErrFoo: 2,
	}
}

func _r(d Direction) {
	// Raw constants (i.e. not identifier or sel expr)
	// in case clauses do not count.

	switch d { // want "^missing cases in switch of type x.Direction: x.S, x.directionInvalid$"
	case N:
	case E:
	case 3:
	case W:
	case 5:
	}

	_ = map[Direction]int{ // want "^missing keys in map of key type x.Direction: x.S, x.directionInvalid$"
		N: 1,
		E: 2,
		3: 3,
		W: 4,
		5: 5,
	}
}

func _s(u bar.Uppercase) {
	switch u {
	case bar.ReallyExported:
	}

	_ = map[bar.Uppercase]int{
		bar.ReallyExported: 1,
	}
}

func _ufunc0() Direction            { return directionInvalid }
func _ufunc1(d Direction) Direction { return d }

func _u() {
	switch _ufunc0() { // want "^missing cases in switch of type x.Direction: x.W, x.directionInvalid$"
	case N, E, S:
	}

	switch _ufunc1(Direction(0)) { // want "^missing cases in switch of type x.Direction: x.W, x.directionInvalid$"
	case N, E, S:
	}
}

func _w() {
	// assignment in switch
	switch x := _ufunc0(); x { // want "^missing cases in switch of type x.Direction: x.E, x.S, x.W$"
	case N, directionInvalid:
	}
}

func mapTypeAlias() {
	type myMapAlias = map[Direction]int

	_ = myMapAlias{ // want "^missing keys in map of key type x.Direction: x.S, x.directionInvalid$"
		N: 1,
		E: 2,
		W: 4,
	}
}

func customMapType() {
	type myMap map[Direction]int

	_ = myMap{ // want "^missing keys in map of key type x.Direction: x.S, x.directionInvalid$"
		N: 1,
		E: 2,
		W: 4,
	}
}

func mapKeyFuncCall() {
	_ = map[Direction]int{ // want "^missing keys in map of key type x.Direction: x.N, x.E, x.W, x.directionInvalid$"
		_ufunc0():               1,
		_ufunc1(Direction(100)): 1,
		S:                       1,
	}
}
