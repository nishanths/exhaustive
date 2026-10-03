// This testdata file contains a variety of valid and invalid
// declarations for use in tests. It is not meant to represent
// enumerated types and constants seen in production code.

package typ

import "unsafe"

type (
	t1 int // want t1:"^permittedtype$" t1:"^enumerated$" t1:"^elements:bit0 = 1, mask0 = 0, bit1 = 2, mask1 = 1, bit3 = 8, mask3 = 7$"
	t2 t1  // want t2:"^permittedtype$" t2:"^enumerated$" t2:"^elements:x1 = 10, x2 = 20, x7 = 70, x8 = 80, x08 = 80, x008 = 80, x9 = 90, x10 = 100, x11 = 110, x12 = 120, x13 = 130, x14 = 140, x15 = 141, x16 = 16$"
	t3 t2  // want t3:"^permittedtype$" t3:"^enumerated$" t3:"^elements:y1 = 10, y2 = 20, y5 = 50$"

	// test cases for all permitted underlying types
	_ bool       // want _:"^permittedtype$"
	_ byte       // want _:"^permittedtype$"
	_ complex64  // want _:"^permittedtype$"
	_ complex128 // want _:"^permittedtype$"
	_ float32    // want _:"^permittedtype$"
	_ float64    // want _:"^permittedtype$"
	_ int        // want _:"^permittedtype$"
	_ int8       // want _:"^permittedtype$"
	_ int16      // want _:"^permittedtype$"
	_ int32      // want _:"^permittedtype$"
	_ int64      // want _:"^permittedtype$"
	_ rune       // want _:"^permittedtype$"
	_ string     // want _:"^permittedtype$"
	_ uint       // want _:"^permittedtype$"
	_ uint8      // want _:"^permittedtype$"
	_ uint16     // want _:"^permittedtype$"
	_ uint32     // want _:"^permittedtype$"
	_ uint64     // want _:"^permittedtype$"
	_ uintptr    // want _:"^permittedtype$"

	q1 interface{}    // not permitted: underlying type is not *types.Basic
	q2 *int           // ditto
	q3 []byte         // ditto
	q4 unsafe.Pointer // not permitted: underlying type is *types.Basic, but the basic type is not a boolean, numeric, or string type
	q5 struct{}       // not permitted: underlying type is not *types.Basic

	m1 q1 // not permitted: underlying type is not *types.Basic
	m2 q4 // ditto

	s1[E any] int // not permitted: parameterized type

	r1        s1[int]      // want r1:"^permittedtype$"
	r2        s1[struct{}] // want r2:"^permittedtype$"
	r3[E any] s1[E]        // not permitted: parameterized type

	p1        r1           // want p1:"^permittedtype$"
	p2        r2           // want p2:"^permittedtype$"
	p3        r3[int]      // want p3:"^permittedtype$"
	p4        r3[struct{}] // want p4:"^permittedtype$"
	p5[E any] r3[E]        // not permitted: parameterized type

	a1        = t3           // not permitted: type alias
	a2        = p1           // ditto
	a3        = p5[int]      // ditto
	a4        = p5[struct{}] // ditto
	a5[E any] = p5[E]        // ditto
	a6        = q1           // ditto
	a7        = q4           // ditto

	n1        a1           // want n1:"^permittedtype$"
	n2        a2           // want n2:"^permittedtype$"
	n3        a3           // want n3:"^permittedtype$"
	n4        a4           // want n4:"^permittedtype$"
	n5        a5[struct{}] // want n5:"^permittedtype$"
	n6[E any] a5[E]
	n7        a6
	n8        a7
)

const (
	// (adopted from the Go spec document.)
	bit0, mask0 t1 = 1 << iota, 1<<iota - 1 // want bit0:"^enumerated$" bit0:"^elementof:t1$" mask0:"^enumerated$" mask0:"^elementof:t1$"
	bit1, mask1                             // want bit1:"^enumerated$" bit1:"^elementof:t1$" mask1:"^enumerated$" mask1:"^elementof:t1$"
	_, _                                    // not permitted: blank identifier
	bit3, mask3                             // want bit3:"^enumerated$" bit3:"^elementof:t1$" mask3:"^enumerated$" mask3:"^elementof:t1$"
)

