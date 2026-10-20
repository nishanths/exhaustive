package pkgd

type T struct{}
func (T) O() {}
type A = *T
