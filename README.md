# exhaustive [![Godoc][godoc-svg]][repo]

Checks exhaustiveness of enum switch statements in Go source code.

```
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
```

The program can additionally be configured to check for exhaustiveness
of map literals with enum key types. See examples below.

For documentation on flags, the definition of enum, and the definition
of exhaustiveness, see [pkg.go.dev][godoc-doc]. For a changelog, see
[CHANGELOG][changelog] in the wiki.

The package provides an `analysis.Analyzer` value that follows the
guidelines in the [`golang.org/x/tools/go/analysis`][xanalysis] package.
This should make it possible to integrate `exhaustive` with your own
analysis driver programs.

## Examples

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

func f(t token.Token) {
	switch t {
	case token.Add:
	case token.Subtract:
	case token.Multiply:
	default:
	}
}

var m = map[token.Token]string{
	token.Add:      "add",
	token.Subtract: "subtract",
	token.Multiply: "multiply",
}
```

running `exhaustive` with default options will print

```
% exhaustive path/to/pkg/calc
calc.go:6:2: missing cases in switch of type token.Token: Quotient, Remainder
%
```

To additionally check exhaustiveness of map literals, use
`-check=switch,map`.

```
% exhaustive -check=switch,map path/to/pkg/calc
calc.go:6:2: missing cases in switch of type token.Token: Quotient, Remainder
calc.go:14:9: missing keys in map of key type token.Token: Quotient, Remainder
%
```

## Contributing

Issues and pull requests are welcome. Before making a substantial
change please discuss it in an issue.

[repo]: https://pkg.go.dev/github.com/nishanths/exhaustive
[godoc-svg]: https://pkg.go.dev/badge/github.com/nishanths/exhaustive.svg
[godoc-doc]: https://pkg.go.dev/github.com/nishanths/exhaustive#section-documentation
[xanalysis]: https://pkg.go.dev/golang.org/x/tools/go/analysis
[changelog]: https://github.com/nishanths/exhaustive/wiki/CHANGELOG
[issue-typeparam]: https://github.com/nishanths/exhaustive/issues/31
