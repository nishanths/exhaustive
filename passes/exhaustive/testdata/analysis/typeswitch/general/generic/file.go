package generic

type n interface{ n() }
type m interface{ m() }
type mo interface {
	m()
	o()
}

type u struct{}
func (*u) u()

type w struct{}
func (*w) n()

// Note: generic types that implement an interface are not supported.

type l[P1 m, P2 int] struct{}
func (*l[P1, P2]) n()

type (
	a0[P1, P2 any] = w
	a1[P1, P2 any] = l[m, int]

	a2[P1 m]         = l[P1, int]
	a3[P1 m, P2 int] = l[P1, P2]
	a4               = l[m, int]
	a5               = a3[m, int]
	a6[P1 m, P2 int] = a3[P1, P2]
)

func f1() {
	var y n

	switch y.(type) { // want "^missing cases in type switch: \\*w or \\*a0$"
	}

	switch y.(type) {
	case *w:
	}

	switch y.(type) {
	case *a0[string, string]: // note: inconsequential type arguments
	}
}
