package switchfix

import "switch/y"

func _tag_callexpr_func() {
	switch ProducesDirection() { // want "^missing cases in switch of type Direction: S$"
	case N, E, W, directionInvalid:
	}
}

func _tag_callexpr_typeconversion() {
	// no function or method calls -- should add case clause for missing enum members

	var i int
	switch Direction(bar.IntWrapper(i)) { // want "^missing cases in switch of type Direction: S$"
	case N, E, W, directionInvalid:
	}
}

func _tag_callexpr_mix() {
	switch Direction(bar.IntWrapper(int(ProducesDirection()))) { // want "^missing cases in switch of type Direction: S$"
	case N, E, W, directionInvalid:
	}
}

func _tag_callexpr_builtin() {
	var a, b []int
	switch Direction(copy(a, b)) { // want "^missing cases in switch of type Direction: S$"
	case N, E, W, directionInvalid:
	}
}
