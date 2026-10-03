package a

import _ "test/depchain/b"
import _ "test/depchain/c"

type A1 int // want A1:"^permittedtype$" A1:"^enumerated$" A1:"^elements:AX0 = 0, AX1 = 1, AX2 = 2$"

const (
	AX0 A1 = iota // want AX0:"^enumerated$" AX0:"^elementof:A1$"
	AX1           // want AX1:"^enumerated$" AX1:"^elementof:A1$"
	AX2           // want AX2:"^enumerated$" AX2:"^elementof:A1$"
)
