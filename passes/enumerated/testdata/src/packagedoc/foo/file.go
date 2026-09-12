package foo

type T1 int // want T1:"^permittedtype$" T1:"^enumerated$" T1:"^elements:X0 = 0, X1 = 1, X2 = 2, XTWO = 2$"

const (
	X0   T1   = iota // want X0:"^enumerated$" X0:"^elementof:T1$"
	X1               // want X1:"^enumerated$" X1:"^elementof:T1$"
	X2               // want X2:"^enumerated$" X2:"^elementof:T1$"
	XTWO = X2        // want XTWO:"^enumerated$" XTWO:"^elementof:T1$"
)
