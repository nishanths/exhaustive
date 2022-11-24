package x

import barpkg "general/y"

const (
	PlainIntA = 1
	PlainIntB = 2
)

func _c() {
	// Tagless switch -- should be ignored.

	var p barpkg.Phylum
	switch {
	case p == barpkg.Chordata:
	case p == barpkg.Echinodermata:
	}
}

func _d() {
	// Tag value is of unnamed type -- should be ignored.

	var a int
	switch a {
	case PlainIntA:
	}

	_ = map[int]int{
		PlainIntA: 1,
	}
}

type NamedButNotEnum int

func _e() {
	// Tag value is a named type, but the named type isn't an enum -- should be
	// ignored.

	var a NamedButNotEnum
	switch a {
	case 1:
	}

	_ = map[NamedButNotEnum]int{
		1: 1,
	}
}

func emptyMapShouldBeIgnored() {
	_ = map[barpkg.Phylum]int{}
}
