// want package:"Direction:N,E,S,W,directionInvalid"

package x

func _h() {
	// Same as _a(), but with ignore directive -- so there should be no reporting.

	var d Direction
	//exhaustive:ignore
	switch d {
	case N:
	case S:
	case W:
	default:
	}
}

func _i() {
	// Same as _a(), but with ignore directive -- so there should be no reporting.

	var d Direction
	switch d { //exhaustive:ignore
	case N:
	case S:
	case W:
	default:
	}
}
