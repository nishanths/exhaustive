package x

import barpkg "map/y"

type Stationery int

const (
	Pencil Stationery = iota
	Pen
	Paper
	Eraser
	sharpener
)

// Basic test of same package enum.

var j = map[Stationery]struct{}{ // want "missing keys in map j of key type Stationery: Paper, sharpener"
	Pencil: {},
	Pen:    {},
	Eraser: {},
}

// Basic test of external package enum.

var k = map[barpkg.Grain]struct{}{} // want "missing keys in map k of key type bar.Grain: Rice, Wheat"

// Parenthesized keys.

var l = map[Stationery]struct{}{ // want "missing keys in map l of key type Stationery: Paper, sharpener"
	(Pencil): {},
	(Pen):    {},
	(Eraser): {},
}

// Multiple names/values in ValueSpec.

var m, n = map[Stationery]struct{}{ // want "missing keys in map m of key type Stationery: Paper, sharpener" "missing keys in map n of key type bar.Grain: Wheat"
	Pencil: {},
	Pen:    {},
	Eraser: {},
}, map[barpkg.Grain]struct{}{
	barpkg.Rice: {},
}

// Should not error/panic when there is no value in ValueSpec.

var noValue map[Stationery]struct{}
