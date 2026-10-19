// foo

package a // foo

func f0() {
	_ = map[K]V{}
}

// foo

// b0
var (
	// b1
	_ = map[K]V{} // b1.1
	// b2
	_ = map[K]V{}
	// foo

	// foo

	_ = map[K]V{}
	// foo

	// b3
	_ = map[K]V{}

	_ = func() {
		f(map[K]V{})
	}
)

// foo

func f1() {
	// foo

	// a0
	_ = map[K]V{}

	// a1
	_ = map[K]V{ // foo
	}

	// foo

	// a2
	//
	// a2+continued
	//
	_ = map[K]V{}

	// a3
	_ = map[K]V{
		0: func() {
			// foo

			// a4
			_ = map[K]V{} // foo

			// foo
		},
	}

	_ = map[K]V{
		0: func() {
			// a5
			_ = map[K]V{}
		},
	}

	// foo
	for {
		// foo
	}

	// foo

	_ = map[K]V{}

	// a6
	_ = map[K]V{}

	// foo
}

func f2() {
	// a7
	_ = map[K]V{}
}

func f3() {
	/* a8 */ _ = map[K]V{}

	// a9
	// a9+continued
	/* a9.1
	 * a9.1+continued */_ = map[K]V{}

	// a10
	/* a10.1 */_ = map[K]V{}

	_ = 0 // foo
	_ = map[K]V{}
}

// foo
func f4() { // foo
	_ = map[K]V{}
	// foo
} // foo

// foo
func f5() {
	_ = map[K]V{}
} // foo

func f6() {
	// Comment a11 should be associated only with the first
	// of these composite literals.

	// a11
	_ = map[K]V{}; _ = map[K]V{}
}

// foo
