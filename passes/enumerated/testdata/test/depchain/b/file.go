package b

type B1 int // want B1:"^permittedtype$" B1:"^enumerated$" B1:"^elements:BY0 = 0, BY1 = 1$"

const (
	BY0 B1 = iota // want BY0:"^enumerated$" BY0:"^elementof:B1$"
	BY1           // want BY1:"^enumerated$" BY1:"^elementof:B1$"
)

const (
	BX0 int = iota
	BX1
)

type B2 int // want B2:"^permittedtype$"
