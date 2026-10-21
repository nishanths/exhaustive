/*
Package exhaustive defines a static analysis that checks that
switch statements in Go source code are exhaustive. The analysis
can check expression switches, in which the switch expression
type is an enumerated type, and type switches.

The complete analysis is implemented in multiple passes; see the
packages under the 'passes' directory.
This package scans the syntax tree of Go packages for eligible
switch statements and checks that the switch statements are
exhaustive.

# Expression switches

Expression switch statements in which the switch expression type
is an [enumerated type] or an alias for an enumerated type are
checked.

An expression switch is considered exhaustive if each value of the
enumerated type of the switch expression is included in the case
clauses.

The following expression switch is not exhaustive.

	type vcs int

	const (
		bazaar vcs = iota
		fossil
		git
		mercurial
		subversion
		darcs
	)

	func f(v vcs) {
		switch v {
		case bazaar:
		case fossil:
		case git:
		case mercurial:
		case subversion:
		}
	}

The diagnostic is

	missing cases in expression switch: darcs

If the expression switch statement and the declarations
of the enumerated type (and constants) are in different packages,
then it is necessary to include only the values corresponding to
exported constants.

A case expression must be a constant to contribute towards making
the expression switch exhaustive.

The presence of a default case does not make a switch statement
exhaustive. See flag -d to control this behavior.

# Type switches

Every type switch statement is checked, except for the following
special case exceptions: type switches in which the interface in
the type assertion is the empty interface or the predeclared type
error are not checked.

A type switch is considered exhaustive if each type T or *T that
implements the interface in the type assertion is included in the
case clauses, where T is a non-interface type declared at top
level in the current package or its (recursively) imported
dependencies.

The following type switch is not exhaustive.

	import "go/types"

	func f(t types.Type) bool {
		switch t.(type) {
		case *types.Array,
			*types.Chan,
			*types.Interface,
			*types.Map,
			*types.Pointer,
			*types.Signature,
			*types.Slice,
			*types.Struct,
			*types.Tuple,
			*types.Union:
			return false
		case *types.Alias, *types.Named:
			return true
		default:
			panic("internal error: unhandled type")
		}
	}

The diagnostic (as of go1.27) is

	missing cases in type switch: *types.Basic, *types.TypeParam

If the type switch statement and the type declaration of a type T
are in different packages, then it is only necessary to include T
or *T in the case clauses if T is an exported type.

If an interface type or a type parameter (whose type constraint
is an interface type) is specified as a type in a case clause, it
is equivalent to specifying each type T or *T that implements
that interface.

Generic types that have methods and implement an interface are
not supported by this analysis. These types will not be reported
in diagnostics.

It is not necessary to include a 'nil' case (i.e. "case nil: ")
for a type switch to be exhaustive, but this behavior can be
adjusted via the flags. The behavior regarding the default case,
as described in the section on expression switches, also applies
to type switches.

# Flags

The exhaustive analysis supports the following flags.

	[-casenil] [-d] [-defrequire] [-e] [-g] [-impossible] [-check string] [-constignore regexp] [-typeignore regexp] [-typeonly regexp]

The flags are described below.

The -check flag specifies the kinds of elements in the syntax
tree that the analysis should check are exhaustive. The argument
is a comma-separated list of one or more of the following names.
The default argument is "switch".

	name          description

	switch        expression switch statements
	typeswitch    type switch statements
	mapliteral    composite literals of underlying type map (or pointer to)

If the -d flag is specified, then switch statements that
include a default case are always considered exhaustive.

The -casenil flag applies to type switches; it specifies whether
type switches must include a 'nil' case (i.e. "case nil: ") to
be exhaustive. The -e flag restricts checking to only those
switch statements that have an '//exhaustive:enforce' comment
directive. The -g flag enables checking of switch statements in
generated files. The -defrequire flag requires that each checked
expression switch includes a default case even if it is exhaustive.

The -impossible flag checks for impossible or
unreachable cases in type switches. See the "Examples" section.

The -constignore, -typeignore, and -typeonly flags accept a
regular expression pattern in Go package regexp syntax. These
flags may be repeated to specify multiple patterns. Constant
names matched by a -constignore regexp do not have to be included
in case clauses for an expression switch to be exhaustive.
If the switch expression type name[*] is matched by a -typeignore
regexp then that switch statement will not be checked. If
-typeonly flags are specified then only those switch statements
in which the switch expression type name is matched by a
-typeonly regexp will be checked. If a type name is matched by
both a -typeonly regexp and a -typeignore regexp, the -typeonly
match wins.

These regexp flags are applicable only when the type or constant
declaration in question is declared at top level. The type
name or constant name that the analysis provides to the regexp
matching routine is always fully qualified by the import path.
For example, given the declarations

	package bar // import "example.org/foo/bar"
	type S int
	const Y0 S = 0

the fully qualified name of the type S is
"example.org/foo/bar.S" and of the constant Y0 is
"example.org/foo/bar.Y0".

The enumerated analysis, defined in [package enumerated],
provides additional flags that control the discovery of
enumerated types and enumerated constants. See its documentation
for details.

[*] In a type switch, this instead refers to the type name, if
any, of the interface in the type assertion.

# Comment directives

The analysis supports the following comment directives.

	//exhaustive:ignore
	//exhaustive:enforce

A comment directive has effect if placed in the source code
adjacent to an element that has been configured to be checked
(see flag -check). An element with an 'ignore' comment directive
will not be checked. An 'enforce' directive forces the check of
an element that may otherwise not be checked due to configuration
elsewhere.

The directive must be present in the comment group immediately
before the target switch statement, and there must be no nodes or
empty lines between the comment group and the switch statement.
For a map literal, similar requirements apply. For a map literal,
the directive must be present in the comments of the nearest
ancestor *ast.AssignStmt node, *ast.ValueSpec node, or
*ast.GenDecl node. See the source code for exact details.

In general, the following comment directive placements should
work as expected.

	//exhaustive:ignore
	switch t {
	case X0:
	case X1:
	}

	//exhaustive:ignore
	var m1 = map[T]bool{
		X0: true,
		X1: true,
	}

# Examples

See the code snippets in the earlier sections for basic examples.

The following example shows the functionality of the -impossible
flag. The *loop type implements the stmt interface. The type
switch has a case for stmt before the case for *loop.

	type node interface{ node() }
	type stmt interface{ node; stmt() }

	type loop struct{}
	func (*loop) node() {}
	func (*loop) stmt() {}

	type literal struct{}
	func (*literal) node() {}

	func f(v node) {
		switch v.(type) {
		case stmt:
		case *loop:
		case *literal:
		}
	}

The diagnostic is:

	unreachable type switch case *loop (matched by earlier case stmt)

The following examples illustrate handling of aliases.

In the expression switch below, the enumerated type is
effectively b.B, which is the type denoted by the alias type a.A
of the switch expression. The switch is exhaustive because each
value of the enumerated type b.B (i.e., 0 and 1, corresponding to
b.X0 and b.X1, respectively) is included in the case clauses.

	package a
	import "b"
	type A = b.B
	const X0 = b.X0
	const X1 = b.X1

	package b
	type B int
	const X0 B = 0
	const X1 B = 1

	package foo
	import "a"
	func f(v a.A) {
		switch v {
		case a.X0:
		case a.X1:
		}
	}

In a type switch, an alias type may be used in a case instead of
the original type. The type switch below is exhaustive.

	type I interface{ m() }

	type T struct{}
	func (*T) m() {}

	type U struct{}
	func (*U) m() {}

	type A = T

	func f(v I) {
		switch v.(type) {
		case *A: // note that *A is identical to *T
		case *U:
		}
	}

[package enumerated]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated
[enumerated type]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated
*/
package exhaustive
