package general

import "errors"
import "test/check/typeswitch/general/sub"
import "test/check/typeswitch/dpkg"

type t struct{}

func (t) m()
func (*t) n()
func (t) o()

type u struct{}

func (*u) m()
func (u) n()

func f1() {
	// fundamentals

	var x interface{ m() }

	switch x.(type) {
	case t:
	case *t:
	case *u:
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	}

	switch x.(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*u$"
	case *t:
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*t$"
	case *u:
	}

	switch x.(type) { // want "^missing cases in type switch: t$"
	case *t:
	case *u:
	}

	switch x.(type) { // want "^missing cases in type switch: \\*t$"
	case t:
	case *u:
	}

	switch x.(type) { // want "^missing cases in type switch: \\*u$"
	case t:
	case *t:
	}
}

func f2() {
	// named interface

	type i interface{ m() }
	var x i

	switch x.(type) {
	case t:
	case *t:
	case *u:
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	}

	switch x.(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
	}
}

func f3() {
	var y interface {
		m()
		n()
	}
	var z interface {
		m()
		n()
		o()
	}

	switch y.(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	}

	switch z.(type) { // want "^missing cases in type switch: \\*t$"
	}
}

type v struct{}
func (*v) O() {}
type av = v

func f4() {
	var x interface{ O() }

	switch x.(type) {
	case *v:
	case dpkg.T:
	case *dpkg.T:
	}

	switch x.(type) {
	case *av:
	case dpkg.T:
	case dpkg.A:
	}

	switch x.(type) { // want "^missing cases in type switch: \\*v or \\*av, dpkg.T, \\*dpkg.T or dpkg.A$"
	}

	switch x.(type) { // want "^missing cases in type switch: \\*v or \\*av, \\*dpkg.T or dpkg.A$"
	case dpkg.T:
	}

	switch x.(type) { // want "^missing cases in type switch: \\*v or \\*av, dpkg.T$"
	case dpkg.A:
	}

	switch x.(type) { // want "^missing cases in type switch: \\*dpkg.T or dpkg.A$"
	case *v:
	case dpkg.T:
	}
}

func f5[P any]() {
	// various extra case clauses that do not contribute towards
	// making the type switch exhaustive.

	type e interface{ e() }
	type j interface {
		m()
		n()
	}
	type i interface {
		m()
	}

	var x i

	switch x.(type) { // want "^missing cases in type switch: t, \\*u$"
	case P:
	case error:
	case any:
	case *t:
	case e:
	case j:
	case interface{ e() }:
	case nil:
	}
}

type um interface{ unexportedMethod() }
var _ um = (*k)(nil) // sanity check that the test case is valid
type k struct{}
func (*k) unexportedMethod()

func f6() {
	var x um
	var y sub.UM

	switch x.(type) { // want "^missing cases in type switch: \\*k$"
	}

	switch x.(type) {
	case *k:
	}

	switch y.(type) { // want "^missing cases in type switch: \\*sub.K$"
	}

	switch y.(type) {
	case *sub.K:
	}
}

func f7() {
	// ast variations at type assertion node

	type i interface{ m() }
	var x interface{ m() }
	fi := func() i { return i(nil) }

	switch fi().(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
	}

	switch y := x.(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
		_ = y
	}

	switch n := 1; y := x.(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
		_ = y
		_ = n
	}

	switch y := fi().(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
		_ = y
	}

	switch y := (interface{ m() })((fi())).(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
		_ = y
	}

	s := struct{ o i }{}
	switch s.o.(i).(type) { // want "^missing cases in type switch: \\*t, \\*u$"
	case t:
	}
}

func e1() {
	// the empty interface

	type empty1 interface{}
	type empty2 empty1

	var x interface{}
	var y any
	var z empty1
	var w empty2

	switch x.(type) {
	}
	switch y.(type) {
	}
	switch z.(type) {
	}
	switch w.(type) {
	}
}

func e2() {
	// predeclared error interface

	type myerror error
	type aliaserror = error

	var err1 error
	err2 := errors.New("")
	var err3 myerror
	var err4 myerror

	switch err1.(type) {
	}
	switch err2.(type) {
	}
	switch err3.(type) { // want "^missing cases in type switch: .+, strconv.Error, \\*strconv.Error, .+$"
	}
	switch err4.(type) { // want "^missing cases in type switch: .+, strconv.Error, \\*strconv.Error, .+$"
	}
}
