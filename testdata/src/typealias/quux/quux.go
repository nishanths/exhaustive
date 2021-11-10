package quux

import (
    "typealias/bar"
    "typealias/foo"
)

func x() {
    var v foo.T1 = foo.ReturnsT1()

    switch v { // want "^missing cases in switch of type bar.T2: D, E$"
    case foo.A:
    case bar.B:
    case foo.C:
    case foo.D:
    case foo.F:
    case foo.H:
    }

    var w bar.T2 = foo.ReturnsT1()
    switch w { // want "^missing cases in switch of type bar.T2: D, E$"
    case foo.A:
    case bar.B:
    case foo.C:
    case foo.D:
    case foo.F:
    case foo.H:
    }
}
