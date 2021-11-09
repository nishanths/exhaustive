package otherpkg

import "enum/otherpkg/anotherpkg"

type T5 int
type T7 = T5

type T11 int // want T11:"^T11_A,T11_B$"

const (
	T11_A T11 = iota
	T11_B
)

type T16 = T11

type T19 = anotherpkg.T1
