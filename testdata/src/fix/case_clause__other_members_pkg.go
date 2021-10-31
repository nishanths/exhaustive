package fix

import (
	barpkg "general/y"
)

func _d() {
	var phy barpkg.Phylum
	switch phy { // want "^missing cases in switch of type bar.Phylum: Chordata, Echinodermata$"
	case barpkg.Mollusca:
		{
		}
	}
}
