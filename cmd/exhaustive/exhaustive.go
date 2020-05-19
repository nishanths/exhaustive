// Command exhaustive ... TODO ...
package main

import (
	"github.com/nishanths/exhaustive"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(exhaustive.Analyzer) }
