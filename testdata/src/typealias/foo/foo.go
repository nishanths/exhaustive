package foo

import "typealias/bar"

type T1 = bar.T2

// None of these constants can constitue T2's enum members
// because they are not in the same package as the enum type T2.
const (
    A        = bar.A // matches bar.A by value; can be listed in switch case instead of bar.A
    B        = bar.B // matches bar.B by value; can be listed in switch case instead of bar.B
    C bar.T2 = "+"   // matches bar.C by value; can be listed instead of bar.C in switch
    F T1     = "|"   // matches bar.F by value (type does not matter); can be listed in switch case instead of bar.F
    H bar.T2 = "<"   // matches bar.I by value (name does not matter); can be listed in switch case instead of bar.I

    D bar.T2 = "***" // does not match bar.D
    G bar.T2 = "@@@" // some arbitrary bar.T2 value
)

func ReturnsT1() T1 { return A }
