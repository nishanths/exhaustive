package exp

// NOTE: Feel free to delete code below. This is a temporary scratchpad for use
// when working on bugfixes or new code for exhaustive.

import "github.com/nishanths/exhaustive/testdata/playground/x"

func _o() {
	// Selector isn't of the form "enumPkg.enumMember"

	type holdsW struct {
		W x.Direction
	}

	var d x.Direction
	var h holdsW

	switch d {
	case x.N:
	case x.E:
	case x.S:
	case h.W:
	}
}
