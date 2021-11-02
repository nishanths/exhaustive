// want package:"^River:DefaultRiver,Yamuna,Ganga,Kaveri; State:_,TamilNadu,Kerala,Karnataka,DefaultState$"

package duplicateenumvalue

type River string

const DefaultRiver = Ganga

const (
	Yamuna River = "Yamuna"
	Ganga  River = "Ganga"
	Kaveri River = "Kaveri"
)

type State int

const (
	_ State = iota
	TamilNadu
	Kerala
	Karnataka
)

const DefaultState = TamilNadu
