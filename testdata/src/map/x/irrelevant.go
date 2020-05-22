package x

func _a() {
	// Non-top level map -- should not be checked.

	var a = map[Stationery]int{
		Pen:   1,
		Paper: 0,
	}
	_ = a
}

type NamedButNotEnum int

// Key is a named type, but the named type isn't an enum -- should be
// ignored.

var c = map[NamedButNotEnum]int{
	1: 1,
	2: 0,
}

// Key is unnamed basic type -- should be ignored.

const PlainIntA = 1

var d = map[int]bool{
	PlainIntA: false,
}
