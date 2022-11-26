package typeparam

import (
	"fmt"
	barpkg "general/y"
)

type M int // want M:"^A,B$"

const (
	A M = iota
	B
)

func (M) String() string { return "M-value" }

type N byte // want N:"^C,D$"

const (
	C N = iota
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
}

type L interface {
	M | NotEnumType
	fmt.Stringer
}

func bar0[T barpkg.Phylum | I](v T) {
	switch v {
	case T(A):
	case T(D):
	}
}

func bar1[T barpkg.Phylum | I](v T) {
	switch v {
	case T(A):
	case T(barpkg.Echinodermata):
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

func repeat[T M | K | I | O](v T) {
	switch v {
	}
}
