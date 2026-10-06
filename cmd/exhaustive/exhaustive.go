package main

import (
	"flag"
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

	driver.ParseFlags(progname, enumerated.Analyzer, exhaustive.Analyzer)
	editRequires(exhaustive.Analyzer.Flags.Lookup("check").Value.String())
	os.Exit(driver.Run(flag.Args(), packages.LoadAllSyntax, exhaustive.Analyzer))
}

func editRequires(check string) {
	// see comment on the exhaustive.Analyzer.Requires field
	// for details.
	need := map[string]bool{
		"enumerated": false,
		"finder":     false,
	}

	for _, v := range strings.Split(check, ",") {
		switch v {
		case "switch", "mapliteral":
			need["enumerated"] = true
		case "typeswitch":
			need["finder"] = true
		}
	}

	exhaustive.Analyzer.Requires = slices.DeleteFunc(exhaustive.Analyzer.Requires, func(a *analysis.Analyzer) bool {
		switch {
		case a == enumerated.Analyzer && !need["enumerated"]:
			return true
		case a == finder.Analyzer && !need["finder"]:
			return true
		default:
			return false
		}
	})
}
