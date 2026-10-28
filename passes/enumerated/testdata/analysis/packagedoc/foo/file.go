package foo

type T0 int // want T0:"^permittedtype$" T0:"^enumerated$" T0:"^elements:X0 = 0, X1 = 1, X2 = 2, X02 = 2$"

const (
	X0 T0 = iota // want X0:"^enumerated$" X0:"^elementof:T0$"
	X1           // want X1:"^enumerated$" X1:"^elementof:T0$"
	X2           // want X2:"^enumerated$" X2:"^elementof:T0$"
)

const X02 = X2 // want X02:"^enumerated$" X02:"^elementof:T0$"
