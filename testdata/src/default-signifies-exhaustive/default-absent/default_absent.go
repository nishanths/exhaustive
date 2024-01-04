package absent

import dse "default-signifies-exhaustive"

func _a(t dse.T) {
	switch t { // want "^missing cases in switch of type dse.T: dse.B$"
	case dse.A:
	}
}

func _b(t dse.T) {
	//exhaustive:ignore
	//exhaustive:enforce
	switch t { // want "^failed to parse directives: conflicting directives \"ignore\" and \"enforce\"$"
	case dse.A:
	case dse.B:
	}
}
