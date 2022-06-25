package otherpkg

import (
	d "duplicate-enum-value"
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
	switch s { // want "^missing cases in switch of type duplicateenumvalue.State: TamilNadu\\|DefaultState, Kerala$"
	case d.Karnataka:
	}
}

func _s(c d.Chart) {
	switch c { // want "^missing cases in switch of type duplicateenumvalue.Chart: Pie$"
	case d.Line:
	case d.Sunburst:
	case d.Area:
	}
}
