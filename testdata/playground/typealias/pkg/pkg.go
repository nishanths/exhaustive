package pkg

import (
	"goplay/pkg/subpkg"
	"goplay/thirdparty"
)

// Earlier this package was returning third party X type, but now this package
// wants to return my own X type instead of third party type.
//
// Old API: thirdparty.X (to be unreferenced in new API's package -- this package)
// New API: OwnX
//
// (no line of code) -->
type OwnX = thirdparty.X // -->
// type OwnX int

const (
	TMemI OwnX = thirdparty.TMemI
	TMemJ OwnX = thirdparty.TMemJ
	TMemK OwnX = 1002 // new extra
)

func ReturnsX() OwnX { return TMemI }

// is it okay to return here the extra?
const _ = TMemK // No, because users aren't yet expected to know of OwnX

// Earlier this package defined own Y type, but now it wants to move the
// definition to a dedicated (sub)package.
//
// Old API: OwnY (to be removed in old API's package -- this package)
// New API: subpkg.SubpkgY
//
// type OwnY uint -->
type OwnY = subpkg.SubpkgY // -->
// (no line of code)

const (
	// created these in subpkg, which is going to become the primary package
	// and forwarding to the created ones in subkpkg
	MemM = subpkg.MemM
	MemN = subpkg.MemN
	// at the end these will be become (no line of code)
)

func ReturnsY() subpkg.SubpkgY { return subpkg.MemM }

// is it okay to return here the extra?
const _ = subpkg.MemO // No, because users aren't yet expected to know of SubpkgY
