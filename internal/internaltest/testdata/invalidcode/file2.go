package invalidcode

import "example/invalidcode/sub"

type a0 = sub.S
type a1 = *sub.S
type a2 = *a1

func (a0) _() // error: cannot define new methods on non-local type sub.S
func (a1) _() // error: cannot define new methods on non-local type sub.S
func (a2) _() // error: invalid receiver type a2

type t0 struct{}
type a3 = t0

func (a3) f()

var _ interface{ f() } = t0{} // note: t0 implements this interface because a3 (alias for t0) has method f()
var _ interface{ f() } = a3{}
