package typealias

type T1 = int // not allowed (alias -> valid basic type)
type T2 = T3  // not allowed (alias -> alias -> valid basic type)
type T9 = T8  // not allowed (alias -> alias -> ... -> alias -> valid basic type)
type T4 = T5  // possible    (alias -> named type -> valid basic type)
type T6 = T7  // possible    (alias -> alias -> ... -> alias -> named type -> valid basic type)

type T3 = int
type T8 = T3
type T5 int // NOTE: does not matter right now that T5 has no known members
type T7 = T5

// TODO(testing): The above should hold true if
// T3, T5, etc. are in different packages from T2, T4, etc. respectively.
