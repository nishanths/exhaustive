package sub

type S struct{}

type V struct{}
func (V) n() {}

type n interface{ n() }
type N interface{ n() }

var _ n = (*V)(nil)
var _ N = (*V)(nil)
var _ interface{ n() } = (*V)(nil)
