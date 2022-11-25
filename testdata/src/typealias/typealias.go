package typealias

import (
	"enum/typealias"
	"enum/typealias/anotherpkg"
	"enum/typealias/otherpkg"
)

func t1() typealias.T1 { return 0 }
func t2() typealias.T2 { return 0 }
func t9() typealias.T9 { return 0 }

const (
	C int          = 1
	D int          = 2
	E typealias.T1 = 4
)
const (
	F int          = 1
	G int          = 2
	H typealias.T2 = 4
	I typealias.T3 = 5
)
const (
	J int          = 1
	K int          = 2
	L typealias.T8 = 4
	M typealias.T9 = 5
)

func _a() {
	v := t1()
	switch v {
	}

	_ = map[typealias.T1]int{
		0: 0,
	}
}
func _b() {
	v := t2()
	switch v {
	}

	_ = map[typealias.T2]int{
		0: 0,
	}
}
func _c() {
	switch t9() {
	}

	_ = map[typealias.T9]int{
		0: 0,
	}
}

// --

func t4() typealias.T4   { return 0 }
func t10() typealias.T10 { return 0 }
func t6() typealias.T6   { return 0 }
func t15() typealias.T15 { return 0 }

const (
	N typealias.T4 = 4
	O typealias.T5 = 5
)
const (
	P typealias.T10 = 4
	Q typealias.T11 = 5
)

const (
	R typealias.T6 = 4
	S typealias.T7 = 5
)
const (
	T typealias.T15 = 4
	U typealias.T16 = 5
)

func _d() {
	switch t4() {
	}

	_ = map[typealias.T4]int{
		0: 0,
	}
}
func _e() {
	switch t10() { // want "^missing cases in switch of type typealias.T11: typealias.T11_A, typealias.T11_B$"
	}

	_ = map[typealias.T10]int{ // want "^missing keys in map of key type typealias.T11: typealias.T11_A, typealias.T11_B$"
		0: 0,
	}
}
func _f() {
	switch t6() {
	}

	_ = map[typealias.T6]int{
		0: 0,
	}
}
func _g() {
	switch t15() { // want "^missing cases in switch of type typealias.T11: typealias.T11_A, typealias.T11_B$"
	}

	_ = map[typealias.T15]int{ // want "^missing keys in map of key type typealias.T11: typealias.T11_A, typealias.T11_B$"
		0: 0,
	}
}

// --

func t12() typealias.T12 { return 0 }
func t13() typealias.T13 { return 0 }
func t14() typealias.T14 { return 0 }
func t17() typealias.T17 { return 0 }

const (
	V typealias.T12 = 6
	W otherpkg.T5   = 7
)
const (
	X typealias.T13 = 6
	Y otherpkg.T11  = 7
)
const (
	Z  typealias.T14 = 6
	AA otherpkg.T7   = 7
)
const (
	BB typealias.T17 = 6
	CC otherpkg.T16  = 7
)

func _h() {
	switch t12() {
	}

	_ = map[typealias.T12]int{
		0: 0,
	}
}
func _i() {
	switch t13() { // want "^missing cases in switch of type otherpkg.T11: otherpkg.T11_A, otherpkg.T11_B, otherpkg.T11_C$"
	}

	_ = map[typealias.T13]int{ // want "^missing keys in map of key type otherpkg.T11: otherpkg.T11_A, otherpkg.T11_B, otherpkg.T11_C$"
		0: 0,
	}
}
func _j() {
	switch t14() {
	}

	_ = map[typealias.T14]int{
		0: 0,
	}
}
func _k() {
	switch t17() { // want "^missing cases in switch of type otherpkg.T11: otherpkg.T11_A, otherpkg.T11_B, otherpkg.T11_C$"
	}

	_ = map[typealias.T17]int{ // want "^missing keys in map of key type otherpkg.T11: otherpkg.T11_A, otherpkg.T11_B, otherpkg.T11_C$"
		0: 0,
	}
}

// --

func t18() typealias.T18 { return 0 }

const (
	DD typealias.T18 = 8
	EE otherpkg.T19  = 9
	FF anotherpkg.T1 = 10
)

func _l() {
	v := t18()
	switch v { // want "^missing cases in switch of type anotherpkg.T1: anotherpkg.T1_A$"
	}

	_ = map[typealias.T18]int{ // want "^missing keys in map of key type anotherpkg.T1: anotherpkg.T1_A$"
		0: 0,
	}
}

// --

func d1() typealias.D1 { return struct{}{} }

func _m() {
	v := d1()
	switch v {
	}

	_ = map[typealias.D1]int{
		struct{}{}: 0,
	}
}
