// want package:"Direction:N,E,S,W,directionInvalid"

// a.go is a stable first package file for package fact testing.

package switchfix

type Direction int

var (
	N                Direction = 1
	E                Direction = 2
	S                Direction = 3
	W                Direction = 4
	directionInvalid Direction = 5
)

func ProducesDirection() Direction {
	return N
}
