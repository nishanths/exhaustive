// foo

package a // foo

func f0() {
	switch {
	}
}

// foo

func f1() {
	// foo

	// a0
	switch {
	}

	// a1
	switch { // foo
	// foo
	}

	// foo

	// a2
	//
	// a2+continued
	//
	switch {
	}

	// a3
	switch {
	case 0:
		// foo

		// a4
		switch {
		} // foo

		// foo
	}

	switch {
	case 0:
		// a5
		switch {
		}
	}

	// foo
	for {
		// foo
	}

	// foo

	switch {
	}

	// a6
	switch {
	}

	// foo
}

func f2() {
	// a7
	switch {
	}
}

func f3() {
	/* a8 */ switch {
	}

	// a9
	// a9+continued
	/* a9.1
	 * a9.1+continued */switch {
	}

	// a10
	/* a10.1 */switch {
	}

	_ = 0 // foo
	switch {
	}
}

// foo
func f4() { // foo
	switch {
	}
	// foo
} // foo

// foo
func f5() {
	switch {
	}
} // foo

func f6() {
	// comment a11 should be associated only with the first
	// of these switches.

	// a11
	switch {}; switch {}
}

// foo
