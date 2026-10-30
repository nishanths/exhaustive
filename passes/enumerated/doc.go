// Package enumerated defines a static analysis that finds
// declarations of enumerated types and enumerated constants.
//
// # Enumerated types
//
// A [defined type] is allowed to be an enumerated type if it
// satisfies the following conditions. The underlying type of the
// defined type must be either a numeric or string type. The type
// definition must not specify type parameters. A predeclared type
// cannot directly be an enumerated type.
//
// These examples show allowed and disallowed enumerated types:
//
//	type T0 int         // T0: allowed
//	type T1 int32       // T1: allowed
//	type T2 rune        // T2: allowed
//	type T3 string      // T3: allowed
//	type T4 T0          // T4: allowed
//	type T5 T4          // T5: allowed
//
//	type Q0 []int       // Q0: disallowed: underlying type []int is not a numeric or string type
//	type S0[E any] int  // S0: disallowed: type definition specifies type parameters
//	int32               // int32: disallowed: predeclared type
//	rune                // rune: disallowed: predeclared type
//
// Each enumerated type must have at least one enumerated constant
// of the type, otherwise the type is disregarded by the analysis.
//
// # Enumerated constants
//
// Declared constants whose type is an allowed enumerated type
// form enumerated constants of the type. An enumerated type and
// the constants of the type must be declared in the same [block] to
// be considered by the analysis. Constants with blank identifier
// names are disregarded.
//
// In the following example, the type T0 is an enumerated type. The
// constants X0, X1, and X2 are enumerated constants of the type.
//
//	type T0 int
//
//	const (
//		X0 T0 = iota
//		X1
//		X2
//	)
//
// The enumerated constants of a given enumerated type may be
// declared across multiple [ConstDecl] productions. The value in
// the constant declaration may be any expression permitted by
// the language. Multiple enumerated constants of a given
// enumerated type may have the same value.
//
// # Value sets
//
// The "value set" of an enumerated type T consists of the values of
// each enumerated constant of the type T.
//
// For example, the value set of the enumerated type T0 in the above
// code snippet consists of the values 0, 1, and 2, corresponding to
// the constants X0, X1, and X2.
//
// # Flags
//
// The synopsis of the flags is:
//
//	[-B] [-i] [-p]
//
// The flags are described below.
//
// -B:
// Disregard enumerated constants that are declared using
// bitwise operators (& | ^ &^ << >>) in their values.
//
// -i:
// Only consider enumerated constants that are declared
// using iota in their values.
//
// -p:
// Only consider enumerated types and enumerated constants
// that are declared at top level.
//
// [block]: https://golang.org/ref/spec#Blocks
// [ConstDecl]: https://golang.org/ref/spec#Constant_declarations
// [defined type]: https://golang.org/ref/spec#Type_definitions
package enumerated
