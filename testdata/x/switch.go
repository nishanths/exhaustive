package x

import (
	bar "github.com/nishanths/exhaustive/testdata/y"
	barpkg "github.com/nishanths/exhaustive/testdata/y"
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
	// Basic same package.

	var d Direction
	switch d {
	case N:
	case S:
	case W:
	}
}

func _b() {
	// Basic external package.

	var p bar.Phylum
	switch p {
	case bar.Chordata:
	case bar.Echinodermata:
	}
}

func _j() {
	// Named import.

	var p barpkg.Phylum
	switch p {
	case barpkg.Chordata:
	case barpkg.Echinodermata:
	}
}

func _h() {
	var p bar.Phylum
	switch p {
	}
}

/*
func _f() {
	// Multiple values in single case.

	var d Direction
	switch d {
	case E, directionInvalid, S:
	default:
	case N:
	}
}

func _g() {
	// Non top-level switch.

	var d Direction
	if true {
		switch d {
		case (N):
		case (E):
		case (W):
		}
	}
}
*/
