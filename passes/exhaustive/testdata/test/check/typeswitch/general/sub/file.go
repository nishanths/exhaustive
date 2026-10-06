package sub

type UM interface{ unexportedMethod() }
var _ UM = (*K)(nil) // sanity check that the test case is valid
type K struct{}
func (*K) unexportedMethod()
