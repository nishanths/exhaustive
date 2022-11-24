// Package switchtest is used by switch_test.go.
package switchtest

import "fmt"

func switchWithDefault(b Biome) {
	switch b {
	case Tundra:
		fmt.Println("hi")
	case Desert:
		_ = 42
	default:
		panic("boom")
	}
}

func switchWithoutDefault(b Biome) {
	switch b {
	case Tundra:
		fmt.Println("hi")
	case Desert:
		_ = 42
	}
}

func switchParen(b Biome) {
	switch b {
	case ((Tundra)), (Desert):
	}
}

func switchNotIdent(b Biome) {
	switch b {
	case 1, 2:
	case 3:
	case Savanna:
	}
}
