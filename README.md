## exhaustive [![Godoc][2]][1]

Check exhaustiveness of enum switch statements in Go source code.

```
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
```

For docs, see package comment at [pkg.go.dev][6]. It describes the
flags, the definition of enum, and the definition of exhaustiveness used
by this package.

For the changelog, see [CHANGELOG][changelog] in the wiki.

The package provides an `Analyzer` that follows the guidelines in the
[`go/analysis`][3] package; this should make it possible to integrate
exhaustive with your own analysis driver program.

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
	default: ...
	}
}
```

running exhaustive

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
