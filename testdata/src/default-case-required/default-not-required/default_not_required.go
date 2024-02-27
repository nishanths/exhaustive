package notrequired

import dcr "default-case-required"

func _a(t dcr.T) {
	// No diagnostic because neither fDefaultCaseRequired is true
	// nor the enforcement comment is present.
	switch t {
	case dcr.A:
	case dcr.B:
	}
}

func _b(t dcr.T) {
	//exhaustive:enforce-default-case-required this is a comment showing that we can turn it on for select switches
	switch t { // want "^missing default case in switch of type dcr.T$"
	case dcr.A:
	case dcr.B:
	}
}

func _c(t dcr.T) {
	//exhaustive:enforce-default-case-required this is a comment showing that we can turn it on for select switches
	switch t { // want "^missing default case in switch of type dcr.T$"
	case dcr.A:
	case dcr.B:
	}
}

func _d(t dcr.T) {
	//exhaustive:enforce-default-case-required this is a comment showing that we can turn it on for select switches
	switch t { // want "^missing default case in switch of type dcr.T$"
	case dcr.A:
	case dcr.B:
	}
}

func _e(t dcr.T) {
	//exhaustive:enforce-default-case-required this is happy because it has a default
	switch t {
	case dcr.A:
	case dcr.B:
	default:
	}
}

func _f() {
	// should not report because these are not enum switch
	// statements.
	//exhaustive:enforce-default-case-required
	var x int
	switch x {
	case 0:
	}

	//exhaustive:enforce-default-case-required
	switch {
	case x == 0:
	}
}

func _g(t dcr.T) {
	//exhaustive:ignore-default-case-required
	//exhaustive:enforce-default-case-required
	switch t { // want "^failed to parse directives: conflicting directives \"ignore-default-case-required\" and \"enforce-default-case-required\"$"
	case dcr.A:
	case dcr.B:
	default:
	}
}
