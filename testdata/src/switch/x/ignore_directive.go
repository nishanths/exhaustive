// want package:"Direction:N,E,S,W,directionInvalid"

package x

func _h() {
	// Same as _a(), but with ignore directive -- so there should be no reporting.
	// (Verified manually, by removing directive, that the diagnostic is reported
	// when the directive does not exist.)

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
	// (Verified manually, by removing directive, that the diagnostic is reported
	// when the directive does not exist.)

	var d Direction
	switch d { //exhaustive:ignore
	case N:
	case S:
	case W:
	default:
	}
}
