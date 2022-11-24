package absent

import "default-signifies-exhaustive"

func _a(t dse.T) {
	switch t { // want "^missing cases in switch of type dse.T: dse.B$"
	case dse.A:
	}
}
