package iotavalue

import "enum/typ"

type t int // want t:"^permittedtype$" t:"^enumerated$" t:"^elements:x0 = 0, x1 = 1$"

const (
	x0 t = iota // want x0:"^enumerated$" x0:"^elementof:t$"
	x1          // want x1:"^enumerated$" x1:"^elementof:t$"
)

type q int // want q:"^permittedtype$"

const (
	y1 q = 1
	y2 q = 2
)

func f1() {
	type s int // want s:"^permittedtype$" s:"^enumerated$" s:"^elements:m0 = 0, m1 = 1$"
	const (
		m0 s = iota // want m0:"^enumerated$" m0:"^elementof:s$"
		m1          // want m1:"^enumerated$" m1:"^elementof:s$"
	)

	type u int // want u:"^permittedtype$" u:"^enumerated$" u:"^elements:l0 = 0, l2 = 6, l6 = 18$"
	const (
		l0 u = iota * 3 // want l0:"^enumerated$" l0:"^elementof:u$"
		_
		l2   // want l2:"^enumerated$" l2:"^elementof:u$"
		k0 u = 0
		k1 u = 1
		k2 u = 2
		l6 u = iota * 3 // want l6:"^enumerated$" l6:"^elementof:u$"
	)

	const (
		a0 int = iota
		a1
	)

	type w int // want w:"^permittedtype$" w:"^enumerated$" w:"^elements:b0 = 0, b1 = 1$"
	const (
		b0 w = iota // want b0:"^enumerated$" b0:"^elementof:w$"
		b1          // want b1:"^enumerated$" b1:"^elementof:w$"
	)

	type v int // want v:"^permittedtype$"
	const (
		i0 v = v(a0)
		i1 v = v(a1)
	)

	type p int // want p:"^permittedtype$" p:"^enumerated$" p:"^elements:j0 = 0, j1 = 1, j2 = 2$"
	const (
		j0 p = p(b0) // want j0:"^enumerated$" j0:"^elementof:p$"
		j1 p = p(b1) // want j1:"^enumerated$" j1:"^elementof:p$"
		j2 p = iota  // want j2:"^enumerated$" j2:"^elementof:p$"
	)

	type q int // want q:"^permittedtype$" q:"^enumerated$" q:"^elements:h1 = 1, h2 = 2, h3 = 3, h4 = 4$"
	const (
		h1 q = q(typ.V0) + 1 // want h1:"^enumerated$" h1:"^elementof:q$"
		h2 q = q(typ.V1) + 1 // want h2:"^enumerated$" h2:"^elementof:q$"
		h3 q = q(typ.V2) + 1 // want h3:"^enumerated$" h3:"^elementof:q$"
		h4 q = q(typ.V3) + 1 // want h4:"^enumerated$" h4:"^elementof:q$"
	)
}

func f2() {
	const iota = 9999 // shadows predeclared identifier
	type s int        // want s:"^permittedtype$"
	const (
		m0 s = iota
		m1
	)

}
