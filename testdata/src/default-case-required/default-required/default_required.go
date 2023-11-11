package required

import "default-case-required"

func _a(t dcr.T) {
	// expect a diagnostic when fDefaultCaseRequired is true.
	switch t { // want "^missing default case in switch of type dcr.T$"
	case dcr.A:
	case dcr.B:
	}
}

func _b(t dcr.T) {
	//exhaustive:ignore-default-case-required this is a comment showing that we can turn it off for select switches
	switch t {
	case dcr.A:
	case dcr.B:
	}
}

func _c(t dcr.T) {
	//exhaustive:ignore-default-case-required this comment is discarded in facvor of the enforcement
	//exhaustive:enforce-default-case-required this helps override the above
	switch t { // want "^missing default case in switch of type dcr.T$"
	case dcr.A:
	case dcr.B:
	}
}

func _d(t dcr.T) {
	// this is happy even with enforcement because we have a default
	switch t {
	case dcr.A:
	case dcr.B:
	default:
	}
}
