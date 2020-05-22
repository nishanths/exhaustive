# exhaustive

[![Godoc](https://godoc.org/github.com/nishanths/exhaustive?status.svg)](https://godoc.org/github.com/nishanths/exhaustive)

The `exhaustive` package and command line program can be used to find
enum switch statements that are not exhaustive.

An enum switch statment is exhaustive if it has cases for each of the enum's members.

## Install

```
go get github.com/nishanths/exhaustive/...
```

## Docs

See Godoc: https://godoc.org/github.com/nishanths/exhaustive

The `exhaustive` package provides a valid "pass", similar to the passes defined in the [`go/analysis`](http://godoc.org/golang.org/x/tools/go/analysis) package. This makes it easy to integrate the package into an existing analysis driver program.

## Example

Running the `exhaustive` command on the following code:

```go
package ecosystem

// Biome is an enum type with 3 members.
type Biome int

const (
	Tundra Biome = iota
	Savanna
	Desert
)
```
```go
package pkg

func BiomeDescription(b ecosystem.Biome) string {
	switch b {
	case Tundra:
		return "the tundra is extremely cold"
	case Desert:
		return "deserts are arid"
	}
}
```

would print:

```
missing cases in switch of type ecosystem.Biome: Savanna
```

## Usage

```
Usage: exhaustive [-flags] [packages...]

Flags:
  -default-signifies-exhaustive
    	switch statements are considered exhaustive if a 'default' case is present
  -fix
    	apply all suggested fixes
```

## License

BSD 2-Clause
