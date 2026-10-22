package main

import (
	"log"
	"os"
	"slices"
	"strings"

	"github.com/nishanths/exhaustive/internal/driver"
	"github.com/nishanths/exhaustive/internal/passes/finder"
	"github.com/nishanths/exhaustive/passes/enumerated"
	"github.com/nishanths/exhaustive/passes/exhaustive"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

const progname = "exhaustive"

func main() {
	log.SetFlags(0)
	log.SetPrefix(progname + ": ")

	fs := driver.ParseFlags(progname, enumerated.Analyzer, exhaustive.Analyzer)
	check := exhaustive.Analyzer.Flags.Lookup("check").Value.String()
	editRequires(check)
	os.Exit(driver.Run(fs.Args(), packages.LoadAllSyntax, exhaustive.Analyzer))
}

func editRequires(check string) {
	// see comment on the exhaustive.Analyzer.Requires field
	// for details.
	remove := map[*analysis.Analyzer]bool{
		enumerated.Analyzer: true,
		finder.Analyzer:     true,
	}

	for _, v := range strings.Split(check, ",") {
		switch v {
		case "switch", "mapliteral":
			remove[enumerated.Analyzer] = false
		case "typeswitch":
			remove[finder.Analyzer] = false
		}
	}

	exhaustive.Analyzer.Requires = slices.DeleteFunc(exhaustive.Analyzer.Requires, func(a *analysis.Analyzer) bool {
		v, ok := remove[a]
		if !ok {
			// unknown/unhandled analyzer.
			// return false (i.e. do not remove).
			return false
		}
		return v
	})
}
