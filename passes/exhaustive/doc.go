/*
Package exhaustive defines a static analysis that checks that
switch statements in Go source code are exhaustive. It can check
expression switches, in which the switch expression type is an
enumerated type, and type switches.

The complete analysis is implemented in multiple passes; see the
packages in the 'passes/*' directories.
This package scans the syntax tree of Go packages for eligible
switch statements and checks that the switch statements are
exhaustive.

# Expression switches

An expression switch is eligible to be checked if the switch
expression type is either an [enumerated type] or an alias for an
enumerated type.

An expression switch is exhaustive if every value of the
enumerated type of the switch expression is included in the case
clauses. If the expression switch statement and the declarations
of the enumerated type (and constants) are in different
packages, then it is necessary to include only the values
corresponding to exported constants. A case expression must be a
[constant expression] to contribute towards making the
expression switch exhaustive.

Given the declarations

	type T int

	const (
		X0 T = iota
		X1
		X2
	)

the following expression switch is not exhaustive.

	var t T

	switch t {
	case X0:
	case X2:
	}

The diagnostic is

	expression switch not exhaustive: missing cases: X1

Note that including a default case does not make a switch
exhaustive. See flag -d to control this behavior.

# Type switches

Every type switch is eligible to be checked,
with the following exceptions. Type switches in which the
interface in the type assertion is the empty interface or the
predeclared type error are not checked.

A type switch is exhaustive if every type T or *T that implements
the interface in the type assertion is included in the case
clauses, where T is the type name of a top-level, non-interface
defined type, or an alias declaration for such a type, either in
the current package or in its (recursively) imported
dependencies. If the type switch statement and the type
declaration are in different packages, then it is necessary to
include the type in the case clauses only if its name is exported.

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

The diagnostic is (as of go1.27)

	type switch not exhaustive: missing cases: *types.Basic, *types.TypeParam

It is not necessary to include a case clause for nil (i.e. "case
nil: ...") for a type switch to be exhaustive; this behavior can
be controlled via the flags. A similar note about the default
case applies to type switches, as described in the section on
expression switches.

# Additional syntax elements

If configured via flags, the analysis can check whether additional
kinds of elements in the syntax tree, such as map literals where
the key type is an enumerated type, are exhaustive. See flag -check.

# Handling of type aliases

For an expression switch, given the following alias declaration
and enumerated type declaration

	package a
	import "b"
	type T = b.U
	const (
		X0 = b.X0
		X1 = b.X1
	)

	package b
	type U int
	const (
		X0 U = iota
		X1
	)

the following expression switch will be checked and is
exhaustive. The enumerated type for this expression switch is
b.U, which is the type denoted by the alias type a.T of the
switch expression.

	var t a.T
	switch t {
	case a.X0: // or b.X0
	case a.X1: // or b.X1
	}

For a type switch, either the type name of the defined type or
the type name of an alias declaration may be specified in the
case clauses. Given these declarations

	type S struct{}
	func (*S) m()
	type A = S

the following type switch is exhaustive

	var x interface{ m() }
	switch x.(type) {
	case *A: // or *S
	}

# Flags

The exhaustive analyzer supports the following flags.

	[-casenil] [-d] [-defrequire] [-e] [-g] [-check string] [-constignore regexp] [-typeignore regexp] [-typeonly regexp]

The flags are described below.

The -check flag specifies the kinds of elements in the syntax
tree that the analysis should check. The argument is a
comma-separated list of one or more of the following names.
The default argument is "switch".

	Name          Description

	switch        check that expression switch statements are exhaustive.
	mapliteral    check that composite literals of underlying type map (or pointer to) are exhaustive.
	typeswitch    check that type switch statements are exhaustive.

The -d flag changes the behavior of the analysis such that
including a default case makes a switch statement exhaustive
regardless of other cases clauses. The -e flag restricts
checking to only those switch statements that have an
'//exhaustive:enforce' comment directive. The -g flag enables
checking switch statements found in generated files. The
-casenil flag applies to type switches; it specifies whether
type switches must include a 'nil' case.

The -defrequire flag additionally checks that each checked
expression switch statement includes a default case; this flag
is off by default and is currently supported only for expression
switches.

The -typeignore, -typeonly, and -constignore flags accept a
regular expression pattern in Go package regexp syntax; these
flags may be repeated to specify multiple patterns. If the
switch expression type name[^^] is matched by a -typeignore
regexp, then that switch statement will not be checked. If
-typeonly flags are specified then only those switch statements
in which the switch expression type name is matched by a
-typeonly regexp will be checked. If a type name is matched by
both a -typeonly regexp and a -typeignore regexp the -typeonly
match wins. Constant names matched by a -constignore regexp do
not have to be included in case clauses for an expression switch
to be exhaustive.

These regexp flags are applicable only when the type or constant
declaration in question has a top-level declaration. The type
name or constant name that the analysis provides to the regexp
matching routine is always fully qualified by the import path.
For example, given these declarations

	package bar // import "example.org/foo/bar"
	type S int
	const Y0 S = 0

the fully qualified name of the type S is
"example.org/foo/bar.S" and of the constant Y0 is
"example.org/foo/bar.Y0".

The enumerated analyzer, defined in [package enumerated],
provides additional flags that control the discovery of
enumerated types and enumerated constants. See its documentation
for details.

[^^] In a type switch, this instead refers to the type name of
the (possibly) named interface type in the type assertion.

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
expected.

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
[enumerated type]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated
[constant expression]: https://golang.org/ref/spec#Constant_expressions

[basic interface]: https://golang.org/ref/spec#Basic_interfaces
*/
package exhaustive
