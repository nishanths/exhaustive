package app

import "github.com/nishanths/exhaustive/testdata/playground/env"

func readFile(path string) ([]byte, error) {
	switch env.Current() {
	case env.Production:
	case env.Dev:
	default:
	}
	panic("")
}
