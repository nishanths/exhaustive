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
eligible to be checked.

An expression switch is exhaustive if each value of the
enumerated type of the switch expression is included in the case
clauses. If the expression switch statement and the declarations
of the enumerated type (and constants) are in different packages,
then it is necessary to include only the values corresponding to
exported constants. A case expression must be a constant to
contribute towards making the expression switch exhaustive.

Normally, including a default case does not make a switch
statement exhaustive. See flag -d to control this behavior.

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

# Type switches

All valid type switch statements (as of go1.27) are eligible to
be checked, except for the following special cases: Type switches
in which the interface in the type assertion is the empty
interface or the predeclared type error are not checked.

A type switch is exhaustive if every type T or *T that implements
the interface in the type assertion is included in the case
clauses, where T is a non-interface type declared at top level
in the current package or its (recursively) imported dependencies.
If the type switch statement and the type declaration are in
different packages, then it is only necessary to include the type
in the case clauses if its name is exported.

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

Normally, it is not necessary to include a 'nil' case (i.e. "case
nil: ...") for a type switch to be exhaustive, but this behavior
can be adjusted via the flags. The same note about the default
case, as described in the section on expression switches, also
applies to type switches.

# Flags

The exhaustive analyzer supports the following flags.

	[-casenil] [-d] [-defrequire] [-e] [-g] [-check string] [-constignore regexp] [-typeignore regexp] [-typeonly regexp]

The flags are described below.

The -check flag specifies the kinds of elements in the syntax
tree that the analysis should check are exhaustive. The argument
is a comma-separated list of one or more of the following names.
The default argument is "switch". For instance, to check both
expression switches and type switches, specify:
-check "switch,typeswitch".

	name          description

	switch        expression switch statements
	typeswitch    type switch statements
	mapliteral    composite literals of underlying type map (or pointer to)

The -d flag changes the behavior of the analysis such that
including a default case makes a switch statement exhaustive
regardless of other case clauses.

The -casenil flag applies to type switches; it specifies whether
type switches must include a 'nil' case (i.e. "case nil: ...") to
be exhaustive.

The -e flag restricts checking to only those switch statements
that have an '//exhaustive:enforce' comment directive. The -g
flag enables checking of switch statements in generated files.
The -defrequire flag requires that each checked expression switch
includes a default case even if it is exhaustive.

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

The enumerated analyzer, defined in [package enumerated],
provides additional flags that control the discovery of
enumerated types and enumerated constants. See its documentation
for details.

[*] In a type switch, this instead refers to the type name, if
any, of the interface in the type assertion.

# Comment directives

The analysis supports the following comment directives.

	//exhaustive:ignore
	//exhaustive:enforce

A comment directive has effect if placed in the source code on an
element that is configured to be checked (see flag -check). The
analysis does not check an element if it has an 'ignore'
directive. The 'enforce' directive forces the check of an element
that may otherwise not be checked due to configuration elsewhere.

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

	//exhaustive:ignore
	var m1 = map[T]bool{
		X0: true,
		X1: true,
	}

# Examples

For more basic examples see the code snippets in the earlier sections.

The following example illustrates the effect of flag -d. Normally,
the expression switch below would not be exhaustive despite the
inclusion of a default case (missing case for X2). But with flag -d
specified, the expression switch is considered exhaustive.

	type T int

	const (
		X0 T = iota
		X1
		X2
	)

	func f(v T) {
		switch v {
		case X0:
		case X1:
		default:
		}
	}

The following examples illustrate handling of aliases. In the
expression switch below, the enumerated type is effectively b.B,
which is the type denoted by the alias type a.A of the switch
expression. The switch is exhaustive because each value of the
enumerated type b.B (i.e., 0 and 1, corresponding to b.X0 and
b.X1, respectively) is included in the case clauses.

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

In a type switch, an alias type may be specified in the case
clauses. The type switch below is exhaustive.

	type I interface{ m() }

	type T struct{}
	func (*T) m() {}

	type A = T

	func f(v I) {
		switch v.(type) {
		case *A: // *A is identical to *T
		}
	}

[package enumerated]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated
[enumerated type]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated
*/
package exhaustive
