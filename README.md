## exhaustive [![Godoc][godoc-status]][godoc] [![Build Status][build-status]][build]

The `exhaustive` package and command line program can be used to detect enum
switch statements that are not exhaustive.

An enum switch statement is exhaustive if it has cases for each of the enum's
members. See godoc for the definition of enum used by the program.

The `exhaustive` package provides an `Analyzer` type that follows the
guidelines described in the [go/analysis][3] package; this makes it
possible to also integrate `exhaustive` into analysis driver
programs.

### Install

Install the command line program (with Go 1.16 or higher):

```
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
```

Install the package:

```
go get github.com/nishanths/exhaustive
```

### Known issues

The program may not correctly handle enum types that are [type
aliases][4]. See [issue #13][5].

### Docs

https://godoc.org/github.com/nishanths/exhaustive

### Example

Given this enum type:

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

Running the `exhaustive` command will print:

```
calc.go:6:2: missing cases in switch of type token.Token: Quotient, Remainder
```

Enums can also be defined using explicit constant values instead of `iota`.

### License

BSD 2-Clause

[godoc]: https://godoc.org/github.com/nishanths/exhaustive
[godoc-status]: https://godoc.org/github.com/nishanths/exhaustive?status.svg
[build]: https://travis-ci.org/nishanths/exhaustive
[build-status]: https://travis-ci.org/nishanths/exhaustive.svg?branch=master
[3]: https://godoc.org/golang.org/x/tools/go/analysis
[4]: https://go.googlesource.com/proposal/+/master/design/18130-type-alias.md
[5]: https://github.com/nishanths/exhaustive/issues/13
