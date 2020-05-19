# exhaustive

`exhaustive` can be used to ensure that _enum_ `switch` statements in Go code are
exhaustive. Optionally, it can also ensure that map keys in map literals of an
enum key type are exhaustive.

For the purpose of this program, the members of an enum are of the set of constant
values for a named type.

```go
package foo

// The Environment enum has three members: Prod, Staging, Dev.
type Environment string

const (
	Prod    Environment = "production"
	Staging Environment = "staging"
)

const Dev Environment = "development"
```

## Install

```
go get github.com/nishanths/exhaustive/...
```
