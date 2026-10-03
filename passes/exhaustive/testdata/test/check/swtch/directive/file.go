package directive

type t1 int

const (
	x0 t1 = iota
	x1
)

func f1() {
	//exhaustive:badname
	switch { // want "^error parsing comment directives: invalid name \"badname\"$"
	}

	//exhaustive:enforce
	//exhaustive:ignore
	switch { // want "^error parsing comment directives: conflicting directives$"
	}

	//exhaustive:enforce
	switch { // want "^enforce directive present and switch not checked: no switch expression$"
	}

	//exhaustive:enforce
	switch int(0) { // want "^enforce directive present and switch not checked: switch expression type is not an enumerated type$"
	}

	//exhaustive:ignore
	switch t1(0) {
	case x0:
	}

	//exhaustive:enforce
	switch t1(0) { // want "^switch not exhaustive: missing cases: x1$"
	case x0:
	}

	//exhaustive:ignore
	switch t1(0) {
	case x0:
		switch t1(0) { // want "^switch not exhaustive: missing cases: x0$"
		case x1:
		}
	}

	switch t1(0) { // want "^switch not exhaustive: missing cases: x1$"
	case x0:
		//exhaustive:ignore
		switch t1(0) {
		case x1:
		}
	}
}

func g1() {
	//exhaustive:badname
	_ = map[any]any{} // want "^error parsing comment directives: invalid name \"badname\"$"

	//exhaustive:enforce
	//exhaustive:ignore
	_ = map[any]any{} // want "^error parsing comment directives: conflicting directives$"

	//exhaustive:ignore
	var (
		//exhaustive:enforce
		_ = map[t1]bool{x0: true} // want "^error parsing comment directives: conflicting directives$"
		_ = map[t1]bool{x1: true}
	)

	//exhaustive:enforce
	var (
		//exhaustive:ignore
		_ = map[t1]bool{x0: true} // want "^error parsing comment directives: conflicting directives$"
		_ = map[t1]bool{x1: true} // want "^map literal not exhaustive: missing keys: x0$"
	)

	//exhaustive:enforce
	_ = map[int]bool{0: true} // want "^enforce directive present and map literal not checked: key type is not an enumerated type$"

	//exhaustive:ignore
	_ = map[t1]bool{
		x0: true,
	}

	//exhaustive:enforce
	_ = map[t1]bool{ // want "^map literal not exhaustive: missing keys: x1$"
		x0: true,
	}

	//exhaustive:ignore
	_ = map[t1]map[t1]bool{
		//exhaustive:enforce
		x0: {
			x1: true,
		},
	}

	//exhaustive:enforce
	_ = map[t1]map[t1]bool{ // want "^map literal not exhaustive: missing keys: x1$"
		//exhaustive:ignore
		x0: { // want "^map literal not exhaustive: missing keys: x0$"
			x1: true,
		},
	}
}
