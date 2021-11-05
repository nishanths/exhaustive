## exhaustive [![Godoc][2]][1]

The `exhaustive` package and the related command line program (found in
`cmd/exhaustive`) can be used to check exhaustiveness of enum switch
statements in Go code.

Install the command:

```
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
```

See [pkg.go.dev][6] for the flags, the definition of enum, and the
definition of exhaustiveness used by this package.

For changelog, see [CHANGELOG][changelog] in the wiki.

The `exhaustive` package provides an `Analzyer` that follows the
guidelines in the [`go/analysis`][3] package; this should make
it possible to integrate with external analysis driver programs.

### Known issues

The package may not correctly handle enum types that are [type
aliases][4]. See issue [#13][5].

### Example

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
	case token.Add: ...
	case token.Subtract: ...
	case token.Multiply: ...
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

[1]: https://godoc.org/github.com/nishanths/exhaustive
[2]: https://godoc.org/github.com/nishanths/exhaustive?status.svg
[3]: https://pkg.go.dev/golang.org/x/tools/go/analysis
[4]: https://go.googlesource.com/proposal/+/master/design/18130-type-alias.md
[5]: https://github.com/nishanths/exhaustive/issues/13
[6]: https://pkg.go.dev/github.com/nishanths/exhaustive#section-documentation
[changelog]: https://github.com/nishanths/exhaustive/wiki/CHANGELOG
