package general

import "test/general/sub"

type t struct{}    // want t:"^typedecl$"
type u interface{} // want u:"^typedecl$"
type (
	v sub.V // want v:"^typedecl$"
)

type a = t        // want a:"^typedecl$"
type b = struct{} // want b:"^typedecl$"
type c = *a       // want c:"^typedecl$"
type (
	d = *sub.V // want d:"^typedecl$"
)

const n = 0

func f1() {
	type m struct{}
	type a = m
}

func f2[t any]() {}
