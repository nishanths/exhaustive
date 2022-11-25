package x

func _k(d Direction) {
	// Parenthesized values in case statements.
	switch d { // want "^missing cases in switch of type x.Direction: x.S, x.directionInvalid$"
	case (N):
	case (E):
	case (((W))):
	}

	// Parenthesized values in switch tag.
	switch ((d)) { // want "^missing cases in switch of type x.Direction: x.S, x.directionInvalid$"
	case N:
	case E:
	case W:
	}

	_ = map[Direction]int{ // want "^missing keys in map of key type x.Direction: x.S, x.directionInvalid$"
		(N): 1,
		(E): 2,
		(((W))): 3,
	}
}
