package fix

func _a() {
	var d Direction
	switch d { // want "^missing cases in switch of type Direction: E, directionInvalid$"
	case N:
	case S:
	case W:
		_ = 3
	default:
	}
}
