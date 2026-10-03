package scope

import "test/enum/scope/sub"

const (
	v0 sub.S1 = iota
	v1
)

type t int // want t:"^permittedtype$" t:"^enumerated$" t:"^elements:w0 = 0, w1 = 1$"

const (
	w0 t = iota // want w0:"^enumerated$" w0:"^elementof:t$"
	w1          // want w1:"^enumerated$" w1:"^elementof:t$"
)

func f1() {
	const (
		x0 t = iota
		x1
	)
}

func f2() {
	type q int // want q:"^permittedtype$" q:"^enumerated$" q:"^elements:y0 = 0, y1 = 1, y4 = 4$"
	const (
		y0 q = iota // want y0:"^enumerated$" y0:"^elementof:q$"
		y1          // want y1:"^enumerated$" y1:"^elementof:q$"
	)
	{
		const (
			y2 q = iota + 2
			y3
		)
	}
	const (
		y4 q = iota + 4 // want y4:"^enumerated$" y4:"^elementof:q$"
	)
}
