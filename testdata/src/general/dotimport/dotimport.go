package x

import (
	. "general/y"
	bar "general/y"
	barpkg "general/y"
)

func _dot() {
	var p Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: Chordata, Mollusca$"
	case Echinodermata:
	}

	_ = map[Phylum]int{ // want "^missing keys in map of key type bar.Phylum: Chordata, Mollusca$"
		Echinodermata: 1,
	}
}

func _mixed() {
	var p bar.Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: Mollusca$"
	case Echinodermata:
	case barpkg.Chordata:
	}

	_ = map[bar.Phylum]int{ // want "^missing keys in map of key type bar.Phylum: Mollusca$"
		Echinodermata:   1,
		barpkg.Chordata: 2,
	}
}
