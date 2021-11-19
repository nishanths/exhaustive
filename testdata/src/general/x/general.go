package x

import (
	"crypto/elliptic"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	bar "general/y"
	barpkg "general/y"
	"io/fs"
	"net/http"
	"os"
	"reflect"
)

func useComplexPackages() {
	// see issue #25: https://github.com/nishanths/exhaustive/issues/25
	var (
		_ http.Server
		_ tls.Conn
		_ reflect.ChanDir
		_ json.Encoder
		_ elliptic.Curve
	)
	fmt.Println(os.LookupEnv(""))
}

type Direction int // want Direction:"^N,E,S,W,directionInvalid$"

const (
	N                Direction = 1
	E                Direction = 2
	S                Direction = 3
	W                Direction = 4
	directionInvalid Direction = 5
)

func _a() {
	// Basic test of same package enum.
	//
	// Additionally: unexported members should be included in exhaustiveness
	// check since enum is in same package.

	var d Direction
	switch d { // want "^missing cases in switch of type Direction: E, directionInvalid$"
	case N:
	case S:
	case W:
	default:
	}
}

func _b() {
	// Basic test of external package enum.
	//
	// Additionally: unexported members should not be included in exhaustiveness
	// check since enum is in external package.

	var p bar.Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: Mollusca$"
	case bar.Chordata:
	case bar.Echinodermata:
	}
}

func _j() {
	// Named imports still report real package name.

	var p barpkg.Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: Mollusca$"
	case barpkg.Chordata:
	case barpkg.Echinodermata:
	}
}

func _f() {
	// Multiple values in single case.

	var d Direction
	switch d { // want "^missing cases in switch of type Direction: W$"
	case E, directionInvalid, S:
	default:
	case N:
	}
}

func _g() {
	// Switch isn't at top-level of function -- should still be checked.

	var d Direction
	if true {
		switch d { // want "^missing cases in switch of type Direction: S, directionInvalid$"
		case (N):
		case (E):
		case (W):
		}
	}

	switch d { // want "^missing cases in switch of type Direction: E, S, W, directionInvalid$"
	case N:
		switch d { // want "^missing cases in switch of type Direction: N, S, W$"
		case E, directionInvalid:
		}
	}
}

type SortDirection int // want SortDirection:"^Asc,Desc$"

const (
	_ SortDirection = iota // blank identifier need not be listed in switch statement to satisfy exhaustiveness
	Asc
	Desc
)

func _n() {
	var d SortDirection
	switch d {
	case Asc:
	case Desc:
	}
}

func _o() {
	// Selector isn't of the form "enumPkg.enumMember"

	type holdsPhylum struct {
		Mollusca bar.Phylum // can technically hold any Phylum value, but field is named Mollusca
	}

	var p bar.Phylum
	var h holdsPhylum

	switch p { // want "^missing cases in switch of type bar.Phylum: Mollusca$"
	case bar.Chordata:
	case bar.Echinodermata:
	case h.Mollusca:
	}
}

var ErrFoo = errors.New("foo")

func _p() {
	// Switch tag variable's type has nil package (lives in Universe scope).
	// Expect things to not panic and to not fail unexpectedly.

	var err error

	switch err {
	case nil:
	case ErrFoo:
	}
}

func _q() {
	// Type alias:
	// type os.FileMode = fs.FileMode
	//
	// Of interest, note that e.g. listing os.ModeSocket in a case clause is
	// equivalent to listing fs.ModeSocket (both have the same constant value).

	fi, err := os.Lstat(".")
	fmt.Println(err)

	switch fi.Mode() { // want "^missing cases in switch of type fs.FileMode: ModeDevice, ModePerm, ModeSetgid, ModeSetuid, ModeType$"
	case os.ModeDir:
	case os.ModeAppend:
	case os.ModeExclusive:
	case fs.ModeTemporary:
	case fs.ModeSymlink:
	case fs.ModeNamedPipe, os.ModeSocket:
	case fs.ModeCharDevice:
	case fs.ModeSticky:
	case fs.ModeIrregular:
	}
}

func _r(d Direction) {
	// Raw constants (i.e. not identifier or sel expr)
	// in case clauses do not count.

	switch d { // want "^missing cases in switch of type Direction: S, directionInvalid$"
	case N:
	case E:
	case 3:
	case W:
	case 5:
	}
}

func _s(u bar.Uppercase) {
	switch u {
	case bar.ReallyExported:
	}
}