const (
	_              = iota      // not permitted: blank identifier
	x1 t2          = iota * 10 // want x1:"^enumerated$" x1:"^elementof:t2$"
	x2                         // want x2:"^enumerated$" x2:"^elementof:t2$"
	x3 = iota * 10             // not permitted: (untyped constant), the type is not an enumerated type
	x4                         // ditto
	x5 int         = 50        // not permitted: predeclared type int
	x6             = 60        // ditto
	x7 t2          = iota * 10 // want x7:"^enumerated$" x7:"^elementof:t2$"
	x8                         // want x8:"^enumerated$" x8:"^elementof:t2$"
	// multiple names for same value.
	x08  t2 = x8 // want x08:"^enumerated$" x08:"^elementof:t2$"
	x008 t2 = 80 // want x008:"^enumerated$" x008:"^elementof:t2$"
	// mixing of iota and non-iota values.
	x9 t2 = 90 // want x9:"^enumerated$" x9:"^elementof:t2$"
)

const (
	y1 t3  = 10 // want y1:"^enumerated$" y1:"^elementof:t3$"
	y2 t3  = 20 // want y2:"^enumerated$" y2:"^elementof:t3$"
	y3     = 30 // not permitted: untyped constant, type is not an enumerated type
	y4 int = 40 // not permitted: predeclared type int
	y5 t3  = 50 // want y5:"^enumerated$" y5:"^elementof:t3$"
)

// continued from earlier const block.
const (
	x10 t2 = (iota + 10) * 10 // want x10:"^enumerated$" x10:"^elementof:t2$"
	x11 t2 = 110              // want x11:"^enumerated$" x11:"^elementof:t2$"
)

// const declarations without "( ... )" parentheses.
const x12 t2 = (iota + 12) * 10 // want x12:"^enumerated$" x12:"^elementof:t2$"
const x13 t2 = 130              // want x13:"^enumerated$" x13:"^elementof:t2$"
const x14 t2 = (iota + 14) * 10 // want x14:"^enumerated$" x14:"^elementof:t2$"

const astring = "0123456789abcdef"

// constant expressions; using built-in functions, using other
// constant identifiers.
const (
	x15 t2 = t2(max(0, 1)) + x14 // want x15:"^enumerated$" x15:"^elementof:t2$"
	x16 t2 = t2(len(astring))    // want x16:"^enumerated$" x16:"^elementof:t2$"
)

func f1() {
	// complex constants
	type t1 complex64 // want t1:"^permittedtype$" t1:"^enumerated$" t1:"^elements:x0 = \\(0 \\+ -1i\\), x1 = \\(0 \\+ 1i\\)$"
	const (
		x0 t1 = complex(0, -1) // want x0:"^enumerated$" x0:"^elementof:t1$"
		x1 t1 = 1i             // want x1:"^enumerated$" x1:"^elementof:t1$"

	)
}

type M1 int // want M1:"^permittedtype$" M1:"^enumerated$" M1:"^elements:V0 = 0, VV0 = 0, V1 = 1, V2 = 2, vv2 = 2, V3 = 3, vv3 = 3, V4 = 4, V5 = 5$"

const (
	V0, VV0 M1 = iota, iota // want V0:"^enumerated$" V0:"^elementof:M1$" VV0:"^enumerated$" VV0:"^elementof:M1$"
	V1      M1 = iota       // want V1:"^enumerated$" V1:"^elementof:M1$"
	V2, vv2 M1 = iota, iota // want V2:"^enumerated$" V2:"^elementof:M1$" vv2:"^enumerated$" vv2:"^elementof:M1$"
	V3, vv3                 // want V3:"^enumerated$" V3:"^elementof:M1$" vv3:"^enumerated$" vv3:"^elementof:M1$"
	V4      M1 = iota       // want V4:"^enumerated$" V4:"^elementof:M1$"
	V5                      // want V5:"^enumerated$" V5:"^elementof:M1$"
)

// exported alias for unexported type; with constants in same block.
type A2 = t11

const (
	// Note: These constants (except Wa3) are of type t11 and they
	// are declared in the same block as type t11.
	// Thus these are enumerated constants of type t11.
	//
	Wa0 = w0 // want Wa0:"^enumerated$" Wa0:"^elementof:t11$"
	Wa1 = w1 // want Wa1:"^enumerated$" Wa1:"^elementof:t11$"
	Wa2 = w2 // want Wa2:"^enumerated$" Wa2:"^elementof:t11$"
	Wa3 = 3
	// Note: missing Wa4 = w4
)

type t11 int // want t11:"^permittedtype$" t11:"^enumerated$" t11:"^elements:Wa0 = 0, Wa1 = 1, Wa2 = 2, w0 = 0, w1 = 1, w2 = 2, w4 = 4$"

const (
	w0 t11 = iota // want w0:"^enumerated$" w0:"^elementof:t11$"
	w1            // want w1:"^enumerated$" w1:"^elementof:t11$"
	w2            // want w2:"^enumerated$" w2:"^elementof:t11$"
	_
	w4 // want w4:"^enumerated$" w4:"^elementof:t11$"
)
