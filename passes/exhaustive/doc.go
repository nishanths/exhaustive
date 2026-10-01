/*
Package exhaustive checks that expression switch statements, in
which the type of the switch expression is an enumerated type,
are exhaustive.

The analysis consists of two passes. The first pass, implemented
in [package enumerated], finds declarations of enumerated types
and enumerated constants. The second pass, implemented in this
package, scans a syntax tree for eligible switch statements and
checks that these switch statements are exhaustive.

The analysis produces a diagnostic if a checked switch statement
is not exhaustive. For example:

	switch not exhaustive: missing cases: X1, X2

# Expression switch statements

An expression switch statement is eligible to be checked if the
type of its switch expression is an enumerated type (or an alias
for an enumerated type). An expression switch statement is
exhaustive if every value of the [enumerated type] in the switch
expression is included in the case clause expressions.

Given the declarations:

	type T int

	const (
		X0 T = iota
		X1
		X2
	)

the following switch statement will be checked and is exhaustive:

	var t T

	switch t {
	case X0:
	case X1:
	case X2:
	}

If the switch statement is in a different package from the
package that declares the enumerated type and the constants, then
only the values corresponding to exported constants have to be
included in the case expressions for the switch statement to be
exhaustive. The case expression must be a [constant expression] to
contribute towards making the switch statement exhaustive.

Note that, by default, including a default case in a switch
statement does not automatically make a switch statement
exhaustive. See flag -d to control this behavior.

# Handling of type aliases

For example, given an enumerated type

	package newpkg

	type S int

	const (
		X0 S = iota
		X1
	)

and the following alias declaration

	package oldpkg

	import "newpkg"

	type T = newpkg.S

	const (
		X0 = newpkg.X0
		X1 = newpkg.X1
	)

the following switch statement will be checked and is exhaustive.
The effective enumerated type of the switch expression is
newpkg.S, which is the type on the right-hand side of the alias
declaration upon following the alias chain.

	package foo

	import "oldpkg"

	var t oldpkg.T

	switch t {
	case oldpkg.X0: // equivalently newpkg.X0
	case oldpkg.X1: // equivalently newpkg.X1
	}

# Other syntax tree elements

If configured via flags, the analysis can check that additional
kinds of elements in the syntax tree, such as map literals, are
exhaustive. See flag -check.

# Flags

The enumerated analyzer defined in [package enumerated] provides
a set of flags to control the discovery of enumerated types and
enumerated constants. See its documentation for details.

The exhaustive analyzer defined in this package supports the
following flags.

	[-d] [-defrequire] [-e] [-g] [-check string] [-constignore regexp] [-typeignore regexp] [-typeonly regexp]

The flags are described below.

The -check flag specifies the kinds of elements in the syntax
tree that the analysis should check. The argument is a
comma-separated list of one or more of these names: switch,
mapliteral. The default argument is "switch".

	Name          Description

	switch        check that expression switch statements
	              are exhaustive

	mapliteral    check that composite literals of underlying
	              type map (or pointer to) are exhaustive

The -d changes the behavior of the analysis such that including
a default case makes a switch statement exhaustive regardless of
other case clauses. The -e flags restricts checking to only
those switch statements that have an '//exhaustive:enforce'
comment directive. The -g flag enables checking switch statements
found in generated files. The -defrequire flag additionally
checks that each checked switch statement includes a default
case; the default value of the flag is off.

The -typeignore, -typeonly, and -constignore flags accept a
regular expression pattern in Go package regexp syntax; these
flags may be repeated to specify multiple patterns. If the type
name of a switch expression is matched by a -typeignore regexp,
then that switch statement will not be checked. If -typeonly
flags are specified then only those switch statements in which
the type name of the switch expression is matched by a -typeonly
regexp will be checked. By default the analysis checks all
eligible switch statements.

Constant names matched by a -constignore regexp do not have to be
included in case expressions for the switch statement to be
exhaustive.

These regexp flags are applicable only when the type or constant
declaration in question has a package-level declaration. The type
name or constant name that the analyzer provides to the regexp
matching routine is a fully qualified name and always includes
the import path. For example, given these declarations

	package bar // import "example.org/foo/bar"
	type S int
	const Y0 S = 0

the fully qualified name of the type S is "example.org/foo/bar.S"
and the fully qualified name of the constant Y0 is
"example.org/foo/bar.Y0".

# Comment directives

The analysis supports the following comment directives.

	//exhaustive:ignore
	//exhaustive:enforce

TODO: Some comment directives are not documented.

The comment directives have effect if placed in the source on an
element that is configured to be checked (see flag -check). The
analysis does not check a switch statement if it has an 'ignore'
directive. The 'enforce' directive forces the check of a switch
statement that may otherwise be ignored due to configuration
elsewhere.

For switch statements the comment must be associated with the
switch statement node as defined by func NewCommentMap in package
go/ast. The comment has effect only for the switch statement that
it is associated with, and not for any descendant switch
statements. For map literals the analysis considers the line
comments, doc comments, and associated comments of the nearest
ancestor *ast.AssignStmt node or *ast.ValueSpec and *ast.GenDecl
nodes. See the source code for exact details.

In general, the following comment placements should work as
expected:

	//exhaustive:ignore
	switch t {
	case X0:
	case X1:
	}

	//exhaustive:enforce
	var m1 = map[T]bool{
		X0: true,
		X1: true,
	}

[package enumerated]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated
[enumerated type]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated#hdr-Enumerated_types
[constant expression]: https://golang.org/ref/spec#Constant_expressions
*/
package exhaustive
