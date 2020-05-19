# exhaustive

The `exhaustive` command line program can be used to ensure that enum
`switch` statements in Go code are exhaustive. Optionally, it can also
ensure that the keys listed in `map` literals of an enum key-type are exhaustive.

It works only for expression switch statements, not type switch statements.

## Install

```
go get github.com/nishanths/exhaustive/...
```
