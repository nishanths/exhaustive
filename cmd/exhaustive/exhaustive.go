package main

import (
	"flag"
	"log"
	"os"

	"github.com/nishanths/exhaustive/internal/driver"
	"github.com/nishanths/exhaustive/passes/enumerated"
	"github.com/nishanths/exhaustive/passes/exhaustive"

	"golang.org/x/tools/go/packages"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("exhaustive: ")
	driver.ParseFlags("exhaustive", enumerated.Analyzer, exhaustive.Analyzer)
	os.Exit(driver.Run(flag.Args(), packages.LoadAllSyntax, exhaustive.Analyzer))
}
