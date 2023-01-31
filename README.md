# exhaustive [![Godoc][godoc-svg]][godoc]

Package exhaustive defines an analyzer that checks exhaustiveness of switch
statements of enum-like constants in Go source code. The analyzer can
be configured to additionally check exhaustiveness of keys in map literals
whose key type is enum-like.

For flags, the definition of enum, and the definition of exhaustiveness used
by this package, see [docs][godoc-doc]. For a changelog, see
[CHANGELOG][changelog] in the GitHub wiki.

## Usage

Command:

```
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest

exhaustive [flags] [packages]
```

Package:

```
go get github.com/nishanths/exhaustive

import "github.com/nishanths/exhaustive"
```

The `exhaustive.Analyzer` variable follows guidelines in the
[`golang.org/x/tools/go/analysis`][xanalysis] package. This should make it
possible to integrate `exhaustive` with your own analysis driver program.

## Example

Given an enum:

```go
package token // import "example.org/token"

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

import "example.org/token"

func x(t token.Token) {
	switch t {
	case token.Add:
	case token.Subtract:
	case token.Multiply:
	default:
	}
}
```

running `exhaustive` with default options will produce:

```
calc.go:6:2: missing cases in switch of type token.Token: token.Quotient, token.Remainder
```

Specify flag `-check=switch,map` to additionally check exhaustiveness of map
literal keys. For example:

```go
var m = map[token.Token]rune{
	token.Add:       '+',
	token.Multiply:  '*',
	token.Quotient:  '/',
	token.Remainder: '%',
}
```

```
calc.go:14:9: missing keys in map of key type token.Token: token.Subtract
```

## Contributing

Issues and changes are welcome. Please discuss substantial changes
in an issue first.

[godoc]: https://pkg.go.dev/github.com/nishanths/exhaustive
[godoc-svg]: https://pkg.go.dev/badge/github.com/nishanths/exhaustive.svg
[godoc-doc]: https://pkg.go.dev/github.com/nishanths/exhaustive#section-documentation
[godoc-flags]: https://pkg.go.dev/github.com/nishanths/exhaustive#hdr-Flags
[xanalysis]: https://pkg.go.dev/golang.org/x/tools/go/analysis
[changelog]: https://github.com/nishanths/exhaustive/wiki/CHANGELOG
[issue-typeparam]: https://github.com/nishanths/exhaustive/issues/31
