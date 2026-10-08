The exhaustive static analysis checks that
switch statements in Go source code are exhaustive. It can check
expression switches, in which the switch expression type is
an enumerated type, and type switches.

The complete analysis is implemented in multiple passes; see the
packages under the 'passes' directory.
Package enumerated finds declarations of enumerated types and
enumerated constants.
Package exhaustive scans the syntax tree of Go packages for eligible
switch statements and checks that the switch statements are
exhaustive.

Documentation:

<https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated>

<https://pkg.go.dev/github.com/nishanths/exhaustive/passes/exhaustive>

# Usage

The exhaustive command can be installed with 'go install'.

	go install github.com/nishanths/exhaustive/cmd/exhaustive@latest

The synopsis of the command is:

	exhaustive [-B] [-casenil] [-d] [-defrequire] [-e] [-g] [-i] [-p]
	           [-check string] [-constignore regexp] [-typeignore regexp]
	           [-typeonly regexp] [packages]

The flags are documented in the package comments. See links to
documentation above.

The packages in this module can be imported and used from external
analysis driver programs. See <https://golang.org/x/tools/go/analysis>
for details. The analysis driver program may want to make
available to users the flag sets defined by both the
enumerated analyzer and the exhaustive analyzer.

# Examples

Given the expression switch

```
package a

type vcs int

const (
	bazaar vcs = iota
	fossil
	git
	mercurial
	subversion
	darcs
)

func f1(v vcs) {
	switch v {
	case bazaar:
	case fossil:
	case git:
	case subversion:
	}
}
```

the analysis produces the following diagnostic

	$ exhaustive
	a.go:15:2: missing cases in expression: mercurial, darcs

Though it is not so in the example above, in general the
enumerated type declarations and the switch statement can be in
different packages.

Given the type switch

```
package b

import "go/types"

func f2(t types.Type) bool {
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
```

the diagnostic (as of go1.27) is

	$ exhaustive -check=typeswitch
	b.go:6:2: missing cases in type switch: *types.Basic, *types.TypeParam
