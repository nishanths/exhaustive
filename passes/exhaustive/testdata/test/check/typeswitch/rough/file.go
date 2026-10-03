package rough

import (
	"go/types"
	"io"
)

// XXX: what about the following declarations:
// type T int
// type S = *T
// type U = *S // and variants of this (i.e. multi-level aliases)
//
// type T int
// type S = *T
// type U = S // and variants of this (i.e. multi-level aliases)
//
// type T int
// type S = T
// type U = *S // and variants of this (i.e. multi-level aliases)
//
// type T int
// type S = T
// type U = S // and variants of this (i.e. multi-level aliases)


type T int
type S *T
type Q *int
type R Q

func (T) _()
func (*T) _()

// illegal:
// XXX: need to make sure in the code
// that S is not considered as implemented by types.Implements.
// func (S) _()
// func (*S) _()

// illegal:
// func (Q) _()
// func (*Q) _()

// illegal:
// func (R) _()
// func (*R) _()

func f1() {
	var x any
	var y interface{}
	
	switch x.(type) {
	}
	
	switch y.(type) {
	}
	
	var e error
	switch e.(type) {
	}
	
	var g error
	switch g.(type) {
	// illegal: impossible type switch case: int
	// case int:
	}
	
	var n int
	_ = n
	// illegal: n (variable of type int) is not an interface
	// switch n.(type) {}
	
	/*
	type s struct {
		o io.Closer
	}
	var sval s
	switch sval.o.(io.ReadCloser).(type) {}
	*/
}

func f2[P any](x any) {
	switch x.(type) {
	case P:
	case nil:
	default:
	}
}

func f3[P any](p P) {
	// illegal:
	// switch p.(type) {
	// }
}

func f4() {
	var t types.Type
	switch t.(type) {
	case *types.Alias:
	case *types.Named:
	}
}

func f5() {
	type IX interface {
		comparable
	}
	
	// illegal:
	// var ix IX
	// switch ix.(type) {
	// }
}

func f6() {
	var a interface{}
	switch a.(type) {
	case int:
	case error:
	case any:
	case io.Closer:
	case interface { Close() error }:
	}
	
	/*
	var b io.Closer
	switch b.(type) {
	case io.ReadCloser:
	}
	*/
}

type U int

func (u U) Close() error { return nil }
func (u U) Read([]byte) (int, error) { return 0, nil }

func f7() {
	var _ io.Closer = U(0)
	var _ io.Reader = U(0)

	var _ io.Closer = (*U)(nil)
	var _ io.Reader = (*U)(nil)
}

type A = *B
type B int
func (*B) Close() error

func f8() {
	var a io.Closer = A(nil)
	_ = a
}

func f9() {
	type b int
	type a = b
	
	// illegal: duplicate case b in type switch
	// var x any
	// switch x.(type) {
	// case a:
	// case b:
	// }
}
