package driver

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/packages"
)

const (
	includeTests = true
	contextLines = -1
)

func ParseFlags(progname, synopsis string, desc []byte, as ...*analysis.Analyzer) *flag.FlagSet {
	fs := flag.NewFlagSet(progname, flag.ContinueOnError)

	var help bool
	fs.BoolVar(&help, "h", help, "print help information")
	fs.BoolVar(&help, "help", help, "print help information")

	// add flags from analyzers
	for _, a := range as {
		a.Flags.VisitAll(func(f *flag.Flag) {
			if fs.Lookup(f.Name) != nil {
				log.Fatalf("internal error: conflicting flag -%s", f.Name)
			}
			fs.Var(f.Value, f.Name, f.Usage)
		})
	}

	fs.Usage = func() {
		printFlagUsage(os.Stderr, synopsis)
	}

	err := fs.Parse(os.Args[1:])
	if err != nil {
		// Note: The usage message would have already been printed by
		// package flag.
		os.Exit(2)
	}
	if help {
		printFlagUsage(os.Stdout, synopsis)
		_, _ = fmt.Fprintln(os.Stdout)
		_, _ = fmt.Fprintf(os.Stdout, "%s", desc)
		_, _ = fmt.Fprintln(os.Stdout)
		os.Exit(0)
	}
	return fs
}

func Run(args []string, loadMode packages.LoadMode, as ...*analysis.Analyzer) (exitStatus int) {
	exitAtLeast := func(s int) {
		if s > exitStatus {
			exitStatus = s
		}
	}

	cfg := &packages.Config{
		Mode:  loadMode | packages.NeedModule,
		Tests: includeTests,
	}
	pkgs, err := packages.Load(cfg, args...)
	if err != nil {
		log.Printf("error loading packages: %s", err)
		exitAtLeast(1)
		return
	}

	if len(pkgs) == 0 {
		log.Println("no packages matched")
		exitAtLeast(1)
		return
	}

	if packages.PrintErrors(pkgs) > 0 {
		exitAtLeast(1)
		// Do not return, the analysis can still proceed.
	}

	graph, err := checker.Analyze(as, pkgs, nil)
	if err != nil {
		log.Println(err)
		exitAtLeast(1)
		return
	}

	if err := graph.PrintText(os.Stderr, contextLines); err != nil {
		log.Println(err)
		exitAtLeast(1)
		return
	}

	var errors, diags int
	for act := range graph.All() {
		if act.Err != nil {
			errors++
			continue
		}
		if act.IsRoot {
			diags += len(act.Diagnostics)
		}
	}
	// Note: These exit status values match those used by
	// analysis/singlechecker.
	if errors > 0 {
		exitAtLeast(1)
	}
	if diags > 0 {
		exitAtLeast(3)
	}
	return
}

func printFlagUsage(w io.Writer, synopsis string) {
	indent := strings.Repeat(" ", len("usage: "))
	first := true

	for line := range strings.Lines(synopsis) {
		switch {
		case first:
			_, _ = fmt.Fprint(w, "usage: ")
			_, _ = fmt.Fprint(w, line)
			first = false
		case line == "\n":
			_, _ = fmt.Fprint(w, line)
		default:
			_, _ = fmt.Fprint(w, indent+line)
		}
	}
}
