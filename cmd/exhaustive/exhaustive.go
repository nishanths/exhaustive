// Command exhaustive can find enum switch statements that are non-exhaustive.
//
// For details see the doc comment for the exhaustive package that this command
// is based upon at https://godoc.org/github.com/nishanths/exhaustive.
//
// Usage
//
//
//   Usage: exhaustive [-flags] [packages...]
//
//   Flags:
//     -default-signifies-exhaustive
//       	switch statements are considered exhaustive if a 'default' case is present
//     -fix
//      	apply all suggested fixes
//
//   Examples:
//     exhaustive code.org/proj/...
//     exhaustive -fix example.org/foo/pkg example.org/foo/bar
package main

import (
	"github.com/nishanths/exhaustive"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(exhaustive.Analyzer) }
