// This file contains tests for code snippets presented in source
// code comments in package exhaustive, etc.

package internaltest

import (
	"slices"
	"testing"
)

type t1 struct{}

func (t1) o() {}

type t2 struct{}

func (t2) o() {}

func f00[P interface{ o() }]() []string {
	var ret []string
	var x interface{ o() } = t2{}

	switch x.(type) {
	case P:
		ret = append(ret, "P")
	case t2:
		ret = append(ret, "t2")
	}

	switch x.(type) {
	case P:
		ret = append(ret, "P")
	case interface{ o() }:
		ret = append(ret, "o")
	}

	return ret
}

// similar to func f00, except that the type constraint has
// a union of type terms instead of an interface.
func f01[P t1 | t2]() []string {
	var ret []string
	var x interface{ o() } = t2{}

	switch x.(type) {
	case P:
		ret = append(ret, "P")
	case t2:
		ret = append(ret, "t2")
	}

	return ret
}

func TestT0(t *testing.T) {
	// comment in function in checkImpossibleCases
	// in file exhaustive/check.go.

	got := f00[t1]()
	want := []string{"t2", "o"}
	if !slices.Equal(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}

	got = f01[t1]()
	want = []string{"t2"}
	if !slices.Equal(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
