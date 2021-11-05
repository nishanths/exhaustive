package pkg

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

	// should not report missing DefaultState, since it has same value as TamilNadu
	switch s {
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

	var s duplicateenumvalue.State
	switch s { // want "^missing cases in switch of type duplicateenumvalue.State: DefaultState|TamilNadu, Kerala$"
	case duplicateenumvalue.Karnataka:
	}
}
