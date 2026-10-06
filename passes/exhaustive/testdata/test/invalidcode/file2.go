package invalidcode

import "test/invalidcode/sub"

type a = sub.S
type b = *sub.S
type c = *b

func (a) _() // invalid: cannot define new methods on non-local type sub.S
func (b) _() // invalid: cannot define new methods on non-local type sub.S
func (c) _() // invalid: invalid receiver type c
