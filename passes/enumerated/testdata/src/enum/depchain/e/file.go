package e

type E1 int // want E1:"^permittedtype$" E1:"^enumerated$" E1:"^elements:V0 = 0, V1 = 1, V2 = 2$"

const (
	V0 E1 = iota // want V0:"^enumerated$" V0:"^elementof:E1$"
	V1           // want V1:"^enumerated$" V1:"^elementof:E1$"
	V2           // want V2:"^enumerated$" V2:"^elementof:E1$"
)
