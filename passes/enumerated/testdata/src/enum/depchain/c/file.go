package c

import "enum/depchain/d"

type C1 = d.D1

const (
	X0 = d.X0
	X1 = d.X1
)

type C2 int // want C2:"^permittedtype$" C2:"^enumerated$" C2:"^elements:Y0 = 0, Y1 = 1$"

const (
	Y0 C2 = iota // want Y0:"^enumerated$" Y0:"^elementof:C2$"
	Y1           // want Y1:"^enumerated$" Y1:"^elementof:C2$"
)
