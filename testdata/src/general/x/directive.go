package x

func _() {
	var d Direction

	// This comment should not produce an diagnostic (note unknown prefix "exhauster:"
	// instead of "exhaustive:").
	//exhauster:foo
	switch d {
	case N:
	case S:
	case W:
	case E:
	case directionInvalid:
	default:
	}

	// This comment should not produce an diagnostic (note unknown prefix "exhauster:"
	// instead of "exhaustive:").
	//exhauster:foo
	_ = map[Direction]int{
		N:                1,
		S:                2,
		W:                3,
		E:                4,
		directionInvalid: 0,
	}
}
