package typeparam

import barpkg "general/y"

type M int

const (
	A M = iota
	B
)

type N uint

const (
	C N = iota
	D
)

func foo[T barpkg.Phylum](v T) {
	switch v {
	case T(barpkg.Chordata):
	}
}

func bar[T M | N](v T) {
	switch v {
	case T(A):
	case T(D):
	}
}

func bazz[T ~int](v T) {
	switch v {
	case T(A):
	}
}

type X interface {
	M | N
}

type O string

func qux[T X](v T) {
	switch v {
	case T(A):
	}
}
