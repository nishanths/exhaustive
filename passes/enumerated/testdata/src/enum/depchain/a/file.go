package a

import _ "enum/depchain/c"

type A1 int // want A1:"^permittedtype$" A1:"^enumerated$" A1:"^elements:Y0 = 0, Y1 = 1$"

const (
	Y0 A1 = iota // want Y0:"^enumerated$" Y0:"^elementof:A1$"
	Y1           // want Y1:"^enumerated$" Y1:"^elementof:A1$"
)
