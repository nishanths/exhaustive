package typeparam

import (
	"fmt"
	y "general/y"
)

type M int // want M:"^A,B$"

const (
	A M = iota * 100
	B
)

func (M) String() string { return "M-value" }

type N byte // want N:"^C,D$"

const (
	C N = iota * 100
	D
)

func (N) String() string { return "N-value" }

type O string // want O:"^E$"

const (
	E O = "hello"
)

func (O) String() string { return "O-value" }

type NotEnumType int

func (NotEnumType) String() string { return "NotEnumType-value" }

type I interface {
	N | J
}

type J interface {
	O
}

type K interface {
	M
	fmt.Stringer
	comparable
}

type L interface {
	M | NotEnumType
	fmt.Stringer
}

// "^missing cases in switch of type bar.Phylum|typeparam.N|typeparam.O: bar.Echinodermata, bar.Mollusca, typeparam.C, typeparam.E$"
// "^missing cases in switch of type bar.Phylum\\|typeparam.N|typeparam.O|typeparam.M: bar.Echinodermata, bar.Mollusca, typeparam.E$"

func bar0[T y.Phylum | I | M](v T) {
	switch v { // want `^missing cases in switch of type bar.Phylum\|typeparam.N\|typeparam.O\|typeparam.M: bar.Echinodermata, bar.Mollusca, typeparam.D\|typeparam.B, typeparam.E$`
	case T(A):
	}
}

/*
func bar1[T y.Phylum | I](v T) {
	switch v {
	case T(A):
	case T(y.Echinodermata):
	}
}

func foo0[T M](v T) {
	switch v {
	case T(B):
	}
}

func foo1[T K](v T) {
	switch v {
	case T(A):
	}
}

func fooNot0[T M | NotEnumType](v T) {
	switch v {
	case T(B):
	}
}

func fooNot1[T L](v T) {
	switch v {
	case T(B):
	}
}

func repeat[T I | O](v T) {
	switch v {
	}
}
*/
