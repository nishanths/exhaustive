package anotherpkg

type T1 rune // want T1:"^T1_A$"

const (
	T1_A T1 = iota
)
