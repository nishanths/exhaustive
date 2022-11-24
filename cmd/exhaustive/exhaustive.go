// Command exhaustive is a command line interface for the exhaustive
// package at github.com/nishanths/exhaustive.
//
// # Usage
//
//	exhaustive [flags] [packages]
package main

import (
	"github.com/nishanths/exhaustive"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(exhaustive.Analyzer) }
