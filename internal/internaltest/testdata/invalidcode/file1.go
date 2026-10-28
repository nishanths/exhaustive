// Package invalidcode contains examples of valid and invalid Go
// code for reference purposes.
package invalidcode

import "cmp"

func f1() {
	var e error
	var x interface{}

	switch e.(type) {
	case int: // error: impossible type switch case: int
	}

	var n int
	switch n.(type) { // error: n (variable of type int) is not an interface
	}

	switch (&x).(type) { // error: (&x) (value of type *interface{}) is not an interface
	}
}

func f2() {
	type b struct{}
	type a = b

	var x any
	switch x.(type) {
	case a:
	case b: // error: duplicate case b in type switch
	}
}

func f3[P any](p P) {
	switch p.(type) { // error: cannot use type switch on type parameter value p (variable of type P constrained by any)
	}

}

func f4() {
	type i interface{ comparable }
	type j interface{ int }
	var z any

	switch z.(type) {
	case cmp.Ordered: // error: cannot use type cmp.Ordered outside a type constraint: interface contains type constraints
	case i: // error: cannot use type i outside a type constraint: interface is (or embeds) comparable
	case j: // error: cannot use type j outside a type constraint: interface contains type constraints
	}
}
