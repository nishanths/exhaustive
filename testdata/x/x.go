package x

import (
	bar "github.com/nishanths/exhaustive/testdata/y"
)

type Direction int

var (
	N                Direction = 1
	E                Direction = 2
	S                Direction = 3
	W                Direction = 4
	directionInvalid Direction = 5
)

func _foo() {
	var d Direction
	switch d {
	case N:
	}
}

func _bar() {
	var p bar.Phylum

	switch p {
	case bar.Chordata:
	case bar.Echinodermata:
	}
}

var (
	something map[string]int

	p, q = map[bar.Phylum]string{
		bar.Chordata: "c",
		bar.Mollusca: "m",
	}, map[Direction]string{
		N: "n",
		S: "s",
	}
)
