package enum

// import (
// 	"enum/otherpkg"
// )

// type T2 int

// type T1 = T2

// const (
// 	T1_A T1 = iota
// 	T1_B
// )

// type T3 = otherpkg.T4
// type T5 = otherpkg.T5

// type T6 = int

type T1 = int // not allowed (alias -> valid basic type)
type T2 = T3  // not allowed (alias -> alias -> valid basic type)
type T9 = T8  // not allowed (alias -> alias -> ... -> alias -> valid basic type)
type T4 = T5  // possible    (alias -> named type -> valid basic type)
type T6 = T7  // possible    (alias -> alias -> ... -> alias -> named type -> valid basic type)

type T3 = int
type T8 = T3
type T5 int // NOTE: does not matter right now that T5 has no known members
type T7 = T5

type X1 = X2

type X2 int

const (
	X1_A X1 = iota
	X1_B
)
