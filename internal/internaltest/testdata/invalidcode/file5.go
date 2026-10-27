package invalidcode

import "example/invalidcode/sub"

type n interface{ n() }

func f8() {
	var x n

	switch x.(type) {
	case sub.V:
		// error: impossible type switch case: sub.V
		// x (variable of type interface{n()}) cannot have dynamic type sub.V (unexported method n)
	}
}

type v struct{}
func (v) n() {}

var _ n = (*v)(nil)
var _ sub.N = (*v)(nil) // error: cannot use (*v)(nil) (value of type *v) as sub.N value in variable declaration: *v does not implement sub.N (unexported method n)

func f9() {
	// Go spec since go1.27: "Interface methods cannot
	// declare type parameters [...]."
	type _ interface {
		x[E any]() E
		// error: interface method must have no type parameters
		// error: undefined E
	}
	var _ interface {
		y[E any]() E
		// error: interface method must have no type parameters
		// error: undefined E
	}
}

func f10[P any]() {
	type _ interface {
		x() P // valid
	}
	var _ interface {
		y() P // valid
	}
}

