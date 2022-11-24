# exhaustive [![Godoc][godoc-svg]][repo]

Checks exhaustiveness of enum switch statements in Go source code.

The repository consists of an importable Go package and a command line
program. The package provides an `analysis.Analyzer` value that follows
the guidelines in the [`golang.org/x/tools/go/analysis`][xanalysis]
package. This should make it possible to integrate exhaustive with your
own analysis driver programs.

To install the command line program, run:

```
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
```

For documentation on the command's flags, definition of enums, and
definition of exhaustiveness, see [pkg.go.dev][godoc]. For a changelog,
see [CHANGELOG][changelog] in the wiki.

The program may additionally be configured to check for exhaustiveness
of map literals with enum key types. See examples below.

## Bugs

`exhaustive` does not report missing cases for a switch statement that
switch on a type-parameterized type. For details see [this
issue][issue-typeparam].

## Examples

### Switch statement

Given the enum

```go
package token

type Token int

const (
	Add Token = iota
	Subtract
	Multiply
	Quotient
	Remainder
)
```

and code that switches on the enum

```go
package calc

import "token"

func doSomething(t token.Token) {
	switch t {
	case token.Add:
	case token.Subtract:
	case token.Multiply:
	default:
	}
}

var tokenNames = map[token.Token]string{
	token.Add:      "add",
	token.Subtract: "subtract",
	token.Multiply: "multiply",
}
```

running exhaustive with default options will print:

```
% exhaustive path/to/pkg/calc
calc.go:6:2: missing cases in switch of type token.Token: Quotient, Remainder
%
```

### Map literal

To additionally check exhaustiveness of map literals, use the `-check`
flag.

```
% exhaustive -check=switch,map path/to/pkg/calc
calc.go:6:2: missing cases in switch of type token.Token: Quotient, Remainder
calc.go:14:18: missing keys in map of key type token.Token: Quotient, Remainder
%
```

## Contributing

Issues and pull requests are welcome. Before making a substantial
change, please discuss it in an issue.

[repo]: https://pkg.go.dev/github.com/nishanths/exhaustive
[godoc-svg]: https://pkg.go.dev/github.com/nishanths/exhaustive?status.svg
[godoc]: https://pkg.go.dev/github.com/nishanths/exhaustive
[xanalysis]: https://pkg.go.dev/golang.org/x/tools/go/analysis
[changelog]: https://github.com/nishanths/exhaustive/wiki/CHANGELOG
[issue-typeparam]: https://github.com/nishanths/exhaustive/issues/31
