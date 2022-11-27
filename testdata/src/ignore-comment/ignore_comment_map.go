package ignorecomment

func callMe(a string, x map[Direction]int) map[Direction]int { return x }
func makeErr(a string, x map[Direction]int) error            { return nil }

var _ = map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
	N: 1,
}

//exhaustive:ignore
var _ = map[Direction]int{
	N: 1,
}

//exhaustive:ignore
var (
	_ = map[Direction]int{
		N: 1,
	}
	_ = &map[Direction]int{
		N: 1,
	}
	_ = map[Direction]int{
		N: 1,
	}[N]
	_ = callMe("something", map[Direction]int{
		N: 1,
	})
)

var (
	//exhaustive:ignore
	_ = map[Direction]int{
		N: 1,
	}
	_ = &map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
		N: 1,
	}
	_ = callMe("something", map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
		N: 1,
	})
)

func returnMap() map[Direction]int {
	switch 0 {
	case 1:
		// some other comment
		//exhaustive:ignore
		// some other comment
		return map[Direction]int{
			N: 1,
		}

	case 2:
		//exhaustive:ignore ... more arbitrary comment content (e.g. an explanation) ...
		return map[Direction]int{
			N: 1,
		}

	case 3:
		return map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			N: 1,
		}

	case 4:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return map[Direction]int{ //exhaustive:ignore // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			N: 1,
		}

	case 5:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			//exhaustive:ignore
			N: 1,
		}
	}
	return nil
}

func returnValueFromMap(d Direction) int {
	switch 0 {
	case 1:
		// some other comment
		//exhaustive:ignore
		// some other comment
		return map[Direction]int{
			N: 1,
		}[d]

	case 2:
		//exhaustive:ignore ... more arbitrary comment content (e.g. an explanation) ...
		return map[Direction]int{
			N: 1,
		}[d]

	case 3:
		return map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			N: 1,
		}[d]

	case 4:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return map[Direction]int{ //exhaustive:ignore // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			N: 1,
		}[d]

	case 5:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			//exhaustive:ignore
			N: 1,
		}[d]
	}
	return 0
}

func returnFuncCallWithMap() error {
	switch 0 {
	case 1:
		// some other comment
		//exhaustive:ignore
		// some other comment
		return makeErr("something", map[Direction]int{
			N: 1,
		}).(error)

	case 2:
		//exhaustive:ignore ... more arbitrary comment content (e.g. an explanation) ...
		return makeErr("something", map[Direction]int{
			N: 1,
		}).(error)

	case 3:
		return makeErr("something", map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			N: 1,
		}).(error)

	case 4:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return makeErr("something", map[Direction]int{ //exhaustive:ignore // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			N: 1,
		}).(error)

	case 5:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return makeErr("something", map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			//exhaustive:ignore
			N: 1,
		}).(error)
	}
	return nil
}

func returnPointerToMap() *map[Direction]int {
	switch 0 {
	case 1:
		// some other comment
		//exhaustive:ignore
		// some other comment
		return &map[Direction]int{
			N: 1,
		}

	case 2:
		//exhaustive:ignore ... more arbitrary comment content (e.g. an explanation) ...
		return &map[Direction]int{
			N: 1,
		}

	case 3:
		return &map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			N: 1,
		}

	case 4:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return &map[Direction]int{ //exhaustive:ignore // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			N: 1,
		}

	case 5:
		// this should report: according to go/ast, the comment is not considered to
		// be associated with the return node.
		return &map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			//exhaustive:ignore
			N: 1,
		}
	}
	return nil
}

func assignMapLiteral() {
	// some other comment
	//exhaustive:ignore
	// some other comment
	_ = map[Direction]int{
		N: 1,
	}

	//exhaustive:ignore ... more arbitrary comment content (e.g. an explanation) ...
	_ = map[Direction]int{
		N: 1,
	}

	//exhaustive:ignore
	a := map[Direction]int{
		N: 1,
	}

	//exhaustive:ignore
	b, ok := map[Direction]int{
		N: 1,
	}, 10

	_, _, _ = a, b, ok

	_ = map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
		N: 1,
	}

	// this should report: according to go/ast, the comment is not considered to
	// be associated with the assign node.
	_ = map[Direction]int{ //exhaustive:ignore // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
		N: 1,
	}

	// this should report: according to go/ast, the comment is not considered to
	// be associated with the assign node.
	_ = map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
		//exhaustive:ignore
		N: 1,
	}
}

func assignValueFromMapLiteral(d Direction) {
	// some other comment
	//exhaustive:ignore
	// some other comment
	_ = map[Direction]int{
		N: 1,
	}[d]

	//exhaustive:ignore ... more arbitrary comment content (e.g. an explanation) ...
	_ = map[Direction]int{
		N: 1,
	}[d]

	//exhaustive:ignore
	a := map[Direction]int{
		N: 1,
	}[N]

	//exhaustive:ignore
	b, ok := map[Direction]int{
		N: 1,
	}[N]

	_, _, _ = a, b, ok

	// this should report.
	_ = map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
		N: 1,
	}[d]

	// this should report: according to go/ast, the comment is not considered to
	// be associated with the assign node.
	_ = map[Direction]int{ //exhaustive:ignore // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
		N: 1,
	}[d]

	// this should report: according to go/ast, the comment is not considered to
	// be associated with the assign node.
	_ = map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
		//exhaustive:ignore
		N: 1,
	}[d]
}

func localVarDeclaration() {
	var _ = map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
		N: 1,
	}

	//exhaustive:ignore
	var _ = map[Direction]int{
		N: 1,
	}

	//exhaustive:ignore
	var (
		_ = map[Direction]int{
			N: 1,
		}
		_ = &map[Direction]int{
			N: 1,
		}
		_ = map[Direction]int{
			N: 1,
		}[N]
		_ = callMe("something", map[Direction]int{
			N: 1,
		})
	)

	var (
		//exhaustive:ignore
		_ = map[Direction]int{
			N: 1,
		}
		_ = &map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			N: 1,
		}
		_ = callMe("something", map[Direction]int{ // want "^missing keys in map of key type ignorecomment.Direction: ignorecomment.E, ignorecomment.S, ignorecomment.W, ignorecomment.directionInvalid$"
			N: 1,
		})
	)
}
