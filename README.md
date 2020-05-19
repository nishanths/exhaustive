The `exhaustive` command can be used to ensure that _enum_ `switch` statements in Go code are
exhaustive. Optionally, it can also ensure that map keys listed in map literals of an
enum key type are exhaustive.

## Example

## Install

```
go get github.com/nishanths/exhaustive/...
```

## Documentation

See [godoc](https://godoc.org/github.com/nishanths/exhaustive/cmd/exhaustive) for usage and more documentation.

For the purpose of this program, the members of an enum are of the set of package-level constant
values for a named type.

```go
// Biome is an enum type with three members: Tundra, Savanna, Desert.
type Biome string

const (
	Tundra  Biome = "tundra"
	Savanna Biome = "savanna"
)

const Desert Biome = "desert"
```
