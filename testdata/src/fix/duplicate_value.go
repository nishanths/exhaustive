package fix

type River string

const DefaultRiver = Ganga

const (
	Yamuna River = "Yamuna"
	Ganga  River = "Ganga"
	Kaveri River = "Kaveri"
)

func _l() {
	var r River

	// value "Ganga", shared by two enum members, is missing.
	switch r { // want "^missing cases in switch of type River: DefaultRiver|Ganga$"
	case Yamuna, Kaveri:
	}
}
