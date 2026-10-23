package main

import (
	"log"
	"os"

	"github.com/nishanths/exhaustive/internal/driver"
	"github.com/nishanths/exhaustive/passes/enumerated"
	"github.com/nishanths/exhaustive/passes/exhaustive"

	"golang.org/x/tools/go/packages"
)

const progname = "exhaustive"

func main() {
	log.SetFlags(0)
	log.SetPrefix(progname + ": ")

	fs := driver.ParseFlags(progname, enumerated.Analyzer, exhaustive.Analyzer)
	if err := editRequires(); err != nil {
		log.Fatal(err)
	}
	os.Exit(driver.Run(fs.Args(), packages.LoadAllSyntax, exhaustive.Analyzer))
}

func editRequires() error {
	check := exhaustive.Analyzer.Flags.Lookup("check").Value.String()
	r, err := exhaustive.Requires(check)
	if err != nil {
		return err
	}
	exhaustive.Analyzer.Requires = r
	return nil
}
