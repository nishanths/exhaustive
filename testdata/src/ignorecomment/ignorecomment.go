package ignorecomment

func _a() {
	var d Direction

	// this should not report.
	// some other comment
	//exhaustive:ignore
	// some other comment
	switch d {
	case N:
	case S:
	case W:
	default:
	}

	// this should report.
	switch d { // want "^missing cases in switch of type Direction: E, directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}
}

func _b() {
	var d Direction

	// this should not report.
	switch d { //exhaustive:ignore
	case N:
	case S:
	case W:
	default:
	}

	// this should report.
	switch d { // want "^missing cases in switch of type Direction: E, directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}
}
