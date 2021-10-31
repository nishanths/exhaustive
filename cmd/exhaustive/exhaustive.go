// Command exhaustive checks exhaustiveness of enum switch statements.
//
// For documentation see https://godoc.org/github.com/nishanths/exhaustive.
package main

import (
	"github.com/nishanths/exhaustive"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(exhaustive.Analyzer) }
