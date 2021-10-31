package fix

import (
	"fmt"
)

func _g() {
	var d Direction
	switch d { // want "^missing cases in switch of type Direction: E, directionInvalid$"
	case N:
	case S:
	case W:
	}
}

var _ = fmt.Printf
