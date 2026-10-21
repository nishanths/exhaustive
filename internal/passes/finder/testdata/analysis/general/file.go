package general

import "test/analysis/general/sub"

type t struct{}    // want t:"^typedecl$"
type u interface{} // want u:"^typedecl$"
type (
	v sub.V // want v:"^typedecl$"
)

type l[P1, P2 any] struct { x P1; y P2 } // want l:"^typedecl$"

type a = t        // want a:"^typedecl$"
type b = struct{} // want b:"^typedecl$"
type c = *a       // want c:"^typedecl$"
type (
	d = *sub.V // want d:"^typedecl$"
)

type e0 = l[any, any] // want e0:"^typedecl$"
type e1[P2 any] = l[any, P2] // want e1:"^typedecl$"

const n = 0

func f1() {
	type m struct{}
	type a = m
}

func f2[t any]() {}
