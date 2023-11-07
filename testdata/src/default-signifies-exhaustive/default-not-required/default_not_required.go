package present

import "default-signifies-exhaustive"

func _a(t dse.T) {
	// No diagnostic because default-require is not set.
	switch t {
	case dse.A:
	case dse.B:
	}
}

func _b(t dse.T) {
	//exhaustive:default-require-enforce this is a comment showing that we can turn it on for select switches
	switch t { // want "^missing default in switch of type dse.T$"
	case dse.A:
	case dse.B:
	}
}

func _c(t dse.T) {
	//exhaustive:default-require-enforce this is happy because it has a default
	switch t {
	case dse.A:
	case dse.B:
	default:
	}
}
