package absent

import dse "default-signifies-exhaustive"

func _a(t dse.T) {
	switch t { // want "^missing cases in switch of type dse.T: B$"
	case dse.A:
	}
}
