/*
Package enumerated defines a static analysis that finds
declarations of enumerated types and enumerated constants.

# Enumerated Types and Enumerated Constants

A [defined type] is allowed to be an enumerated type if it
satisfies the following conditions. The underlying type of the
defined type must be either a numeric or string type. The type
definition must not specify type parameters. A predeclared type
cannot directly be an enumerated type.

These examples show allowed and disallowed enumerated types:

	type T0 int         // T0: allowed
	type T1 int32       // T1: allowed
	type T2 rune        // T2: allowed
	type T3 string      // T3: allowed
	type T4 T0          // T4: allowed
	type T5 T4          // T5: allowed

	type Q0 []int       // Q0: disallowed: underlying type []int is not a numeric or string type
	type S0[E any] int  // S0: disallowed: parameterized
	int32               // int32: disallowed: predeclared type
	rune                // rune: disallowed: predeclared type

Declared constants whose type is an allowed enumerated type
form enumerated constants of the type. An enumerated type and
the constants of the type must be declared in the same [block] to
be considered by the analysis. Constants with blank identifier
names are ignored. Each enumerated type must have at least one
enumerated constant of the type, otherwise the type is
disregarded by the analysis.

The possible values of an enumerated type are the combined set of
values of each enumerated constant of the type.

In the following example, the type T0 is an enumerated type. The
constants X0, X1, and X2 are enumerated constants of the type.

	type T0 int

	const (
		X0 T0 = iota
		X1
		X2
	)

The enumerated constants of a given enumerated type may be
declared across multiple [ConstDecl] productions. The constant
value of an enumerated constant can be any value allowed by the
language (e.g. iota, constant literals, other constant
identifiers). The analysis permits multiple enumerated constants
of a given enumerated type to have the same value.

# Flags

The analysis flags can control the discovery of enumerated types
and enumerated constants.

The synopsis of the flags is:

	[-B] [-i] [-p]

The -i flag directs the analysis to only consider constants
that are declared using iota in their values.

The -B flag directs the analysis to ignore constants that are
declared using bitwise operators (& | ^ &^ << >>) in their
values.

The -p flag directs the analysis to only consider types and
constants that are declared at top level.

[block]: https://golang.org/ref/spec#Blocks
[ConstDecl]: https://golang.org/ref/spec#Constant_declarations
[defined type]: https://golang.org/ref/spec#Type_definitions
*/
package enumerated
