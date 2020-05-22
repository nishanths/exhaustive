The `exhaustive` package and command line program can be used to detect
enum switch statements in Go code that are not exhaustive.

An enum switch statment is exhaustive if it has cases for each of the enum's members.
Exhaustive switches are useful for ensuring at compile time that all enum cases are
properly handled. They can be useful, for instance, to draw attention to switch
statements that need to be updated when a new member is added to an existing enum.

For the purpose of this program, an enum type is a package-level named integer, float, or
string type. An enum type has associated with one or more enum members that are variables
of the enum type.

## Install

```
go get github.com/nishanths/exhaustive/...
```

## Docs

See Godoc: https://godoc.org/github.com/nishanths/exhaustive.

The `exhaustive` package provides a valid "pass", similar to the passes defined in the [`go/analysis`](http://godoc.org/golang.org/x/tools/go/analysis) package. This makes it easy to integrate the package into an existing analysis driver program.

## Example

Running the `exhaustive` command on the following code:

```go
package environment

// Biome is an enum type with 3 members.
type Biome int

const (
	Tundra Biome = iota
	Savanna
	Desert
)
```

```go
package foo

func BiomeDescription(b environment.Biome) {
	switch b {
	case Tundra:
		println("The tundra is extremely cold")
	case Desert:
		println("Deserts are arid")
	}
}
```

would print:

```
missing cases in switch of type environment.Biome: Savanna
```

## Usage

```
exhaustive [-flags] [packages...]
```

The relevant flags are:

```
-default-signifies-exhaustive
    switch statements are considered exhaustive if a 'default' case is present
-fix
    apply all suggested fixes
```
