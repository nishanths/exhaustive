package x

import (
	bar "github.com/nishanths/exhaustive/testdata/y"
)

type Dir int

var (
	N          Dir = 1
	E          Dir = 2
	S          Dir = 3
	W          Dir = 4
	invalidDir Dir = 5
)

func _foo() {
	var d Dir
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

// Is this doc comment?
var (
	something map[string]int

	// Is this doc comment?
	//
	//exhaustive:ignore
	p, q = map[bar.Phylum]string{
		bar.Chordata:      "c",
		bar.Echinodermata: "e",
	}, map[Dir]string{
		N: "n",
	} // line
)
