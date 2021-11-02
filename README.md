# exhaustive [![Godoc][godoc-status]][godoc] [![Build Status][build-status]][build]

The `exhaustive` package and the related command line program
(`cmd/exhaustive`) can be used to check exhaustiveness of enum switch
statements in Go code.

An enum switch statement is exhaustive if it has cases for each of the
enum's members. See Godoc for the definition of enum used by this
package.

## Docs

https://godoc.org/github.com/nishanths/exhaustive

## Known issues

The package may not correctly handle enums that are [type
aliases][4]. See issue [#13][5].

## Install command line program

Install latest tagged release:

```
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
```

Install latest `master`:

```
go install github.com/nishanths/exhaustive/cmd/exhaustive@master
```

## Example

Given this enum:

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

And code elsewhere that switches on the enum:

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

Running the `exhaustive` command on the `calc` package will print:

```
calc.go:6:2: missing cases in switch of type token.Token: Quotient, Remainder
```

Enums can also be defined using explicit constant values instead of `iota`.

## Integrate with analyzer driver programs

The `exhaustive` package provides an `Analyzer` that follows the
guidelines described in the [go/analysis][3] package; this should make
it possible to integrate `exhaustive` into analysis driver
programs.

## License

BSD 2-Clause

[godoc]: https://godoc.org/github.com/nishanths/exhaustive
[godoc-status]: https://godoc.org/github.com/nishanths/exhaustive?status.svg
[build]: https://travis-ci.org/nishanths/exhaustive
[build-status]: https://travis-ci.org/nishanths/exhaustive.svg?branch=master
[3]: https://godoc.org/golang.org/x/tools/go/analysis
[4]: https://go.googlesource.com/proposal/+/master/design/18130-type-alias.md
[5]: https://github.com/nishanths/exhaustive/issues/13
