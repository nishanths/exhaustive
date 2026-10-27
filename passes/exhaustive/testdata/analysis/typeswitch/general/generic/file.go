package generic

type n interface{ n() }
type m interface{ m() }
type mo interface {
	m()
	o()
}
type mno interface {
	m()
	n()
	o()
}

type u struct{}
type w struct{}
type l0[P1 m, P2 int] l1[P1, P2]
type l1[P1, P2 any] struct{}
type l2[P1, P2 any] struct{}

func (*u) u()
func (*w) n()
func (l0[P1, P2]) n()
func (*l1[A, B]) n()

type (
	a00             = w
	a01[P1, P2 any] = w // inconsequential type parameters
)

type (
	b00[P1 m, P2 int]      = l0[P1, P2]
	b01[P1 m, P2 int]      = *l0[P1, P2]
	b02                    = l0[m, int]
	b03                    = *l0[m, int]
	b04[P1 mo, P2 float64] = l0[m, int]
	b05[P1 m]              = l0[P1, int]

	c00[P1 mo, P2 int]          = l0[P1, P2]
	c01[P1 mo, P2 int]          = *l0[P1, P2]
	c02                         = l0[mo, int]
	c03                         = *l0[mo, int]
	c04[P1 float64, P2 float64] = l0[mo, int]
	c05[P1 mo]                  = l0[P1, int]

	b10[P1 m, P2 int] = b00[P1, P2]
	b11[P1 m, P2 int] = b01[P1, int]
	c12[P1 m, P2 int] = *b01[P1, int]
	b13               = b00[m, int]
	b14               = b01[m, int]
	b15               = b02
	b16               = *b02
	b17               = b03
	c18               = *b03
	b19               = b04[mo, float64]
	b20               = b05[m]
	c21               = b00[mo, int]
	c22               = b05[mo]

	b30 = b02
	b31 = *b02
	b32 = b03
	c33 = *b03

	c40[P1 mo, P2 int] = l0[P1, P2]
	c41                = c40[mno, int]

	d00[P1, P2 any] = *l1[P1, P2]
	d01             = *l1[any, any]
	e00             = *l1[int, int]
)

func f1() {
	var y n

	switch y.(type) { // want "^missing cases in type switch: \\*w or \\*a00 or \\*a01\\[…\\]$"
	}

	switch y.(type) {
	case *w:
	}

	switch y.(type) {
	case *a00:
	}

	switch y.(type) {
	case *a01[float64, float64]: // inconsequential type arguments
	}

	switch y.(type) {
	case l0[m, int]:
	case *l0[m, int]:
	case *l1[any, any]:
	case *w:
	}

	switch y.(type) {
	case b00[m, int]:
	case b01[m, int]:
	case d00[any, any]:
	case *a00:
	}

	switch y.(type) {
	case b05[m]:
	case *b19:
	case d01:
	case *a01[float64, float64]: // inconsequential type arguments
	}
}

var (
	_ a01[any, any]         = w{}
	_ a01[float64, float64] = w{}
	_ a01[any, any]         = a01[float64, float64]{}
)
