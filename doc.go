/*
Package exhaustive provides an analyzer that checks exhaustiveness of enum
switch statements in Go source code.

Definition of enum

The Go language spec does not provide an explicit definition for an enum. For
the purpose of this analyzer, an enum type is a named type (a.k.a. defined type)
whose underlying type is an integer (includes byte and rune), a float, or a
string type. An enum type has associated with it constants of the named type;
these constants constitute the enum's members.

In the example below, Biome is an enum type with 3 members.

    type Biome int

    const (
        Tundra  Biome = 1
        Savanna Biome = 2
        Desert  Biome = 3
    )

For a constant to be an enum member, it must be declared in the same scope as
the enum type. This implies that all members of an enum type must be present
in the same package as the enum type (if a constant of the enum type is defined
in a package different from the enum type's package, the constant will not
constitute an enum member).

Enum member constants for a given enum type don't necessarily have to all be
declared in the same const block. The constant values may be specified using
using iota, using explicit values, or by any means of declaring a valid const.

Definition of exhaustiveness

An switch statement that switches on a value of an enum type is exhaustive if
all of the enum type's members (by constant value) are listed in the switch
statement cases.

For an enum type defined in the same package as the switch statement, both
exported and unexported enum members must be listed to satisfy exhaustiveness.
For an enum type defined in an external package, it is sufficient that only the
exported enum members be listed to satisfy exhaustiveness.

Exhaustiveness and type aliases

The type alias proposal says that in a type alias declaration:

    type T1 = T2

T1 is merely an alternate spelling for T2, and nearly all analysis of code
involving T1 proceeds by first expanding T1 to T2 [*]. For this analyzer, it
means that a switch statement that switches on a value of type T1 is, in effect,
switching on a value of type T2.

If T2 or its underlying type were, then a switch statement that switches on a
value of type T1 (which, in effect, is type T2) is exhaustive if all of the type
T2's enum members are listed in the switch statement cases.

Note that the switch statement cases don't necessarily have to list only T2
constants; it is valid to list T1 constants (or a mixture of T2 and T1
constants) since T1 and T2 are type aliases.

If a T1 constant is listed, the analzyer attempts to match the T1 constant with
a T2 enum member by name and value. For example, the switch statement below
is exhaustive even though it lists no T2 members explicitly:

    package bar
    type T2 rune // T2 is an enum type with enum members: A, B.
    const (
        A = 'a'
        B = 'b'
    )

    package foo
    import "bar"
    type T1 = bar.T2 // T1 is an alias for enum type T2.
    const (
       A = bar.A
       B = bar.B
    )
    func ReturnsT1() T1 { ... }

    package main
    import "foo"
    func x() {
        switch foo.ReturnsT1() { // Switch tag's type is bar.T2, according to Go type analysis.
        case foo.A:              // The analyzer considers foo.A equivalent to bar.A (same name, same value).
        case foo.B:              // The analyzer considers foo.B equivalent to bar.B (same name, same value).
        }
    }

As a side note, it is worth highlighting that only constants declared in the
same package as the type T2 can constitute T2's enum members (as defined in
section 'Definition of enum').

[*] https://go.googlesource.com/proposal/+/master/design/18130-type-alias.md#proposal

Flags

The notable flags used by the analyzer are described below.
All of these flags are optional.

    Flag                            Type    Default value
    -check-generated                bool    false
    -default-signifies-exhaustive   bool    false
    -ignore-enum-members            string  (none)
    -package-scope-only             bool    false

If the -check-generated flag is enabled, switch statements in generated Go
source files are also checked. Otherwise switch statements in generated files
are ignored by default.

If the -default-signifies-exhaustive flag is enabled, the presence of a 'default'
case in switch statements satisfies exhaustiveness, even if all enum members are
not listed. It is recommended that you do not enable this flag; enabling it
generally defeats the purpose of exhaustiveness checking.

The -ignore-enum-members flag specifies a regular expression, in the syntax
accepted by the regexp package. Enum members matching the regular expression
are ignored, i.e.  matching enum member names don't have to be listed in switch
statements to satisfy exhaustiveness. The specified regular expression is
matched against an enum member name inclusive of the enum package import path:
for example, if the import path is "example.com/pkg" and the member name is
"Tundra", the supplied regular expression will be matched against the string
"example.com/pkg.Tundra".

If the -package-scope-only flag is enabled, the analyzer only finds enums
defined in package scopes, and consequently only switch statements that switch
on package-scoped enums will be checked for exhaustiveness. By default, the
analyzer finds enums defined in all scopes, and checks switch statements that
switch on all these enums.

Skipping analysis

To skip checking of a specific switch statement, associate the following comment
with the switch statement. Note the lack of whitespace between the comment
marker ("//") and the comment text.

    //exhaustive:ignore

For example:

    switch v { //exhaustive:ignore

To ignore specific enum members, see the -ignore-enum-members flag.

By default, the analyzer skips checking of switch statements in generated
Go source files. Use the -check-generated flag to change this behavior.
See https://golang.org/s/generatedcode for the definition of generated file.
*/
package exhaustive
