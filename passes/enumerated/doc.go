/*
Package enumerated finds declarations of enumerated types and
enumerated constants.

# Enumerated Types and Enumerated Constants

A [defined type] is permitted to be an enumerated type if it
satisfies the following conditions. The underlying type of the
defined type must be either a numeric or string type.
The type definition must not specify type parameters. The type
must be a new type; the predeclared types are
not directly permitted to be enumerated types.

Examples of types permitted to be enumerated types:

	type T1 int         // T1: permitted
	type T2 int32       // T2: permitted
	type T3 rune        // T3: permitted
	type T4 string      // T4: permitted
	type T5 T1          // T5: permitted
	type T6 T5          // T6: permitted

Examples of types not permitted to be enumerated types:

	type Q1 []int       // Q1: not permitted: underlying type []int is not a numeric or string type
	type S1[E any] int  // S1: not permitted: parameterized
	int32               // int32: not permitted: predeclared type
	rune                // rune: not permitted: predeclared type

A declared constant whose type is a permitted enumerated type forms
an enumerated constant of that type. A given type and the
constants must be declared in the same [block] to be considered
by the analysis. A constant with the [blank identifier]
as its name is ignored. An enumerated type must have at least one enumerated
constant of the type. A type without any enumerated constants
is not considered an enumerated type by the analysis.

In the following example, the type T1 is an enumerated type. The
constants X0, X1, and X2 are enumerated constants of the type.

	type T1 int

	const (
		X0 T1 = iota
		X1
		X2
	)

The possible values of an enumerated type are the combined set of
values of the enumerated constants of the type.

The enumerated constants of a given enumerated type may be
declared across multiple [ConstDecl] productions. The constant
value of an enumerated constant may be any value allowed by the
language. It is permitted for multiple enumerated constants of a
given enumerated type to have the same value.

# Flags

The analyzer's flags can control the discovery of enumerated
types and enumerated constants.

The synopsis of the flags is:

	[-B] [-i] [-p]

The -i flag directs the analysis to only consider constants
that are declared using iota in their values.

The -B flag directs the analysis to ignore constants that are
declared using bitwise operators (& | ^ &^ << >>) in their
values.

The -p flag directs the analysis to only consider types and
constants that are declared at package level.

[block]: https://golang.org/ref/spec#Blocks
[ConstDecl]: https://golang.org/ref/spec#Constant_declarations
[defined type]: https://golang.org/ref/spec#Type_definitions
[blank identifier]: https://golang.org/ref/spec#Blank_identifier
*/
package enumerated
