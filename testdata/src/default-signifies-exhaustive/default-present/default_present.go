package present

import "default-signifies-exhaustive"

func _a(t dse.T) {
	// expect no diagnostics, since default case is present,
	// even though members are missing in the switch.

	switch t {
	case dse.A:
	default:
	}
}
