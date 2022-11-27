//go:build go1.18
// +build go1.18

package typeparam

import (
	y "general/y"
)

type M uint8 // want M:"^A,B$"
const (
	_ M = iota * 100
	A
	B
)

func (M) String() string { return "" }

type N uint8 // want N:"^C,D$"
const (
	_ N = iota * 100
	C
	D
)

type O byte // want O:"^E1,E2$"
const (
	E1 O = 'h'
	E2 O = 'e'
)

type P float32 // want P:"^F$"
const (
	F P = 1.1234
)

type Q string // want Q:"^G$"
const (
	G Q = "world"
)

type NotEnumType uint8

type Stringer interface {
	String() string
}

type II interface{ N | JJ }
type JJ interface{ O }
type KK interface {
	M
	Stringer
	error
	comparable
}
type LL interface {
	M | NotEnumType
	Stringer
	error
}
type MM interface {
	M
}
type QQ interface {
	Q
}

func _a[T y.Phylum | M](v T) {
	switch v { // want `^missing cases in switch of type bar.Phylum\|typeparam.M: bar.Chordata, bar.Mollusca, typeparam.B$`
	case T(A):
	case T(y.Echinodermata):
	}

	switch M(v) { // want "^missing cases in switch of type typeparam.M: typeparam.A, typeparam.B$"
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type bar.Phylum\|typeparam.M: bar.Chordata, bar.Mollusca, typeparam.B$`
		T(A):               struct{}{},
		T(y.Echinodermata): struct{}{},
	}
}

func _b[T N | MM](v T) {
	switch v { // want `^missing cases in switch of type typeparam.N\|typeparam.M: typeparam.D\|typeparam.B$`
	case T(A):
	}

	switch M(v) { // want "^missing cases in switch of type typeparam.M: typeparam.A, typeparam.B$"
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.N\|typeparam.M: typeparam.D\|typeparam.B$`
		T(A): struct{}{},
	}
}

func _c[T O | M | N](v T) {
	switch v { // want `^missing cases in switch of type typeparam.O\|typeparam.M\|typeparam.N: typeparam.E1, typeparam.E2, typeparam.A\|typeparam.C$`
	case T(B):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.O\|typeparam.M\|typeparam.N: typeparam.E1, typeparam.E2, typeparam.A\|typeparam.C$`
		T(B): struct{}{},
	}
}

func _d[T y.Phylum | II | M](v T) {
	switch v { // want `^missing cases in switch of type bar.Phylum\|typeparam.N\|typeparam.O\|typeparam.M: bar.Chordata, bar.Mollusca, typeparam.D\|typeparam.B, typeparam.E1, typeparam.E2$`
	case T(A):
	case T(y.Echinodermata):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type bar.Phylum\|typeparam.N\|typeparam.O\|typeparam.M: bar.Chordata, bar.Mollusca, typeparam.D\|typeparam.B, typeparam.E1, typeparam.E2$`
		T(A):               struct{}{},
		T(y.Echinodermata): struct{}{},
	}
}

func _e[T M](v T) {
	switch v { // want `^missing cases in switch of type typeparam.M: typeparam.A$`
	case T(M(B)):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.M: typeparam.A$`
		T(M(B)): struct{}{},
	}
}

func repeat0[T II | O](v T) {
	switch v { // want `^missing cases in switch of type typeparam.N\|typeparam.O: typeparam.C, typeparam.D, typeparam.E2$`
	case T(E1):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.N\|typeparam.O: typeparam.C, typeparam.D, typeparam.E2$`
		T(E1): struct{}{},
	}
}

func repeat1[T MM | M](v T) {
	switch v { // want `^missing cases in switch of type typeparam.M: typeparam.A$`
	case T(B):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.M: typeparam.A$`
		T(B): struct{}{},
	}
}

func _mixedTypes0[T M | QQ](v T) {
	// expect no diagnostic because underlying basic kinds are not same:
	// uint8 vs. string
	switch v {
	case T(A):
	}

	_ = map[T]struct{}{
		T(A): struct{}{},
	}
}

func _mixedTypes1[T MM | QQ](v T) {
	switch v {
	case T(A):
	}

	_ = map[T]struct{}{
		T(A): struct{}{},
	}
}

func _notEnumType0[T M | NotEnumType](v T) {
	// expect no diagnostic because not type elements are enum types.
	switch v {
	case T(B):
	}

	_ = map[T]struct{}{
		T(B): struct{}{},
	}
}

func _notEnumType1[T LL](v T) {
	switch v {
	case T(A):
	}

	_ = map[T]struct{}{
		T(A): struct{}{},
	}
}
