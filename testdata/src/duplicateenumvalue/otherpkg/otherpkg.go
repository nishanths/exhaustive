package otherpkg

import (
	d "duplicateenumvalue"
)

func _p() {
	var r d.River

	// should not report missing DefaultRiver, since it has same value as Ganga
	switch r {
	case d.Ganga, d.Yamuna, d.Kaveri:
	}

	// should not report missing Ganga, since it has same value as DefaultRiver
	switch r {
	case d.DefaultRiver, d.Yamuna, d.Kaveri:
	}
}

func _q() {
	var s d.State

	// should not report missing DefaultState, since it has same value as TamilNadu
	switch s {
	case d.TamilNadu, d.Kerala, d.Karnataka:
	}
}

func _r() {
	// should report correctly (in union '|' form) when same-valued names are
	// missing.

	var r d.River
	switch r { // want "^missing cases in switch of type duplicateenumvalue.River: DefaultRiver\\|Ganga, Kaveri$"
	case d.Yamuna:
	}

	var s d.State
	switch s { // want "^missing cases in switch of type duplicateenumvalue.State: DefaultState\\|TamilNadu, Kerala$"
	case d.Karnataka:
	}
}
