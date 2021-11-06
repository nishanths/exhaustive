/*
Package exhaustive provides an analyzer that checks exhaustiveness of enum
switch statements in Go code.

Definition of enum

The Go language spec does not provide an explicit definition for enums. For the
purpose of this analyzer, an enum type is a named (defined) type whose
underlying type is an integer (includes byte and rune), a float, or a string
type. An enum type must have associated with it one or more constants of the
named type. These constants constitute the enum's members.

In the example below, Biome is an enum type with 3 members.

	type Biome int

	const (
			Tundra  Biome = 1
			Savanna Biome = 2
			Desert  Biome = 3
	)

For a constant to be an enum member, it must be declared in the same scope as
the enum type. Enum member constants don't necessarily have to all be declared
in the same const block. Enum member constant values may be specified using iota
or using explicit values (like in the example).

The analyzer's behavior is undefined for type aliases. This may change in the
future.

Definition of exhaustiveness

An enum switch statement is exhaustive if all of the enum's members are listed
in the switch statement's cases.

For an enum type defined in the same package as the switch statement, both
exported and unexported enum members must be listed to satisfy exhaustiveness.
For an enum type defined in an external package, it is sufficient that only the
exported enum members be listed to satisfy exhaustiveness.

Flags

The notable flags used by the analyzer are described below.
All of these flags are optional.

	Flag name						Type	Default value
	-check-generated				bool	false
	-default-signifies-exhaustive	bool	false
	-ignore-enum-members			string	(none)
	-package-scope-only				bool	false

If the -check-generated flag is enabled, switch statements in generated Go
source files are also checked. Otherwise switch statements in generated files
are ignored by default.

If the -default-signifies-exhaustive flag is enabled, the presence of a "default"
case in switch statements satisfies exhaustiveness, even if all enum members are
not listed. It is recommended that you do not enable this flag; enabling it
generally defeats the purpose of exhaustiveness checking.

The -ignore-enum-members flag specifies a regular expression, in the syntax
accepted by the Go regexp package. Enum members matching the regular expression
are ignored, i.e. these enum member names don't have to be listed in switch
statements to satisfy exhaustiveness. The specified regular expression is
matched against an enum member name inclusive of the enum package import path:
for example, "example.com/pkg.Tundra" where where the import path is
"example.com/pkg" and the enum member name is "Tundra".

If the -package-scope-only flag is enabled, the analyzer only finds enums
defined in package scopes, and consequently only switch statements that switch
on package-scoped enums will be checked for exhaustiveness. By default, the
analyzer finds enums defined in all scopes, and checks switch statements that
switch on all these enums.

Skipping analysis

To skip analysis of a specific switch statement, associate the following
comment with the switch statement. Note the lack of whitespace
between the comment marker ("//") and the comment text.

	//exhaustive:ignore

To ignore specific enum members, see the -ignore-enum-members flag.

By default, the analyzer skips analysis of switch statements in generated
Go source files. Use the -check-generated flag to change this behavior.
See https://golang.org/s/generatedcode for the definition of generated file.
*/
package exhaustive
