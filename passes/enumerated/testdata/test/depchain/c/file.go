package c

import "test/depchain/d"

type C1 = d.D1

const (
	CX0 = d.DX0
	CX1 = d.DX1
)

type C2 int // want C2:"^permittedtype$" C2:"^enumerated$" C2:"^elements:CY0 = 0, CY1 = 1$"

const (
	CY0 C2 = iota // want CY0:"^enumerated$" CY0:"^elementof:C2$"
	CY1           // want CY1:"^enumerated$" CY1:"^elementof:C2$"
)
