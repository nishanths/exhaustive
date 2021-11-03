# exhaustive [![Godoc][godoc-status]][godoc]

The `exhaustive` package and the related command line program (found in
`cmd/exhaustive`) can be used to check exhaustiveness of enum switch
statements in Go code. An enum switch statement is exhaustive if all of
the enum's members are listed in the switch statement's cases.

```
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
```

See [pkg.go.dev](https://pkg.go.dev/github.com/nishanths/exhaustive#section-documentation)
for the flags, the definition of enum, and the definition of exhaustiveness
used by this package.

## Known issues

The package may not correctly handle enums that are [type
aliases][4]. See issue [#13][5].

## Example

Given the enum

```diff
package token

type Token int

const (
	Add Token = iota
	Subtract
	Multiply
+	Quotient
+	Remainder
)
```

and the switch statement

```
package calc

import "token"

func processToken(t token.Token) {
	switch t {
	case token.Add:
		...
	case token.Subtract:
		...
	case token.Multiply:
		...
	}
}
```

running `exhaustive` on the `calc` package

```
exhaustive ./calc/...
```

will print

```
calc.go:6:2: missing cases in switch of type token.Token: Quotient, Remainder
```

[godoc]: https://godoc.org/github.com/nishanths/exhaustive
[godoc-status]: https://godoc.org/github.com/nishanths/exhaustive?status.svg
[4]: https://go.googlesource.com/proposal/+/master/design/18130-type-alias.md
[5]: https://github.com/nishanths/exhaustive/issues/13
