package x

func _k(d Direction) {	
	// Parenthesized values in case statements.
	switch d { // want "^missing cases in switch of type Direction: S, directionInvalid$"
	case (N):
	case (E):
	case (((W))):
	}

	// Parenthesized values in switch tag.
	switch ((d)) { // want "^missing cases in switch of type Direction: S, directionInvalid$"
	case N:
	case E:
	case W:
	}
}
