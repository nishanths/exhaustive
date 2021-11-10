## exhaustive [![Godoc][2]][1]

Check exhaustiveness of enum switch statements in Go source code.

```
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
```

For docs on the flags, the definition of enum, and the definition of
exhaustiveness, see [pkg.go.dev][6].

For the changelog, see [CHANGELOG][changelog] in the wiki.

The package provides an `Analyzer` that follows the guidelines in the
[`go/analysis`][3] package; this should make it possible to integrate
exhaustive with your own analysis driver program.

## Example

Given the enum

```go
package env

type Environment string

const (
	Production Environment = "production"
	Staging    Environment = "staging"
	Dev        Environment = "dev"
)

func Current() Environment { /* ... */ }
```

and the switch statement

```go
package app

import "example/pkg/env"

func f() {
	switch env.Current() {
	case env.Production:
	case env.Dev:
	default:
	}
}
```

running exhaustive will print

```
app.go:6:2: missing cases in switch of type env.Environment: Staging
```

## Contributing

Issues and pull requests are welcome. Before making a substantial
change, please discuss it in an issue.

[1]: https://godoc.org/github.com/nishanths/exhaustive
[2]: https://godoc.org/github.com/nishanths/exhaustive?status.svg
[3]: https://pkg.go.dev/golang.org/x/tools/go/analysis
[6]: https://pkg.go.dev/github.com/nishanths/exhaustive#section-documentation
[changelog]: https://github.com/nishanths/exhaustive/wiki/CHANGELOG
