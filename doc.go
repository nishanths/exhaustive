/*
Package exhaustive provides an analyzer that checks exhaustiveness of enum
switch statements in Go source code.

Definition of enum

The Go language spec does not provide an explicit definition for an enum. For
the purpose of this analyzer, an enum type is any named type (a.k.a. defined
type) whose underlying type is an integer (includes byte and rune), a float, or
a string type. An enum type has associated with it constants of this named type;
these constants constitute the enum members.

In the example below, Biome is an enum type with 3 members.

    type Biome int

    const (
        Tundra  Biome = 1
        Savanna Biome = 2
        Desert  Biome = 3
    )

For a constant to be an enum member for an enum type, the constant must be
declared in the same scope as the enum type. Note that the scope requirement
implies that only constants declared in the same package as the enum type's
package can constitute the enum members for the enum type.

Enum member constants for a given enum type don't necessarily have to all be
declared in the same const block. Constant values may be specified using iota,
using explicit values, or by any means of declaring a valid Go const. It is
allowed for multiple enum members for a given enum type to have the same
constant value.

Definition of exhaustiveness

A switch statement that switches on a value of an enum type is exhaustive if all
of the enum type's members are listed in the switch statement's cases. If
multiple enum member constants have the same constant value, it is sufficient
that any one of these same-valued members is listed.

For an enum type defined in the same package as the switch statement, both
exported and unexported enum members must be listed to satisfy exhaustiveness.
For an enum type defined in an external package, it is sufficient that only the
exported enum members are listed.

Only identifiers denoting constants (e.g. Tundra) and qualified identifiers
denoting constants (e.g. mypkg.Constant) listed in switch statement cases can
contribute towards satisfying exhaustiveness. Literal constant values (e.g. 42,
"Sunday"), struct fields (e.g. obj.f), etc. will not.

Type aliases

The analyzer handles type aliases for an enum type in the following manner.
Consider the example below. T2 is a enum type, and T1 is an alias for T2. Note
that we don't term T1 itself an enum type; it is only an alias for an enum
type.

    package pkg
    type T1 = otherpkg.T2
    const (
        A = otherpkg.A
        B = otherpkg.B
    )

    package otherpkg
    type T2 int
    const (
        A T2 = 1
        B T2 = 2
    )

A switch statement that switches on a value of type T1 (which, in reality, is
just an alternate spelling for type T2) is exhaustive if all of T2's enum
members are listed in the switch statement's cases. Additionally, the same
conditions described in the previous section for same-valued enum members and
for exported/unexported enum members apply here.

It is worth noting that, though T1 and T2 are identical types, only constants
declared in the same scope as type T2's scope can constitute the enum type T2's
enum members. In the example, otherpkg.A and otherpkg.B are T2's enum
members. T1, as mentioned earlier, is not an enum type; consequently
the concept of enum members does not apply to it.

Advanced notes

Recall from an earlier section that for a constant to be an enum member for an
enum type, the constant must be declared in the same scope as the enum type.
However it is valid, both to the Go type checker and to this analyzer, for any
constant of the right type to be listed in the cases of an enum switch statement
(it does not necessarily have to be a constant declared in the same scope/package
as the enum type's scope/package).

Such a constant can contribute towards satisfying switch statement
exhaustiveness if it has the same constant value as an actual enum member
constant—the constant can take the place of the same-valued enum member constant
in the switch statement's cases. This behavior is particularly useful when a
type alias is involved: A forwarding const declaration (such as pkg.A, in type
T1's package) can take the place of the actual enum member const (such as
otherpkg.A, in type T2's package) in the switch statement's cases.

Flags

Notable flags for the analyzer are described below.
All of these flags are optional.

    flag                            type    default value

    -check-generated                bool    false
    -default-signifies-exhaustive   bool    false
    -ignore-enum-members            string  (none)
    -package-scope-only             bool    false

If the -check-generated flag is enabled, switch statements in generated Go
source files are also checked. Otherwise, by default, switch statements in
generated files are not checked.

If the -default-signifies-exhaustive flag is enabled, the presence of a
'default' case in switch statements always satisfies exhaustiveness, even if all
enum members are not listed. It is recommended that you do not enable this flag;
enabling it generally defeats the purpose of exhaustiveness checking.

The -ignore-enum-members flag specifies a regular expression in Go syntax. Enum
members matching the regular expression are ignored, i.e. matching enum member
names don't have to be listed in switch statement cases to satisfy
exhaustiveness. The specified regular expression is matched against an enum
member name inclusive of the enum package import path: for example, if the enum
package import path is "example.com/pkg" and the member name is "Tundra", the
supplied regular expression will be matched against the string
"example.com/pkg.Tundra".

If the -package-scope-only flag is enabled, the analyzer only finds enums
defined in package scopes, and consequently only switch statements that switch
on package-scoped enums will be checked for exhaustiveness. By default, the
analyzer finds enums defined in all scopes, and checks switch statements that
switch on all these enums.

Skipping analysis

To skip checking of a specific switch statement, associate the following comment
with the switch statement.

    //exhaustive:ignore

For example:

    //exhaustive:ignore
    switch v {

Note the lack of whitespace between the comment marker ("//") and the comment
text.

To ignore specific enum members, see the -ignore-enum-members flag.

By default, the analyzer skips checking of switch statements in generated
Go source files. Use the -check-generated flag to change this behavior.
See https://golang.org/s/generatedcode for the definition of generated file.
*/
package exhaustive
