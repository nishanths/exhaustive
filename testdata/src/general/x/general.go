package x

import (
	bar "general/y"
	barpkg "general/y"
)

type Direction int

var (
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
	switch d { // want "^missing cases in switch of type Direction: E, directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}
}

func _b() {
	// Basic test of external package enum.
	//
	// Additionally: unexported members should not be included in exhaustiveness
	// check since enum is in external package.

	var p bar.Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: Mollusca$"
	case bar.Chordata:
	case bar.Echinodermata:
	}
}

func _j() {
	// Named imports still report real package name.

	var p barpkg.Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: Mollusca$"
	case barpkg.Chordata:
	case barpkg.Echinodermata:
	}
}

func _k(d Direction) {
	// Parenthesized values in case statements.

	switch d { // want "^missing cases in switch of type Direction: S, directionInvalid$"
	case (N):
	case (E):
	case (W):
	}
}

func _f() {
	// Multiple values in single case.

	var d Direction
	switch d { // want "^missing cases in switch of type Direction: W$"
	case E, directionInvalid, S:
	default:
	case N:
	}
}

func _g() {
	// Switch isn't at top-level of function -- should still be checked.

	var d Direction
	if true {
		switch d { // want "^missing cases in switch of type Direction: S, directionInvalid$"
		case (N):
		case (E):
		case (W):
		}
	}

	switch d { // want "^missing cases in switch of type Direction: E, S, W, directionInvalid$"
	case N:
		switch d { // want "^missing cases in switch of type Direction: N, S, W$"
		case E, directionInvalid:
		}
	}
}

type SortDirection int

const (
	_ SortDirection = iota // blank identifier need to be listed in switch statement to satisfy exhaustiveness
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

	switch p { // want "^missing cases in switch of type bar.Phylum: Mollusca$"
	case bar.Chordata:
	case bar.Echinodermata:
	case h.Mollusca:
	}
}
