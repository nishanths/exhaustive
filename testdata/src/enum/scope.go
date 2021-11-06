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
