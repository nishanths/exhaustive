// Package exhaustive defines a static analysis that checks that
// switch statements in Go source code are exhaustive. The analysis
// can check expression switches, if the switch expression
// type is an enumerated type, and type switches.
//
// The complete analysis is implemented in multiple passes; see the
// packages under the 'passes' directory of the module.
// This package scans the syntax tree of Go packages for switch
// statements and checks that the switch statements are exhaustive.
//
// # Expression switches
//
// Given an [expression switch] in which the type of the switch
// expression is an [enumerated type] T, the analysis considers
// the expression switch to be exhaustive if every value in the
// [value set] of T is listed in the cases.
// If the expression switch and the declaration of T (and its
// constants) are in different packages, then it is sufficient
// that the values corresponding to exported constants are listed
// in the cases.
//
// A given case expression must be a constant to satisfy a
// corresponding value in the value set of T.
//
// The following expression switch is not exhaustive.
//
//	type vcs int
//
//	const (
//		bazaar vcs = iota
//		fossil
//		git
//		mercurial
//		subversion
//		darcs
//	)
//
//	func f(v vcs) {
//		switch v {
//		case bazaar:
//		case fossil:
//		case git:
//		case mercurial:
//		case subversion:
//		}
//	}
//
// The analysis reports a diagnostic if an expression switch is not exhaustive.
//
// The diagnostic for the above expression switch is
//
//	missing cases in expression switch: darcs
//
// The presence of a default case will not satisfy missing cases.
// See flag -d to control this behavior.
//
// # Type switches
//
// Given a [type switch] in a package p, the "knowable type set" of
// the type switch is the set of non-interface types T or *T that
// implements the interface type I of the switch expression,
// where T is a type declared at top level in package p or its
// (recursively) imported dependencies. If the type switch and
// the type declaration of T are in different packages, then the
// types T and *T are included in the knowable type set only if T
// is exported.
//
// The analysis considers a type switch to be exhaustive if
//
//   - every type in its knowable type set is listed in the cases; and
//
//   - none of the types listed in the cases are ineffectual.
//     (See the "Ineffectual cases" section for details.)
//
// If an interface type S, but not a type parameter, is listed in
// a case, it is equivalent to listing in the cases each
// (non-interface) type in the knowable type set that implements S.
//
// The following type switch is not exhaustive.
//
//	import "go/types"
//
//	func f(t types.Type) bool {
//		switch t.(type) {
//		case *types.Array,
//			*types.Chan,
//			*types.Interface,
//			*types.Map,
//			*types.Pointer,
//			*types.Signature,
//			*types.Slice,
//			*types.Struct,
//			*types.Tuple,
//			*types.Union:
//			return false
//		case *types.Alias, *types.Named:
//			return true
//		default:
//			panic("internal error: unhandled type")
//		}
//	}
//
// The analysis reports a diagnostic if a type switch is not exhaustive.
//
// The diagnostic (as of go1.27) for the above type switch is
//
//	missing cases in type switch: *types.Basic, *types.TypeParam
//
// The analysis does not support generic types that implement an
// interface. Diagnostics do not mention these types as missing in
// the cases of a type switch.
//
// As special cases, if the interface type of the switch
// expression is either the empty interface or
// the error built-in interface, the analysis does not evaluate
// missing cases in these type switches.
//
// The presence of a default case will not satisfy missing cases.
// See flag -d to control this behavior.
// It is not necessary for a type switch to include a nil case
// ("case nil:") to be exhaustive, but this behavior can be
// adjusted via the flags.
//
// # Ineffectual cases
//
// The type switch analysis requires that none of the types
// listed in the cases are ineffectual. To illustrate, consider
// the following type switch.
//
//	import "go/ast"
//
//	func count(file *ast.File) {
//		var decls, fields, funcs int
//		for n := range ast.Preorder(file) {
//			switch n.(type) {
//			case ast.Decl:
//				decls += 1
//			case *ast.Field:
//				fields += 1
//			case *ast.FuncDecl:
//				funcs += 1
//			default:
//				// do not care about other types
//			}
//		}
//		/* ... */
//	}
//
// Note that the case that contains *ast.FuncDecl is present
// after the case that contains ast.Decl. Also note that
// *ast.FuncDecl implements the ast.Decl interface.
//
// If the interface variable n of the switch expression has
// dynamic type *ast.FuncDecl, then the ast.Decl case matches
// and control does not reach the *ast.FuncDecl case. If n does
// not have dynamic type *ast.FuncDecl (including when n is a nil
// interface value), then control may reach the *ast.FuncDecl
// case but the dynamic type of n is not matched by the
// *ast.FuncDecl type in the case.
//
// The diagnostic for the above type switch is
//
//	order of cases makes it impossible to match type *ast.FuncDecl in the case (due to type ast.Decl in earlier case)
//
// The analysis detects ineffectual cases in limited scenarios.
// The analysis does not handle scenarios where the type listed
// in a case is the type constraint interface for a type
// parameter.
// The analysis does not check for issues that are already
// reported as errors by 'go build'.
//
// # Flags
//
// The exhaustive analysis supports the following flags.
//
//	[-casenil] [-d] [-defrequire] [-e] [-g] [-check string]
//	[-constignore regexp] [-typeignore regexp] [-typeonly regexp]
//
// The flags are described below.
//
// -casenil:
// Type switches must include a nil case to be exhaustive.
//
// -check string:
// Analyze the specified elements.  The argument is a
// comma-separated list of one or more of the following
// element names: switch, typeswitch, mapliteral.  The
// default argument is "switch".  The analysis checks
// that the specified elements are exhaustive.
//
//	switch        analyze expression switch statements
//	typeswitch    analyze type switch statements
//	mapliteral    analyze composite literals of underlying type map (or pointer to)
//
// -constignore regexp:
// Constant names matched by the regexp pattern do not
// have to be listed in the cases of an expression switch.
// The flag may be repeated to specify multiple regexp
// patterns.
//
// -d:
// Including a default case in a switch statement satisfies
// all missing cases.
//
// -defrequire:
// Check that analyzed expression switches always include
// a default case.
//
// -e:
// Analyze a switch statement only if it has an
// '//exhaustive:enforce' comment directive.
//
// -g:
// Analyze switch statements in generated files.
//
// -typeignore regexp:
// Do not analyze a switch statement if the type name (if
// there is one) of the switch expression is matched by
// the regexp pattern.  The flag may be repeated to specify
// multiple regexp patterns.
//
// -typeonly regexp:
// Analyze a switch statement only if the type name (if
// there is one) of the switch expression is matched by
// the regexp pattern.  The flag may be repeated to specify
// multiple regexp patterns.  If a type name is matched
// by both a -typeonly regexp and a -typeignore regexp,
// the -typeonly match wins.
//
// The regexp arguments to the flags must be in Go package regexp syntax.
//
// The -constignore, -typeignore, and -typeonly flags are
// applicable only when the type or constant in question is
// declared at top level.
// The type name or constant name that the analysis provides to
// the regular expression match routine is always fully qualified
// by the import path.
// For example, given the declarations
//
//	package bar // import "example.org/foo/bar"
//	type S int
//	const Y0 S = 0
//
// the fully qualified name of the type S is
// "example.org/foo/bar.S" and of the constant Y0 is
// "example.org/foo/bar.Y0".
//
// The enumerated analysis, defined in [package enumerated],
// provides additional flags that control the discovery of
// enumerated types and enumerated constants. See its
// documentation for details.
//
// # Comment directives
//
// The analysis supports the following comment directives.
//
//	//exhaustive:ignore
//	//exhaustive:enforce
//
// A comment directive has effect if placed in the source code
// adjacent to an element that has been configured to be analyzed
// (see flag -check). An element with an 'ignore' directive will
// not be analyzed. An 'enforce' directive forces analysis of an
// element that may otherwise not be analyzed due to
// configuration elsewhere.
//
// The directive must be present in the comment group immediately
// before the target switch statement, and there must be no nodes or
// empty lines between the comment group and the switch statement.
// For a map literal, similar requirements apply. For a map literal,
// the directive must be present in the comments of the nearest
// ancestor *ast.AssignStmt node, *ast.ValueSpec node, or
// *ast.GenDecl node. See the source code for exact details.
//
// In general, the following comment directive placements should
// work as expected.
//
//	func f(v V) {
//		//exhaustive:ignore
//		switch v {
//		case X0:
//		case X1:
//		}
//	}
//
//	//exhaustive:ignore
//	var m = map[T]bool{
//		X0: true,
//		X1: true,
//	}
//
// [package enumerated]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated
// [expression switch]: https://golang.org/ref/spec#Expression_switches
// [type switch]: https://golang.org/ref/spec#Type_switches
// [enumerated type]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated
// [value set]: https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated#hdr-Value_sets
package exhaustive
