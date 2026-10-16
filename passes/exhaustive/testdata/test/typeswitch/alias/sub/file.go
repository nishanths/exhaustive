package sub

// exported alias to unexported defined type.
type V = v
type VV = *v

type v struct{}
func (v) N() {}
