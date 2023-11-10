package notrequired

import "default-case-required"

func _a(t dcr.T) {
	// No diagnostic because default-require is not set.
	switch t {
	case dcr.A:
	case dcr.B:
	}
}

func _b(t dcr.T) {
	//exhaustive:default-require-enforce this is a comment showing that we can turn it on for select switches
	switch t { // want "^missing default in switch of type dcr.T$"
	case dcr.A:
	case dcr.B:
	}
}

func _c(t dcr.T) {
	//exhaustive:default-require-ignore this comment is discarded in facvor of the enforcement
	//exhaustive:default-require-enforce this is a comment showing that we can turn it on for select switches
	switch t { // want "^missing default in switch of type dcr.T$"
	case dcr.A:
	case dcr.B:
	}
}

func _d(t dcr.T) {
	//exhaustive:default-require-enforce this is a comment showing that we can turn it on for select switches
	//exhaustive:default-require-ignore this comment is discarded in facvor of the enforcement
	switch t { // want "^missing default in switch of type dcr.T$"
	case dcr.A:
	case dcr.B:
	}
}

func _e(t dcr.T) {
	//exhaustive:default-require-enforce this is happy because it has a default
	switch t {
	case dcr.A:
	case dcr.B:
	default:
	}
}
