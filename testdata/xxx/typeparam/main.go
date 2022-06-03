package typeparam

// Testing instructions:
//  $ go build ./cmd/exhaustive
//  $ ./exhaustive ./testdata/xxx/typeparam

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

func foo[T M](v T) {
	switch v {
	case T(A):
	}
}

func bar[T M | N](v T) {
	switch v {
	case T(A):
	case T(D):
	}
}
