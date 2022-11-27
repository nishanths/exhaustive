package enforcecomment

func callMe(a string, x map[Direction]int) map[Direction]int { return x }
func makeErr(a string, x map[Direction]int) error            { return nil }

var _ = map[Direction]int{
	N: 1,
}

//exhaustive:enforce
var _ = map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
	N: 1,
}

//exhaustive:enforce
var (
	_ = map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}
	_ = &map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}
	_ = map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}[N]
	_ = callMe("something", map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	})
)

var (
	//exhaustive:enforce
	_ = map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}
	_ = &map[Direction]int{
		N: 1,
	}
	_ = callMe("something", map[Direction]int{
		N: 1,
	})
)

func returnMap() map[Direction]int {
	switch 0 {
	case 1:
		// some other comment
		//exhaustive:enforce
		// some other comment
		return map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		}

	case 2:
		//exhaustive:enforce ... more arbitrary comment content (e.g. an explanation) ...
		return map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		}

	case 3:
		return map[Direction]int{
			N: 1,
		}

	case 4:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return map[Direction]int{ //exhaustive:enforce
			N: 1,
		}

	case 5:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return map[Direction]int{
			//exhaustive:enforce
			N: 1,
		}
	}
	return nil
}

func returnValueFromMap(d Direction) int {
	switch 0 {
	case 1:
		// some other comment
		//exhaustive:enforce
		// some other comment
		return map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		}[d]

	case 2:
		//exhaustive:enforce ... more arbitrary comment content (e.g. an explanation) ...
		return map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		}[d]

	case 3:
		return map[Direction]int{
			N: 1,
		}[d]

	case 4:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return map[Direction]int{ //exhaustive:enforce
			N: 1,
		}[d]

	case 5:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return map[Direction]int{
			//exhaustive:enforce
			N: 1,
		}[d]
	}
	return 0
}

func returnFuncCallWithMap() error {
	switch 0 {
	case 1:
		// some other comment
		//exhaustive:enforce
		// some other comment
		return makeErr("something", map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		}).(error)

	case 2:
		//exhaustive:enforce ... more arbitrary comment content (e.g. an explanation) ...
		return makeErr("something", map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		}).(error)

	case 3:
		return makeErr("something", map[Direction]int{
			N: 1,
		}).(error)

	case 4:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return makeErr("something", map[Direction]int{ //exhaustive:enforce
			N: 1,
		}).(error)

	case 5:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return makeErr("something", map[Direction]int{
			//exhaustive:enforce
			N: 1,
		}).(error)
	}
	return nil
}

func returnPointerToMap() *map[Direction]int {
	switch 0 {
	case 1:
		// some other comment
		//exhaustive:enforce
		// some other comment
		return &map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		}

	case 2:
		//exhaustive:enforce ... more arbitrary comment content (e.g. an explanation) ...
		return &map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		}

	case 3:
		return &map[Direction]int{
			N: 1,
		}

	case 4:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return &map[Direction]int{ //exhaustive:enforce
			N: 1,
		}

	case 5:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return &map[Direction]int{
			//exhaustive:enforce
			N: 1,
		}
	}
	return nil
}

func assignMapLiteral() {
	// some other comment
	//exhaustive:enforce
	// some other comment
	_ = map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}

	//exhaustive:enforce ... more arbitrary comment content (e.g. an explanation) ...
	_ = map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}

	//exhaustive:enforce
	a := map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}

	//exhaustive:enforce
	b, ok := map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}, 10

	_, _, _ = a, b, ok

	_ = map[Direction]int{
		N: 1,
	}

	// this should report: according to go/ast, the comment is not considered to
	// be associated with the assign node.
	_ = map[Direction]int{ //exhaustive:enforce
		N: 1,
	}

	// this should report: according to go/ast, the comment is not considered to
	// be associated with the assign node.
	_ = map[Direction]int{
		//exhaustive:enforce
		N: 1,
	}
}

func assignValueFromMapLiteral(d Direction) {
	// some other comment
	//exhaustive:enforce
	// some other comment
	_ = map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}[d]

	//exhaustive:enforce ... more arbitrary comment content (e.g. an explanation) ...
	_ = map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}[d]

	//exhaustive:enforce
	a := map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}[N]

	//exhaustive:enforce
	b, ok := map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}[N]

	_, _, _ = a, b, ok

	// this should report.
	_ = map[Direction]int{
		N: 1,
	}[d]

	// this should report: according to go/ast, the comment is not considered to
	// be associated with the assign node.
	_ = map[Direction]int{ //exhaustive:enforce
		N: 1,
	}[d]

	// this should report: according to go/ast, the comment is not considered to
	// be associated with the assign node.
	_ = map[Direction]int{
		//exhaustive:enforce
		N: 1,
	}[d]
}

func localVarDeclaration() {
	var _ = map[Direction]int{
		N: 1,
	}

	//exhaustive:enforce
	var _ = map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
		N: 1,
	}

	//exhaustive:enforce
	var (
		_ = map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		}
		_ = &map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		}
		_ = map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		}[N]
		_ = callMe("something", map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		})
	)

	var (
		//exhaustive:enforce
		_ = map[Direction]int{ // want "^missing keys in map of key type enforcecomment.Direction: enforcecomment.E, enforcecomment.S, enforcecomment.W, enforcecomment.directionInvalid$"
			N: 1,
		}
		_ = &map[Direction]int{
			N: 1,
		}
		_ = callMe("something", map[Direction]int{
			N: 1,
		})
	)
}
