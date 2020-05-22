package switchfix

import (
	barpkg "switch/y"
)

func _caseclause_othermemberspkg() {
	var phy barpkg.Phylum
	switch phy { // want "missing cases in switch of type bar.Phylum: Chordata, Echinodermata"
	case barpkg.Mollusca:
		{
		}
	}
}
