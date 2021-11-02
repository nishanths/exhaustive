package x

type Direction int

var (
	N                Direction = 1
	E                Direction = 2
	S                Direction = 3
	W                Direction = 4
	directionInvalid Direction = 5
)

func _a() {
	// Basic same package.

	var d Direction
	switch d {
	case N:
	case S:
	case W:
	}
}
