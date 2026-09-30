/*
Package enumerated finds declarations of enumerated types and
enumerated constants.

# Enumerated types

An enumerated type is a [defined type] whose underlying type is
either a boolean, numeric, or string type. The type definition
must not specify type parameters. An enumerated type cannot
directly be one of the predeclared types; it must be a new,
distinct type.

Examples of valid and invalid types:

	type T1 int         // T1: valid
	type T2 string      // T2: valid
	type T3 T1          // T3: valid
	type Q1 []int       // Q1: invalid: underlying type []int is not a boolean, numeric, or string type
	type S1[E any] int  // S1: invalid: parameterized
	int                 // int: invalid: predeclared type

The possible values of a given enumerated type are the values of
each enumerated constant of that type. A type that is valid but
has an empty set of values is not considered an enumerated type
by this package.

# Enumerated constants

Declared constants where the type of the constant is an
enumerated type form enumerated constants of that type. The type
and the constants must be declared in the same [block] to be
considered by the analysis. Constants declared with the blank
identifier are ignored. In the following example, the type T1 is
an enumerated type and the constants X0, X1, and X2 are
enumerated constants of that type.

	type T1 int

	const (
		X0 T1 = iota
		X1
		X2
	)

The enumerated constants of a given enumerated type may be
declared across multiple [ConstDecl] productions. The constant
value of an enumerated constant may be any allowed [constant
expression]; this includes literal values, values generated with
iota, and values containing constant identifiers. It is
permitted for multiple enumerated constants of a given enumerated
type to have the same constant value.

# Flags

The analyzer's flags can control the discovery of enumerated
types and enumerated constants.

The synopsis of the flags is:

	[-B] [-i] [-p]

The -i flag directs the analysis to only consider constants
that are declared using iota in their value expressions.

The -B flag directs the analysis to ignore constants that are
declared using bitwise operators (& | ^ &^ << >>) in their value
expressions.

The -p flag directs the analysis to only consider types and
constants that are declared at package level.

[block]: https://golang.org/ref/spec#Blocks
[ConstDecl]: https://golang.org/ref/spec#Constant_declarations
[constant expression]: https://golang.org/ref/spec#Constant_expressions
[defined type]: https://golang.org/ref/spec#Type_definitions
*/
package enumerated
