package a

import _ "depchain/b"
import _ "depchain/c"

type A1 int // want A1:"^permittedtype$" A1:"^enumerated$" A1:"^elements:AY0 = 0, AY1 = 1$"

const (
	AY0 A1 = iota // want AY0:"^enumerated$" AY0:"^elementof:A1$"
	AY1           // want AY1:"^enumerated$" AY1:"^elementof:A1$"
)
