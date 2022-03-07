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
	//exhaustive:ignore
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
		switch d { // want "^missing cases in switch of type Direction: E, directionInvalid$"
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
       switch d { // want "^missing cases in switch of type Direction: E, directionInvalid$"
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
