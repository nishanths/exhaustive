/*
Package exhaustive provides an analyzer that checks exhaustiveness of enum
switch statements in Go code.

Definition of enum

The Go language spec does not provide an explicit definition for enums.
For the purpose of this program, an enum type is a package-level named type
whose underlying type is an integer (includes byte and rune), a float, or
a string type. An enum type must have associated with it one or more
package-level variables of the named type in the package. These variables
constitute the enum's members.

In the code snippet below, Biome is an enum type with 3 members. You may
also use iota instead of explicitly specifying values, and enum members
don't necessarily have to be all defined in the same var/const block.

  type Biome int

  const (
      Tundra  Biome = 1
      Savanna Biome = 2
      Desert  Biome = 3
  )

Definition of exhaustiveness

An enum switch statement is exhaustive if it has cases for each of the enum's
members.

For an enum type defined in the same package as the switch statement, both
exported and unexported enum members must be present in order to consider the
switch exhaustive. For an enum type defined in an external package it is
sufficient for just the exported enum members to be present in order to consider
the switch exhaustive.

Notable flags

The "-default-signifies-exhaustive" boolean flag indicates to the analyzer
whether switch statements are to be considered exhaustive—even if all enum
members aren't listed in the switch statements cases—as long as a 'default' case
is present. The default value for the flag is false.

The "-check-generated" boolean flag indicates whether to check switch statements
in generated Go source files. The default value for the flag is false.

The "-ignore-pattern" flag specifies a regular expression. Enum members that
match the regular expression do not require a case clause in switch statements
in order for the switch statements to be considered exhaustive. Effectively, the
enum member is ignored when checking for exhaustiveness. The supplied regular
expression is matched against the enum package's import path and the enum member
name combined, e.g. "github.com/foo/bar.Tundra", where the enum package's import
path is "github.com/foo/bar" and the enum member name is "Tundra".

Skipping analysis

If the following directive comment:

  //exhaustive:ignore

is associated with a switch statement, the analyzer skips checking of the switch
statement and no diagnostics are reported. Note the lack of whitespace between
the comment marker ("//") and the comment text.

Additionally, no diagnostics are reported for switch statements in generated
files unless the "-check-generated" flag is enabled. (See
https://golang.org/s/generatedcode for definition of generated file).

Additionally, see the "-ignore-pattern" flag.
*/
package exhaustive

// TODO: add docs for upcoming -by-name flag.
// TODO: add docs to "Definition of exhaustiveness" section for by value vs. by name checking.
