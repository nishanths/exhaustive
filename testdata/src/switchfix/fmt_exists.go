package switchfix

import (
	"fmt"
)

func _fmt_exists() {
	var d Direction
	switch d { // want "missing cases in switch of type Direction: E, directionInvalid"
	case N:
	case S:
	case W:
	}
}

var _ = fmt.Printf
