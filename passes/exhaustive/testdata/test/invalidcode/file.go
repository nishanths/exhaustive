// Package invalid contains examples of invalid Go code for
// reference purposes.
package invalidcode

func f1() {
	var e error
	switch e.(type) {
	case int: // invalid: impossible type switch case: int
	}

	var n int
	switch n.(type) { // invalid: n (variable of type int) is not an interface
	}
}

func f2() {
	type b int
	type a = b

	var x any
	switch x.(type) {
	case a:
	case b: // invalid: duplicate case b in type switch
	}
}

func f3[P any](p P) {
	switch p.(type) { // invalid: cannot use type switch on type parameter value p (variable of type P constrained by any)
	}
}

func f4() {
	type i interface {
		comparable
	}

	var x i // invalid: cannot use type i outside a type constraint: interface is (or embeds) comparable
	_ = x
}

func f5() {
	var x interface{ n() }

	switch x.(type) {
	case sub.V:
		// invalid: impossible type switch case: sub.V
		// x (variable of type interface{n()}) cannot have dynamic type sub.V (unexported method n)
	}
}
