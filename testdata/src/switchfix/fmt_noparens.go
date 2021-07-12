package switchfix

import "switch/y"

func _fmt_noparens() {
	var p bar.Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: Mollusca$"
	case bar.Chordata:
	case bar.Echinodermata:
	}
}
