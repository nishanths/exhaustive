package switchfix

import (
	bar "switch/y"
)

func _caseclause_emptybody() {
	var phy bar.Phylum
	switch phy { // want "^missing cases in switch of type bar.Phylum: Chordata, Echinodermata, Mollusca$"
	}
}
