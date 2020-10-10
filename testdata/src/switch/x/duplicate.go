// want package:"Direction:N,E,S,W,directionInvalid; River:DefaultRiver,Yamuna,Ganga,Kaveri; SortDirection:_,Asc,Desc; State:_,TamilNadu,Kerala,Karnataka,DefaultState"

package x

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

func _p() {
	var r River

	// should not report missing DefaultRiver, since it has same value as Ganga
	switch r {
	case Ganga, Yamuna, Kaveri:
	}

	// should not report missing Ganga, since it has same value as DefaultRiver
	switch r {
	case DefaultRiver, Yamuna, Kaveri:
	}
}

func _q() {
	var s State

	// value-based checks not available for iota enums (implementation detail: since
	// we cannot determine a constant.Value from type information).

	switch s { // want "missing cases in switch of type State: DefaultState"
	case TamilNadu, Kerala, Karnataka:
	}
}

func _r() {
	// should report correctly (in union '|' form) when same-valued names are
	// missing.
	var r River
	switch r { // want "missing cases in switch of type River: DefaultRiver|Ganga, Kaveri"
	case Yamuna:
	}

	// reporting should work correctly when constant.Values are not present also.
	var s State
	switch s { // want "missing cases in switch of type State: DefaultState, Kerala, TamilNadu"
	case Karnataka:
	}
}
