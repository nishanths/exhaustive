// Command exhaustive can find enum switch statements that are non-exhaustive.
//
// For details see the doc comment for the exhaustive package that this command
// is based upon at https://godoc.org/github.com/nishanths/exhaustive.
package main

import (
	"github.com/nishanths/exhaustive"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(exhaustive.Analyzer) }
