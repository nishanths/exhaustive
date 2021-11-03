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

In the code snippet below, Biome is an enum type with 3 members. Enum values may
be specified using iota (they don't have to be explicit values, like in the
snippet), and enum members don't necessarily have to all be defined in the same
var or const block.

  type Biome int

  const (
      Tundra  Biome = 1
      Savanna Biome = 2
      Desert  Biome = 3
  )

Definition of exhaustiveness

An enum switch statement is exhaustive if all of the enum's members are listed
in the switch statement's cases.

For an enum type defined in the same package as the switch statement, both
exported and unexported enum members must be present in order to consider the
switch exhaustive. For an enum type defined in an external package it is
sufficient for just the exported enum members to be present in order to consider
the switch exhaustive.

Exhaustiveness checking strategies

There are two strategies for checking exhaustiveness: the "value" strategy
(which is the default) and the "name" strategy. The "value" strategy requires
that each independent enum value is listed in a switch statement to satisfy
exhaustiveness. The "name" strategy requires that each independent enum member
name is listed in a switch statement to satisfy exhaustiveness. The desired
exhaustiveness checking strategy can be specified using the "-checking-strategy"
flag.

To illustrate the difference between the strategies, consider the following enum
and switch statement.

  type AccessControl string

  const (
      AccessAll     AccessControl = "all"
      AccessAny     AccessControl = "any"
      AccessDefault AccessControl = AccessAll
  )

  func example(v AccessControl) {
      switch v {
          case AccessAll:
          case AccessAny:
      }
  }

The switch statement is not exhaustive when using the "name" strategy (because
the name AccessDefault is not listed), but it is exhaustive when using the
"value" strategy (because AccessDefault and AccessAll have the same value).

Notable flags

The "-default-signifies-exhaustive" boolean flag indicates to the analyzer
whether switch statements are to be considered exhaustive—even if all enum
members aren't listed in the switch statements cases—as long as a 'default' case
is present. The default value for the flag is false.

The "-check-generated" boolean flag indicates whether to check switch statements
in generated Go source files. The default value for the flag is false.

The "-ignore-enum-members" flag specifies a regular expression. Enum members
that match the regular expression do need to be listed in switch statements in
order for switch statements to be considered exhaustive. The supplied
regular expression is matched against the enum package's import path and the
enum member name combined in the following format: <import path>.<enum member
name>. For example: "github.com/foo/bar.Tundra", where the enum package's import
path is "github.com/foo/bar" and the enum member name is "Tundra".

The "-checking-strategy" flag specifies the exhaustiveness checking strategy to
use. The flag value must be either "value" (which is the default) or "name". See
discussion in the "Defintion of exhaustiveness" section for more details.

Skipping analysis

If the following directive comment:

  //exhaustive:ignore

is associated with a switch statement, the analyzer skips checking of the switch
statement and no diagnostics are reported. Note the lack of whitespace between
the comment marker ("//") and the comment text.

Additionally, no diagnostics are reported for switch statements in generated
files unless the "-check-generated" flag is enabled. (See
https://golang.org/s/generatedcode for definition of generated file).

Additionally, see the "-ignore-enum-members" flag.
*/
package exhaustive
