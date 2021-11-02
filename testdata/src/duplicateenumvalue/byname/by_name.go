package byname

import (
	"duplicateenumvalue"
)

// Compare to ../byvalue/by_value.go

func _p() {
	var r duplicateenumvalue.River

	switch r { // want "^missing cases in switch of type duplicateenumvalue.River: DefaultRiver$"
	case duplicateenumvalue.Ganga, duplicateenumvalue.Yamuna, duplicateenumvalue.Kaveri:
	}

	switch r { // want "^missing cases in switch of type duplicateenumvalue.River: Ganga$"
	case duplicateenumvalue.DefaultRiver, duplicateenumvalue.Yamuna, duplicateenumvalue.Kaveri:
	}
}

func _q() {
	var s duplicateenumvalue.State

	switch s { // want "^missing cases in switch of type duplicateenumvalue.State: DefaultState$"
	case duplicateenumvalue.TamilNadu, duplicateenumvalue.Kerala, duplicateenumvalue.Karnataka:
	}
}

func _r() {
	var r duplicateenumvalue.River
	switch r { // want "^missing cases in switch of type duplicateenumvalue.River: DefaultRiver, Ganga, Kaveri$"
	case duplicateenumvalue.Yamuna:
	}

	var s duplicateenumvalue.State
	switch s { // want "^missing cases in switch of type duplicateenumvalue.State: DefaultState, Kerala, TamilNadu$"
	case duplicateenumvalue.Karnataka:
	}
}
