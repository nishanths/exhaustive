package otherpkg

import (
	"fmt"

	"github.com/nishanths/exhaustive/testdata/switchtest"
)

func switchParen(b switchtest.Biome) {
	switch b {
	case (switchtest.Tundra):
		fmt.Println("hi")
	case (switchtest.Desert):
	case (99):
	}
}

func switchNotSelExpr(b switchtest.Biome) {
	switch b {
	case 99:
	case switchtest.Tundra:
	}
}

type local struct {
	inner struct {
		f switchtest.Biome
	}
	f switchtest.Biome
}

func switchNotExpectedSelExpr(b switchtest.Biome, l local) {
	switch b {
	case l.inner.f:
	case switchtest.Desert:
	case l.f:
	}
}
