package enforcecomment

func _a() {
	var d Direction

	switch d {
	case N:
	case S:
	case W:
	default:
	}

	// this should report.
	// some other comment
	//exhaustive:enforce
	// some other comment
	switch d { // want "^missing cases in switch of type enforcecomment.Direction: enforcecomment.E, enforcecomment.directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}
}

func _b() {
	var d Direction

	// this should not report.
	switch d {
	case N:
	case S:
	case W:
	default:
	}

	// this should report
	//exhaustive:enforce
	switch d { // want "^missing cases in switch of type enforcecomment.Direction: enforcecomment.E, enforcecomment.directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}
}

func _nested() {
	var d Direction

	// this should not report.
	switch d {
	case N:
	case S:
	case W:
	default:
		// this should report.
		//exhaustive:enforce
		switch d { // want "^missing cases in switch of type enforcecomment.Direction: enforcecomment.E, enforcecomment.directionInvalid$"
		case N:
		case S:
		case W:
		default:
		}
	}
}

func _reverse_nested() {
	var d Direction

	// this should report.
	//exhaustive:enforce
	switch d { // want "^missing cases in switch of type enforcecomment.Direction: enforcecomment.E, enforcecomment.directionInvalid$"
	case N:
	case S:
	case W:
	default:
		// this should not report.
		switch d {
		case N:
		case S:
		case W:
		default:
		}
	}
}
