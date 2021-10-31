package fix

import (
	bar "general/y"
)

func _c() {
	var phy bar.Phylum
	switch phy { // want "^missing cases in switch of type bar.Phylum: Chordata, Echinodermata, Mollusca$"
	default:
		print("...")
	}
}
