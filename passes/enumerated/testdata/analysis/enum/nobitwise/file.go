package nobitwise

func f1() {
	type t1 int // want t1:"^permittedtype$" t1:"^enumerated$" t1:"^elements:x0 = 0, x1 = 1$"
	const (
		x0 t1 = 0 // want x0:"^enumerated$" x0:"^elementof:t1"
		x1 t1 = 1 // want x1:"^enumerated$" x1:"^elementof:t1"
	)

	type t2 int // want t2:"^permittedtype$"
	const (
		y0 t2 = ^0
		y1 t2 = ^1
	)

	type t3 int // want t3:"^permittedtype$"
	const (
		z0 t3 = ^0
		z1 t3 = t3(y1)
	)

	type t4 int // want t4:"^permittedtype$" t4:"^enumerated$" t4:"^elements:w0 = 0, w2 = 2$"
	const (
		w0 t4 = 0 // want w0:"^enumerated$" w0:"^elementof:t4"
		w1 t4 = ^0
		w2 t4 = 2 // want w2:"^enumerated$" w2:"^elementof:t4"
	)
}
