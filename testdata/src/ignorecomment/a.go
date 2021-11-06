package ignorecomment

type Direction int // want Direction:"^N,E,S,W,directionInvalid$"

const (
	N                Direction = 1
	E                Direction = 2
	S                Direction = 3
	W                Direction = 4
	directionInvalid Direction = 5
)
