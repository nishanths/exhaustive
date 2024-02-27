package required

import dcr "default-case-required"

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

func _e() {
	// should not report because these are not enum switch
	// statements.
	var x int
	switch x {
	case 0:
	}

	switch {
	case x == 0:
	}
}

func _f(t dcr.T) {
	//exhaustive:enforce-default-case-required
	//exhaustive:ignore-default-case-required
	switch t { // want "^failed to parse directives: conflicting directives \"ignore-default-case-required\" and \"enforce-default-case-required\"$"
	case dcr.A:
	case dcr.B:
	default:
	}
}
