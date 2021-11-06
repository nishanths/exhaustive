package enum

type T uint // want T:"^A,B$"

const (
	A T = iota
	B
)

func F() {
	const (
		AA T = iota
		BB
	)

	f := func() {
		const (
			CC T = iota
			DD
		)

		type T uint // want T:"^C,D,E,F$"
		const (
			C T = iota
			D
		)
		const E T = 42

		for {
			const (
				EE T = iota
				FF
			)
			break
		}

		const F T = 43
	}
	_ = f
}

func F2() {
	if true {
		type T uint // want T:"^A,B$"
		const (
			A T = iota
			B
		)
	}
}

// Members must be in the same scope.

type PkgRequireSameLevel uint // want PkgRequireSameLevel:"^PA$"

const (
	_  PkgRequireSameLevel = 100
	PA PkgRequireSameLevel = 200
)

func f() {
	const (
		_  PkgRequireSameLevel = 42
		PC PkgRequireSameLevel = 0
		PD PkgRequireSameLevel = 1
	)
}

type PkgRequireSameLevel_2 uint

func g() {
	const PE PkgRequireSameLevel_2 = 9

	for {
		type InnerRequireSameLevel uint // want InnerRequireSameLevel:"^IX,IY$"

		const (
			_  InnerRequireSameLevel = 100
			IX InnerRequireSameLevel = 200
			IY InnerRequireSameLevel = 200
		)

		if true {
			const (
				_  InnerRequireSameLevel = 42
				IM InnerRequireSameLevel = 0
				IN InnerRequireSameLevel = 1
			)
		}
	}
}
