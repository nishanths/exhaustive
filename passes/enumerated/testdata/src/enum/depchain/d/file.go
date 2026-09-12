package d

type D1 int // want D1:"^permittedtype$" D1:"^enumerated$" D1:"^elements:X0 = 0, X1 = 1$"

const (
	X0 D1 = iota // want X0:"^enumerated$" X0:"^elementof:D1$"
	X1           // want X1:"^enumerated$" X1:"^elementof:D1$"
)

type D2 int // want D2:"^permittedtype$" D2:"^enumerated$" D2:"^elements:Z0 = 0, Z1 = 1, Z2 = 2$"

const (
	Z0 D2 = iota // want Z0:"^enumerated$" Z0:"^elementof:D2$"
	Z1           // want Z1:"^enumerated$" Z1:"^elementof:D2$"
	Z2           // want Z2:"^enumerated$" Z2:"^elementof:D2$"
)
