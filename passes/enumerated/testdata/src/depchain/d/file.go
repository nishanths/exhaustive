package d

type D1 int // want D1:"^permittedtype$" D1:"^enumerated$" D1:"^elements:DX0 = 0, DX1 = 1$"

const (
	DX0 D1 = iota // want DX0:"^enumerated$" DX0:"^elementof:D1$"
	DX1           // want DX1:"^enumerated$" DX1:"^elementof:D1$"
)

type D2 int // want D2:"^permittedtype$" D2:"^enumerated$" D2:"^elements:DZ0 = 0, DZ1 = 1, DZ2 = 2$"

const (
	DZ0 D2 = iota // want DZ0:"^enumerated$" DZ0:"^elementof:D2$"
	DZ1           // want DZ1:"^enumerated$" DZ1:"^elementof:D2$"
	DZ2           // want DZ2:"^enumerated$" DZ2:"^elementof:D2$"
)
