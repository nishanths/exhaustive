# exhaustive [![Godoc][godoc-svg]][repo]

Package exhaustive defines an analyzer that checks exhaustiveness of switch
statements of enum-like constants in Go source code. The analyzer can be
configured to additionally check exhaustiveness of map literals whose key type
is enum-like.

For documentation on the flags, the definition of enum, and the definition of
exhaustiveness, see [pkg.go.dev][godoc-doc]. For a changelog, see
[CHANGELOG][changelog] in the GitHub wiki.

The exported `analysis.Analyzer` uses the
[`golang.org/x/tools/go/analysis`][xanalysis] API. This should make it
possible to integrate `exhaustive` in your own analysis driver program.

## Install

```
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
```

## Usage

```
exhaustive [flags] [packages]
```

## Example

Given the enum:

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

and code that switches on the enum:

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

running `exhaustive` with default options will print:

```
$ exhaustive
calc.go:6:2: missing cases in switch of type token.Token: token.Quotient, token.Remainder
$
```

To additionally check exhaustiveness of map literal keys, use
`-check=switch,map`:

```
$ exhaustive -check=switch,map
calc.go:6:2: missing cases in switch of type token.Token: token.Quotient, token.Remainder
calc.go:14:9: missing keys in map of key type token.Token: token.Quotient, token.Remainder
$
```

## Contributing

Issues and changes are welcome. Please discuss substantial changes
in an issue first.

[repo]: https://pkg.go.dev/github.com/nishanths/exhaustive
[godoc-svg]: https://pkg.go.dev/badge/github.com/nishanths/exhaustive.svg
[godoc-doc]: https://pkg.go.dev/github.com/nishanths/exhaustive#section-documentation
[xanalysis]: https://pkg.go.dev/golang.org/x/tools/go/analysis
[changelog]: https://github.com/nishanths/exhaustive/wiki/CHANGELOG
[issue-typeparam]: https://github.com/nishanths/exhaustive/issues/31
