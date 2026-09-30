The exhaustive static analysis checks that expression switch
statements, in which the type of the switch expression is an
enumerated type, are exhaustive.

The analysis is implemented in multiple passes; see the packages
in the 'passes/\*' directories.
Package enumerated find declarations of enumerated types and
enumerated constants.
Package exhaustive scans a syntax tree for eligible
expression switch statements and checks that the switch
statements are exhaustive.

Documentation:

<https://pkg.go.dev/github.com/nishanths/exhaustive/passes/enumerated>

<https://pkg.go.dev/github.com/nishanths/exhaustive/passes/exhaustive>

# Usage

The exhaustive command can be installed with 'go install'.

	go install github.com/nishanths/exhaustive/cmd/exhaustive@latest

The synopsis of the command is:

	exhaustive [-B] [-d] [-defrequire] [-e] [-g] [-i] [-p] [-check string]
	           [-constignore regexp] [-typeignore regexp] [-typeonly regexp] [packages]

The flags are documented in the package comments. See links to
documentation above.

The packages in this module can be imported and used from external
analysis driver programs. See <https://golang.org/x/tools/go/analysis>
for details. The analysis driver program may want to make
available to users the flag sets defined by both the
enumerated analyzer and the exhaustive analyzer.

# Examples

Given the following Go source code:

```
package example

type vcs int

const (
	bzr vcs = iota
	fossil
	git
	hg
	svn
	darcs
)

func f(v vcs) {
	switch v {
	case bzr:
	case fossil:
	case git:
	case svn:
	}
}
```

the exhaustive command produces the following diagnostic:

	example.go:15:2: switch not exhaustive: missing cases: hg, darcs

Though it is not so in the example, in general the enumerated
type declarations and the switch statements can be in different
packages.
