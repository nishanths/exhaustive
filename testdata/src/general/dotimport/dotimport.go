package x

import (
	. "general/y"
	barpkg "general/y"
)

func _dot() {
	var p Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: Chordata, Mollusca$"
	case Echinodermata:
	}
}

func _mixed() {
	var p Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: Mollusca$"
	case Echinodermata:
	case barpkg.Chordata:
	}
}
