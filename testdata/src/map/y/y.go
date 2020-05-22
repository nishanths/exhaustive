// want package:"Grain:Wheat,Rice,corn"

package bar

type Grain int

const (
	Wheat = Grain(iota)
	Rice
	corn
)
