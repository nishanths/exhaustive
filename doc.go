/*
Package exhaustive provides an analyzer that checks exhaustiveness of enum
switch statements in Go code.

Definition of enum

The Go language spec does not provide an explicit definition for enums. For the
purpose of this analyzer, an enum type is a package-level named type whose
underlying type is an integer (includes byte and rune), a float, or a string
type. An enum type must have associated with it one or more package-level
variables of the named type in the same package. These variables constitute the
enum's members.

In the code snippet below, Biome is an enum type with 3 members.

  type Biome int

  const (
      Tundra  Biome = 1
      Savanna Biome = 2
      Desert  Biome = 3
  )

Enum member values may also be specified using iota; they don't necessarily have
to be explicit values, like in the snippet. Enum members don't necessarily have
to all be defined in the same var or const block.

Definition of exhaustiveness

An enum switch statement is exhaustive if all of the enum's members are listed
in the switch statement's cases.

For an enum type defined in the same package as the switch statement, both
exported and unexported enum members must be present in order to consider the
switch statement exhaustive. For an enum type defined in an external package, it
is sufficient for just the exported enum members to be present in order to
consider the switch statement exhaustive.

Exhaustiveness checking strategies

There are two strategies for checking exhaustiveness: the "name" strategy and
the "value" strategy (which is the default). The name strategy requires that
each independent enum member name is listed in a switch statement to satisfy
exhaustiveness. On the other hand, the value strategy only requires that each
independent enum value is listed in a switch statement to satisfy
exhaustiveness.

To illustrate the difference between the two strategies, consider the
enum and the switch statement in the code snippet below.

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

The switch statement is not exhaustive when using the name strategy (because the
name AccessDefault is not listed in the switch), but it is exhaustive when using
the value strategy (because AccessDefault and AccessAll have the same value, and
it suffices that one of them is listed in the switch).

Notable flags

The notable flags used by the analyzer are:

  -default-signifies-exhaustive

If enabled, the presence of a "default" case in switch statements satisfies
exhaustiveness, even if all enum members are not listed.

  -check-generated

If enabled, switch statements in generated Go source files are also checked.
Otherwise switch statements in generated files are ignored by default.

  -ignore-enum-members <regex>

Specifies a regular expression; enum members matching the regular expression are
ignored. Ignored enum members don't have to be present in switch statements to
satisfy exhaustiveness. The regular expression is matched against enum member
names inclusive of the enum package import path, e.g.
"github.com/foo/bar.Tundra" where the enum package import path is
"github.com/foo/bar" and the enum member name is "Tundra".

  -checking-strategy <strategy>

Specifies the exhaustiveness checking strategy, which must be one of "name" or
"value" (default). For details see section: Exhaustiveness checking strategies.

Skipping analysis

If the following comment:

  //exhaustive:ignore

is associated with a switch statement, the analyzer skips inspection of the
switch statement and no diagnostics are reported. Note the lack of whitespace
between the comment marker ("//") and the comment text.

Additionally, no diagnostics are reported for switch statements in generated
files unless the -check-generated flag is enabled. See
https://golang.org/s/generatedcode for the definition of generated file.

Additionally, see the -ignore-enum-members flag, which can be used
to ignore specific enum members.
*/
package exhaustive
