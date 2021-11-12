// Command exhaustive checks exhaustiveness of enum switch statements.
//
// Usage
//
// The command line usage is:
//
//    exhaustive [flags] [packages]
//
// The program checks exhaustiveness of enum switch statements found in the
// specified packages. The enums required for the analysis don't necessarily
// have to be declared in the specified packages.
//
// For more about specifying packages, see 'go help packages'.
//
// For help, run 'exhaustive -h'.
//
// For more documentation, see https://godocs.io/github.com/nishanths/exhaustive.
package main

import (
	"github.com/nishanths/exhaustive"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(exhaustive.Analyzer) }
