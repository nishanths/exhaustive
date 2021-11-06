package pkgscope

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
	type T int

	const (
		C T = iota
		D
	)

	var t T
	switch t {
	case C, D:
	}

	switch t {
	case C:
	}

	type Q string

	const (
		X Q = "x"
		Y Q = "y"
	)

	var q Q
	switch q {
	case X, Y:
	}

	switch q {
	case X:
	}
}
