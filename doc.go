/*
Package exhaustive provides an analyzer that checks exhaustiveness of enum
switch statements in Go code.

The analyzer also provides fixes to make the offending switch statements
exhaustive (see "Fixes" section). For the related command line program, see the
"cmd/exhaustive" subpackage.

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

The -default-signifies-exhaustive boolean flag indicates to the analyzer whether
switch statements are to be considered exhaustive (even if all enum members
aren't listed in the switch statements cases) as long as a 'default' case is
present. The default value for the flag is false.

The -check-generated boolean flag indicates whether to check switch statements
in generated Go source files. The default value for the flag is false.

The -ignore-pattern flag specifies a regular expression. Enum members that match
the regular expression do not require a case clause in switch statements to
satisfy exhaustiveness. The regular expression is matched against enum member
names inclusive of the enum package's import path. For example, in
"github.com/foo/bar.Tundra", the enum package's import path is
"github.com/foo/bar" and the enum member name is Tundra.

For details on the -fix boolean flag, see the next section.

Fixes

The analyzer suggests fixes for a switch statement if it is not exhaustive.
The suggested fix always adds a single case clause for the missing enum member
values.

  case MissingA, MissingB, MissingC:
      panic(fmt.Sprintf("unhandled value: %v", v))

where v is the expression in the switch statement's tag (in other words, the
value being switched upon). If the switch statement's tag is a function or a
method call the analyzer does not suggest a fix, as reusing the call expression
in the panic/fmt.Sprintf call could be mutative.

The rationale for the fix using panic is that it might be better to fail loudly
on existing unhandled or impossible cases than to let them slip by quietly
unnoticed. An even better fix, of course, may be to manually inspect the sites
reported by the package and handle the missing cases as necessary.

Imports will be adjusted automatically to account for the "fmt" dependency.

Skipping analysis

If the following directive comment:

  //exhaustive:ignore

is associated with a switch statement, the analyzer skips
checking of the switch statement and no diagnostics are reported.

Additionally, no diagnostics are reported for switch statements in
generated files (see https://golang.org/s/generatedcode for definition of
generated file), unless the -check-generated flag is enabled.

Additionally, see the -ignore-pattern flag.
*/
package exhaustive
