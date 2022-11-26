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

	// should not report missing DefaultRiver, since it has same value as Ganga
	_ = map[d.River]int{
		d.Ganga:  1,
		d.Yamuna: 2,
		d.Kaveri: 3,
	}

	// should not report missing Ganga, since it has same value as DefaultRiver
	_ = map[d.River]int{
		d.DefaultRiver: 1,
		d.Yamuna:       2,
		d.Kaveri:       3,
	}
}

func _q() {
	var s d.State

	// should not report missing DefaultState, since it has same value as TamilNadu
	switch s {
	case d.TamilNadu, d.Kerala, d.Karnataka:
	}

	// should not report missing DefaultState, since it has same value as TamilNadu
	_ = map[d.State]int{
		d.TamilNadu: 1,
		d.Kerala:    2,
		d.Karnataka: 3,
	}
}

func _r() {
	// should report correctly (in union '|' form) when same-valued names are
	// missing.

	var r d.River
	switch r { // want `^missing cases in switch of type duplicateenumvalue.River: duplicateenumvalue.DefaultRiver\|duplicateenumvalue.Ganga, duplicateenumvalue.Kaveri$`
	case d.Yamuna:
	}

	var s d.State
	switch s { // want `^missing cases in switch of type duplicateenumvalue.State: duplicateenumvalue.TamilNadu\|duplicateenumvalue.DefaultState, duplicateenumvalue.Kerala$`
	case d.Karnataka:
	}

	_ = map[d.River]int{ // want `^missing keys in map of key type duplicateenumvalue.River: duplicateenumvalue.DefaultRiver\|duplicateenumvalue.Ganga, duplicateenumvalue.Kaveri$`
		d.Yamuna: 1,
	}

	_ = map[d.State]int{ // want `^missing keys in map of key type duplicateenumvalue.State: duplicateenumvalue.TamilNadu\|duplicateenumvalue.DefaultState, duplicateenumvalue.Kerala$`
		d.Karnataka: 1,
	}
}

func _s(c d.Chart) {
	switch c { // want "^missing cases in switch of type duplicateenumvalue.Chart: duplicateenumvalue.Pie$"
	case d.Line:
	case d.Sunburst:
	case d.Area:
	}

	_ = map[d.Chart]int{ // want "^missing keys in map of key type duplicateenumvalue.Chart: duplicateenumvalue.Pie$"
		d.Line:     1,
		d.Sunburst: 2,
		d.Area:     3,
	}
}
