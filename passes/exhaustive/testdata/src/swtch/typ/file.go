package typ

import "swtch/typnew"

type A1 = typnew.S1

const (
	Za0 = typnew.Z0
	Za1 = typnew.Z1
	// Note: alias for Z2 absent.
	Za3 A1        = typnew.Z3
	Za4 typnew.S1 = typnew.Z4
	// Extra values:
	Za7  A1        = 7
	Za8  typnew.S1 = 8
	za9  A1        = 9
	za10 typnew.S1 = 10
)

// Note: The M1 type and related V* constants are a copy of the
// ones in testdata/*/enum/typ/file.go. Keep in sync.
type M1 int

const (
	V0, VV0 M1 = iota, iota //
	V1      M1 = iota       //
	V2, vv2 M1 = iota, iota //
	V3, vv3                 //
	V4      M1 = iota       //
	V5                      //
)
