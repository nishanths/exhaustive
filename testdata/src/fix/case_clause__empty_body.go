package fix

import (
	bar "general/y"
)

func _b() {
	var phy bar.Phylum
	switch phy { // want "^missing cases in switch of type bar.Phylum: Chordata, Echinodermata, Mollusca$"
	}
}
