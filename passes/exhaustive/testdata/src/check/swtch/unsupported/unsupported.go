package unsupported

func xxx(a os) {
	switch a {
	case os(osxtiger):
	case os(osxleopard):
	case os(windows8):
	}
}

func yyy[A osintf](a A) {
	switch a {
	case A(osxtiger):
	case A(osxleopard):
	case A(windows8):
	}
}

type os int

type osintf interface{ darwin | windows }

type darwin os
type windows os

const (
	osxtiger   darwin  = 1
	osxleopard darwin  = 2
	windows7   windows = 3
	windows8   windows = 4
)
