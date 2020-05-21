package x

import bar "switch/y"

// These switches are exhaustive, expect no diagnostics.

func _l() {
	var d Direction
	switch d {
	case N, E, S, W, directionInvalid:
	}
}

func _m() {
	var p bar.Phylum
	switch p {
	case bar.Echinodermata, bar.Mollusca, bar.Chordata:
	}
}
