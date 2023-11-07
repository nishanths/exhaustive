package present

import "default-signifies-exhaustive"

func _a(t dse.T) {
	// expect a diagnostic when fDefaultCaseRequired is true.
	switch t { // want "^missing default in switch of type dse.T$"
	case dse.A:
	case dse.B:
	}
}

func _b(t dse.T) {
	//exhaustive:default-require-ignore this is a comment showing that we can turn it off for select switches
	switch t {
	case dse.A:
	case dse.B:
	}
}
