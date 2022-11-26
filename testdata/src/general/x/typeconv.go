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
}

func callMe(d Direction) Direction { return d }

type _tt struct{}

func (_tt) methodCallMe(d Direction) Direction { return d }
