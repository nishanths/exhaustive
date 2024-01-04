package present

import dse "default-signifies-exhaustive"

func _a(t dse.T) {
	// expect no diagnostics, since default case is present,
	// even though members are missing in the switch.

	switch t {
	case dse.A:
	default:
	}
}

func _b(t dse.T) {
	//exhaustive:enforce
	//exhaustive:ignore
	switch t { // want "^failed to parse directives: conflicting directives \"ignore\" and \"enforce\"$"
	case dse.A:
	case dse.B:
	}
}
