package typealias

import (
	"enum/typealias/otherpkg"
)

type (
	T1 = int // RHS is not allowed enum (alias -> valid basic type)
	T2 = T3  // RHS is not allowed enum (alias -> alias -> valid basic type)
	T9 = T8  // RHS is not allowed enum (alias -> alias -> ... -> alias -> valid basic type)
)

type (
	T4  = T5  // RHS is possible enum (alias -> named type -> valid basic type)
	T10 = T11 // RHS is possible enum (alias -> named type -> valid basic type)
	T6  = T7  // RHS is possible enum (alias -> alias -> ... -> alias -> named type -> valid basic type)
	T15 = T16 // RHS is possible enum (alias -> alias -> ... -> alias -> named type -> valid basic type)
)

// same comments as block above.
type (
	T12 = otherpkg.T5
	T13 = otherpkg.T11
	T14 = otherpkg.T7
	T17 = otherpkg.T16
)

// through multiple packages.
// RHS is possible enum.
type T18 = otherpkg.T19

// Disallowed RHS type.
type D2 struct{}
type D1 = D2

// -- Auxiliary types --

type T3 = int
type T8 = T3

type T5 int // NOTE: T5 has no members
type T7 = T5

type T11 int // want T11:"^T11_A,T11_B$"

const (
	T11_A T11 = iota
	T11_B
)

type T16 = T11
