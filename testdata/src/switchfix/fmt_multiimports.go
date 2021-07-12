package switchfix

import (
	"errors"
	"path"
)

var _ = errors.New
var _ = path.Base

func _fmt_multiimports() {
	var d Direction
	switch d { // want "^missing cases in switch of type Direction: E, directionInvalid$"
	case N:
	case S:
	case W:
	}
}
