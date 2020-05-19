package explay

type Dir int

const (
	N Dir = iota + 1
	E
	S
	W
)

func foo() {
	var d Dir
	switch d {
	case N:
	case E:
	case W:
	}
}
