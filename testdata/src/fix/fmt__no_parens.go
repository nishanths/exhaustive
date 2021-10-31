package fix

import "general/y"

func _j() {
	var p bar.Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: Mollusca$"
	case bar.Chordata:
	case bar.Echinodermata:
	}
}
