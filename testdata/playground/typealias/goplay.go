package main

import (
	"goplay/pkg"
	"goplay/pkg/subpkg"
	"goplay/someone"
	"goplay/thirdparty"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	mainV40()
}

func mainV40() {
	var a thirdparty.X // I don't know/shouldn't have to know at this point that pkg.OwnX exists/will exist
	var b pkg.OwnY     // I don't know/shouldn't have to know at this point that subpkg.SubpkgY exists/will exist

	a = pkg.ReturnsX()
	b = pkg.ReturnsY()

	switch a {
	case thirdparty.TMemI:
	case thirdparty.TMemJ:
	}

	switch b {
	case pkg.MemM:
	case pkg.MemN:
	}

	log.Println(a, b)
}

func mainV45() {
	var a pkg.OwnX       // done to latest
	var b subpkg.SubpkgY // done to latest

	a = pkg.ReturnsX()
	b = pkg.ReturnsY()

	switch a {
	case pkg.TMemI:
	case pkg.TMemJ:
		// case pkg.TMemK:
	}

	switch b {
	case subpkg.MemM:
	case subpkg.MemN:
		// case subpkg.MemO:
	}

	log.Println(a, b)
}

func anotherV41() {
	var c pkg.OwnY // I don't know/shouldn't have to know at this point that someone.SomeonesOwnY exists/will exist
	c = someone.P()
	switch c {
	case pkg.MemM:
	case pkg.MemN:
	}
	log.Println(c)
}

func anotherV4x() {
	var c someone.SomeonesOwnY
	c = someone.P()
	switch c {
	case someone.SomeonesMM:
	case someone.SomeonesMN:
		// case someone.SomeonesMO:
	}
	log.Println(c)
}

// type alias
// assume you use it for code refactoring only
// assume you do one by one (first type change fully done; then only add new enum members for example
// don't add new enum members to old/pas type after new/future type is introduced
