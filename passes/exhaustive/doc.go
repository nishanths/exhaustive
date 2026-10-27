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

Given an [expression switch] in which the switch expression type is
an [enumerated type], the expression switch is considered
exhaustive if every value in the [value set] of the enumerated
type is specified in the cases.
If the expression switch and the declaration of the enumerated
type and its constants are in different packages, then it is
sufficient that the values corresponding to exported constants
are specified in the cases.
The case expressions must be constant values.

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

The analysis reports a diagnostic if an expression switch is not exhaustive.
The diagnostic for the above expression switch is

	missing cases in expression switch: darcs

The presence of a default case does not make an expression switch exhaustive.
See flag -d to control this behavior.

# Type switches

Given a [type switch] in a package p, the "knowable type set" of
the type switch is the set of types T or *T that implements the
interface type I of the expression in the type switch guard,
where T is a non-interface type declared at top level in package
p or its (recursively) imported dependencies. If the type switch
and the type declaration of T are in different packages, then the
types T and *T are included in the knowable type set only if T is
exported.

A type switch is considered exhaustive if every type in its
knowable type set is specified in the cases.

If an interface type J, but not a type parameter, is specified as a
type in a case, it is equivalent to specifying each non-interface
type T or *T that implements J.

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

The analysis reports a diagnostic if a type switch is not exhaustive.
The diagnostic (as of go1.27) for the above type switch is

	missing cases in type switch: *types.Basic, *types.TypeParam

The presence of a default case does not make a type switch exhaustive.
See flag -d to control this behavior.
It is not necessary for a type switch to include a nil case (i.e.
"case nil:") to be exhaustive, but this behavior can be adjusted
via the flags.

As special case exceptions, the analysis does not check type
switches in which the interface type of the expression in the
type switch guard is the empty interface or the predeclared type
error.

The analysis does not support generic types that implement an
interface. Diagnostics do not mention these types as missing in
the cases of a type switch.

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

The -casenil flag specifies whether type switches must include a
nil case (i.e. "case nil:") to be exhaustive.
The -e flag restricts checking to only those switch statements
that have an '//exhaustive:enforce' comment directive.
The -g flag enables checking of switch statements in generated
files.
The -defrequire flag requires that each checked expression switch
includes a default case even if it is exhaustive.

The -impossible flag checks for impossible or unreachable cases
in type switches. See the "Examples" section.

The -constignore, -typeignore, and -typeonly flags accept a
regular expression pattern in Go package regexp syntax. These
flags may be repeated to specify multiple patterns.
Constant names matched by a -constignore regexp do not have to be
specified in the cases of an expression switch to be exhaustive.
If the switch expression type name[*] is matched by a -typeignore
regexp, then that switch statement will not be checked.
If -typeonly flags are specified, then only those switch
statements in which the switch expression type name is matched by
a -typeonly regexp will be checked.
If a type name is matched by both a -typeonly regexp and a
-typeignore regexp, the -typeonly match wins.

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

[*] In a type switch, this instead refers to the type name (if it
has one) of the interface type of the expression in the type
switch guard.

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
flag. The type *Loop implements the interface Statement. The type
switch has a case for Statement before the case for *Loop.

	type Node interface{ node() }
	type Statement interface{ Node; statement() }

	type Literal struct{}
	func (*Literal) node() {}

	type Loop struct{}
	func (*Loop) node() {}
	func (*Loop) statement() {}

	func f(v Node) {
		switch v.(type) {
		case Statement:
		case *Literal:
		case *Loop:
		}
	}

The diagnostic is

	unreachable type switch case *Loop (matched by earlier case Statement)

The following examples illustrate handling of aliases.

In the following expression switch, the enumerated type is b.B,
which is the type denoted by the alias a.A of the switch
expression. The expression switch is exhaustive; every value in
the value set of b.B is specified in the cases.

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

In a type switch, an alias may be used in a case instead of the
original type. The following type switch must specify types *T
and *U in the cases to be exhaustive. The type *A, which is
identical to *T, is specified in a case instead of *T. The type
switch is exhaustive.

	type m interface{ m() }

	type T struct{}
	func (*T) m() {}

	type A = T

	type U struct{}
	func (*U) m() {}

	func f(v m) {
		switch v.(type) {
		case *A:
		case *U:
		}
	}

[expression switch]: https://golang.org/ref/spec#Expression_switches
[type switch]: https://golang.org/ref/spec#Type_switches
[package enumerated]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated
[enumerated type]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated
[value set]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated#hdr-Value_sets
*/
package exhaustive
