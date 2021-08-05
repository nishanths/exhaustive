// Package funccall is used in TestContainsFuncCall.
package funccall

type Int int

func (i Int) s(j int) int {
	return int(i) + j
}

func f(i int) int {
	return i
}

var i int
var integer Int

// Keep this at the end of file (test code assumes this).
var (
	_ = f(i)              // true
	_ = Int(i)            // false
	_ = bool(true)        // false
	_ = Int(f(i))         // true
	_ = f(Int(i))         // true
	_ = integer.s(i)      // true
	_ = Int(integer.s(i)) // true
	_ = f(integer.s(i))   // true
)
