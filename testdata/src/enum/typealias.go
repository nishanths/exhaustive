package enum

import (
	"enum/otherpkg"
)

// type T2 int

// type T1 = T2

// const (
// 	T1_A T1 = iota
// 	T1_B
// )

type T3 = otherpkg.T4
type T5 = otherpkg.T5
