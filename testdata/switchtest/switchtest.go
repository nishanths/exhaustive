package switchtest

import "log"

func switchWithDefault(b Biome) {
	switch b {
	case Tundra:
		log.Println("hi")
	case Desert:
		_ = 42
	default:
		panic("boom")
	}
}

func switchWithoutDefault(b Biome) {
	switch b {
	case Tundra:
		log.Println("hi")
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
