package foo

import "typealias/bar"

type T1 = bar.T2

const (
    A = bar.A
    B = bar.B

    C bar.T2 = "+"   // matches bar.C
    D bar.T2 = "***" // does not match bar.C
    F T1     = "|"   // matches bar.F by name and value (shows that type does not matter)
)

func ReturnsT1() T1 { return A }
