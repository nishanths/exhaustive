# exhaustive [![Godoc][godoc-svg]][godoc]

Package exhaustive defines an analyzer that checks exhaustiveness of switch
statements of enum-like constants in Go source code.

For flags, the definition of enum, and the definition of exhaustiveness used
by this package, see [pkg.go.dev][godoc-doc]. For a changelog, see
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
possible to integrate `exhaustive` in your own analysis driver program.

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

func g(t token.Token) {
	switch t {
	case token.Add:
	case token.Subtract:
	case token.Multiply:
	default:
	}
}
```

running `exhaustive` with default options will print:

```
calc.go:6:2: missing cases in switch of type token.Token: token.Quotient, token.Remainder
```

Specify flag `-check=switch,map` to additionally check exhaustiveness of map
literal keys. For example, given:

```go
var m = map[token.Token]string{
	token.Add:       "add",
	token.Multiply:  "multiply",
	token.Quotient:  "quotient",
	token.Remainder: "remainder",
}
```

and `exhaustive` will print:

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
