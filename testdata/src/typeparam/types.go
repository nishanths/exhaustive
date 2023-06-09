package typeparam

type Stringer interface {
	String() string
}

type M uint8 // want M:"^A,B$"
const (
	_ M = iota * 100
	A
	B
)

func (M) String() string { return "" }

type N uint8 // want N:"^C,D$"
const (
	_ N = iota * 100
	C
	D
)

type O byte // want O:"^E1,E2$"
const (
	E1 O = 'h'
	E2 O = 'e'
)

type P float32 // want P:"^F$"
const (
	F P = 1.1234
)

type Q string // want Q:"^G$"
const (
	G Q = "world"
)

type NotEnumType uint8

type II interface{ N | JJ }
type JJ interface{ O }
type KK interface {
	M
	Stringer
	error
	comparable
}
type LL interface {
	M | NotEnumType
	Stringer
	error
}
type MM interface {
	M
}
type Anon interface {
	interface{ M } | interface{ N }
}
type QQ interface {
	Q
}
