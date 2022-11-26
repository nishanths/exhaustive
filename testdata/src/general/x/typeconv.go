package x

func _t(d Direction) {
	switch d { // want "^missing cases in switch of type x.Direction: x.S, x.directionInvalid$"
	case Direction(N):
	case Direction(E):
	case W:
	}

	switch Direction(d) { // want "^missing cases in switch of type x.Direction: x.S, x.directionInvalid$"
	case N:
	case Direction(E):
	case W:
	}
}
