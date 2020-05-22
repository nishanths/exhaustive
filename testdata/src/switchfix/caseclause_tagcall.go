package switchfix

func directionFunc() Direction {
	return N
}

func _caseclause_tagcall() {
	switch directionFunc() { // want "missing cases in switch of type Direction: E, directionInvalid"
	case N:
	case S:
	case W:
	}
}
