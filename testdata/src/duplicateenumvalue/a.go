package duplicateenumvalue

type River string // want River:"^DefaultRiver,Yamuna,Ganga,Kaveri$"

const DefaultRiver = Ganga

const (
	Yamuna River = "Yamuna"
	Ganga  River = "Ganga"
	Kaveri River = "Kaveri"
)

type State int // want State:"^TamilNadu,Kerala,Karnataka,DefaultState$"

const (
	_ State = iota
	TamilNadu
	Kerala
	Karnataka
)

const DefaultState = TamilNadu
