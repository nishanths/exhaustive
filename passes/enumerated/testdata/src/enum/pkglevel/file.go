package pkglevel

type t int // want t:"^permittedtype$" t:"^enumerated$" t:"^elements:w0 = 0, w1 = 1$"

const (
	w0 t = iota // want w0:"^enumerated$" w0:"^elementof:t$"
	w1          // want w1:"^enumerated$" w1:"^elementof:t$"
)

func f1() {
	type q int
	const (
		x0 q = iota
		x1
	)
}

func f2() {
	const (
		y0 t = iota
		y1
	)
}
