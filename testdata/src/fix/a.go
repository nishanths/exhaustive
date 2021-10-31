// want package:"^Direction:N,E,S,W,directionInvalid; River:DefaultRiver,Yamuna,Ganga,Kaveri$"

// a.go is a stable first package file for package fact testing.

package fix

type Direction int

var (
	N                Direction = 1
	E                Direction = 2
	S                Direction = 3
	W                Direction = 4
	directionInvalid Direction = 5
)

func ReturnsDirection() Direction {
	return N
}
