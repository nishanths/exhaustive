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
	switch d { // want "^missing cases in switch of type ignorecomment.Direction: ignorecomment.E, ignorecomment.directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}
}

func _b() {
	var d Direction

	// this should not report.
	//exhaustive:ignore ... more arbitrary comment content (e.g. an explanation) ...
	switch d {
	case N:
	case S:
	case W:
	default:
	}

	// this should report.
	switch d { // want "^missing cases in switch of type ignorecomment.Direction: ignorecomment.E, ignorecomment.directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}
}

func _c0() {
	var d Direction

	// this should report: according to go/ast, the comment is not considered to
	// be associated with the switch statement node.
	switch d { //exhaustive:ignore // want "^missing cases in switch of type ignorecomment.Direction: ignorecomment.E, ignorecomment.directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}

	// this should report.
	switch d { // want "^missing cases in switch of type ignorecomment.Direction: ignorecomment.E, ignorecomment.directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}
}

func _c1() {
	var d Direction

	// this should report: according to go/ast, the comment is not considered to
	// be associated with the switch statement node.
	switch d { // want "^missing cases in switch of type ignorecomment.Direction: ignorecomment.E, ignorecomment.directionInvalid$"
	//exhaustive:ignore
	case N:
	case S:
	case W:
	default:
	}

	// this should report.
	switch d { // want "^missing cases in switch of type ignorecomment.Direction: ignorecomment.E, ignorecomment.directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}
}

func _d() {
	// this should report.
	switch (func() Direction { // want "^missing cases in switch of type ignorecomment.Direction: ignorecomment.E, ignorecomment.directionInvalid$"
		// this should not report.
		var x Direction
		//exhaustive:ignore
		switch x {
		case N, S:
		}
		return N
	})() {
	case N:
	case S:
	case W:
	default:
	}

	var d Direction

	// this should report.
	switch d { // want "^missing cases in switch of type ignorecomment.Direction: ignorecomment.E, ignorecomment.directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}
}

func _nested() {
	var d Direction

	// this should not report.
	//exhaustive:ignore
	switch d {
	case N:
	case S:
	case W:
	default:
		// this should report.
		switch d { // want "^missing cases in switch of type ignorecomment.Direction: ignorecomment.E, ignorecomment.directionInvalid$"
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
	switch d { // want "^missing cases in switch of type ignorecomment.Direction: ignorecomment.E, ignorecomment.directionInvalid$"
	case N:
	case S:
	case W:
	default:
		// this should not report.
		//exhaustive:ignore
		switch d {
		case N:
		case S:
		case W:
		default:
		}
	}
}
