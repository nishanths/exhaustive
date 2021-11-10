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

func _k(d Direction) {
	// Parenthesized values in case statements.

	switch d { // want "^missing cases in switch of type Direction: S, directionInvalid$"
	case (N):
	case (E):
	case (W):
	}

	// Parenthesized values in switch tag.
	switch d { // want "^missing cases in switch of type Direction: S, directionInvalid$"
	case N:
	case E:
	case W:
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
	// Both os.Mode* and fs.Mode* constants can exist in the case
	// clauses (the Go type system allows it).
	// When checking if exhaustiveness is satisfied, the exhaustive analyzer
	// will "match" same-named and same-valued constants in package os and
	// package fs, for the listed case clause expressions in the switch
	// statement.
	// This means, for example, that listing os.ModeSocket is equivalent to
	// listing fs.ModeSocket (since they have the same name and the same
	// constant value).
	//
	// This test case tests for the above described scenarios.

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
