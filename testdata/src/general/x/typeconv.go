package x

func _t(d Direction) {
	switch d { // want "^missing cases in switch of type x.Direction: x.S, x.directionInvalid$"
	case Direction(N):
	case Direction(int(E)):
	case W:
	}

	switch Direction(int(d)) { // want "^missing cases in switch of type x.Direction: x.S, x.directionInvalid$"
	case N:
	case Direction(Direction(E)):
	case Direction(W):
	}

	switch Direction(int(d)) { // want "^missing cases in switch of type x.Direction: x.N, x.S, x.W, x.directionInvalid$"
	case (_tt{}).methodCallMe(N):
	case Direction(E):
	case callMe(W):
	}

	var _ = map[Direction]struct{}{ // want "^missing keys in map of key type x.Direction: x.N, x.S, x.W, x.directionInvalid$"
		(_tt{}).methodCallMe(N): struct{}{},
		Direction(E):            struct{}{},
		callMe(W):               struct{}{},
	}
}

func callMe(d Direction) Direction { return d }

type _tt struct{}

func (_tt) methodCallMe(d Direction) Direction { return d }
