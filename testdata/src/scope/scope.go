package scope

type T int // want T:"^A,B$"

const (
	A T = iota
	B
)

type Q string // want Q:"^X,Y$"

const (
	X Q = "x"
	Y Q = "y"
)

func _a() {
	type T int // want T:"^C,D$"

	const (
		C T = iota
		D
	)

	var t T
	// must not report diagnostic here
	switch t {
	case C, D:
	}

	switch t { // want "^missing cases in switch of type T: D$"
	case C:
	}

	type Q string // want Q:"^X,Y$"

	const (
		X Q = "x"
		Y Q = "y"
	)

	var q Q
	// must not report diagnostic here
	switch q {
	case X, Y:
	}

	switch q { // want "^missing cases in switch of type Q: Y$"
	case X:
	}
}
