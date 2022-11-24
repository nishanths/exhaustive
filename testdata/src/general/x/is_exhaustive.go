package x

import bar "general/y"

// These are exhaustive, expect no diagnostics.

func _l() {
	var d Direction
	switch d {
	case N, E, S, W, directionInvalid:
	}

	_ = map[Direction]int{
		N:                1,
		E:                2,
		S:                3,
		W:                4,
		directionInvalid: 5,
	}
}

func _m() {
	var p bar.Phylum
	switch p {
	case bar.Echinodermata, bar.Mollusca, bar.Chordata:
	}

	_ = map[bar.Phylum]int{
		bar.Echinodermata: 1,
		bar.Mollusca:      2,
		bar.Chordata:      3,
	}
}
