package e

type E1 int // want E1:"^permittedtype$" E1:"^enumerated$" E1:"^elements:EV0 = 0, EV1 = 1, EV2 = 2$"

const (
	EV0 E1 = iota // want EV0:"^enumerated$" EV0:"^elementof:E1$"
	EV1           // want EV1:"^enumerated$" EV1:"^elementof:E1$"
	EV2           // want EV2:"^enumerated$" EV2:"^elementof:E1$"
)
