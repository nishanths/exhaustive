package switchfix

func _fmt_nogendecl() {
	var d Direction
	switch d { // want "missing cases in switch of type Direction: E, directionInvalid"
	case N:
	case S:
	case W:
	}
}
