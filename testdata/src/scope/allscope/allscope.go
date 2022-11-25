package allscope

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

	switch t { // want "^missing cases in switch of type allscope.T: allscope.D$"
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

	switch q { // want "^missing cases in switch of type allscope.Q: allscope.Y$"
	case X:
	}
}

func _b() {
	type T int // want T:"^C,D$"

	const (
		C T = iota
		D
	)

	// must not report diagnostic here
	_ = map[T]int{
		C: 1,
		D: 2,
	}

	_ = map[T]int{ // want "^missing keys in map of key type allscope.T: allscope.D$"
		C: 1,
	}

	type Q string // want Q:"^X,Y$"

	const (
		X Q = "x"
		Y Q = "y"
	)

	// must not report diagnostic here
	_ = map[Q]int{
		X: 1,
		Y: 2,
	}

	_ = map[Q]int{ // want "^missing keys in map of key type allscope.Q: allscope.Y$"
		X: 1,
	}
}
