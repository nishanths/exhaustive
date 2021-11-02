package byvalue

import (
	"duplicateenumvalue"
)

func _p() {
	var r duplicateenumvalue.River

	// should not report missing DefaultRiver, since it has same value as Ganga
	switch r {
	case duplicateenumvalue.Ganga, duplicateenumvalue.Yamuna, duplicateenumvalue.Kaveri:
	}

	// should not report missing Ganga, since it has same value as DefaultRiver
	switch r {
	case duplicateenumvalue.DefaultRiver, duplicateenumvalue.Yamuna, duplicateenumvalue.Kaveri:
	}
}

func _q() {
	var s duplicateenumvalue.State

	// value-based checks not available for iota enums (implementation detail: since
	// we cannot determine a constant.Value from AST/type information).

	switch s { // want "^missing cases in switch of type duplicateenumvalue.State: DefaultState$"
	case duplicateenumvalue.TamilNadu, duplicateenumvalue.Kerala, duplicateenumvalue.Karnataka:
	}
}

func _r() {
	// should report correctly (in union '|' form) when same-valued names are
	// missing.
	var r duplicateenumvalue.River
	switch r { // want "^missing cases in switch of type duplicateenumvalue.River: DefaultRiver|Ganga, Kaveri$"
	case duplicateenumvalue.Yamuna:
	}

	// reporting should work correctly when constant.Values are not present also.
	var s duplicateenumvalue.State
	switch s { // want "^missing cases in switch of type duplicateenumvalue.State: DefaultState, Kerala, TamilNadu$"
	case duplicateenumvalue.Karnataka:
	}
}
