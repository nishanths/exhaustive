package driver

import (
	"cmp"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/packages"
)

const (
	includeTests = true
	contextLines = -1
)

func ParseFlags(progname string, as ...*analysis.Analyzer) {
	for _, a := range as {
		a.Flags.VisitAll(func(f *flag.Flag) {
			if flag.Lookup(f.Name) != nil {
				log.Fatalf("internal error: conflicting flag -%s", f.Name)
			}
			flag.Var(f.Value, f.Name, f.Usage)
		})
	}
	flag.Usage = func() {
		syn := synopsis(progname, len("usage: "), false, slicemap(as, func(a *analysis.Analyzer) flag.FlagSet { return a.Flags }))
		fmt.Fprintf(os.Stderr, "usage: %s\n", syn)
	}
	flag.Parse()
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
		log.Println("matched no packages")
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

func synopsis(progname string, prefixlen int, includeHelpInvoc bool, sets []flag.FlagSet) string {
	// Sort flags by name, boolean flags first.
	var flags []*flag.Flag
	for _, fs := range sets {
		fs.VisitAll(func(f *flag.Flag) {
			flags = append(flags, f)
		})
	}
	slices.SortStableFunc(flags, func(a, b *flag.Flag) int {
		type boolflag interface{ IsBoolFlag() bool }
		_, aa := a.Value.(boolflag)
		_, bb := b.Value.(boolflag)
		switch {
		case aa && !bb:
			return -1
		case !aa && bb:
			return 1
		default:
			return cmp.Compare(a.Name, b.Name)
		}
	})

	const maxwidth = 80
	indentlen := 8
	if includeHelpInvoc {
		// Multiple invocations need to be printed. The
		// lines must align with program name.
		indentlen = prefixlen + len(progname) + 1
	}
	lastarg := "[packages]" // last argument must never be alone on its own line
	var width int
	var b strings.Builder

	newinvoc := func(first bool) {
		if !first {
			fmt.Fprintf(&b, "\n")
			fmt.Fprintf(&b, "%s", strings.Repeat(" ", prefixlen))
		}
		width = 0
		width += prefixlen
		fmt.Fprintf(&b, "%s", progname)
		width += len(progname)
	}

	printarg := func(arg string, prelast bool) {
		w := width + 1 + len(arg)
		if prelast {
			w = width + 1 + len(arg) + 1 + len(lastarg)
		}
		if w <= maxwidth {
			fmt.Fprintf(&b, " %s", arg)
			width += 1 + len(arg)
			return
		}
		fmt.Fprintf(&b, "\n")
		fmt.Fprintf(&b, "%s%s", strings.Repeat(" ", indentlen), arg)
		width = indentlen + len(arg)
	}

	// Regular invocation.
	newinvoc(true)
	for i, f := range flags {
		if v, _ := flag.UnquoteUsage(f); len(v) != 0 {
			printarg(fmt.Sprintf("[-%s %v]", f.Name, v), i == len(flags)-1)
		} else {
			printarg(fmt.Sprintf("[-%s]", f.Name), i == len(flags)-1)
		}
	}
	printarg(lastarg, false)

	// Help invocation.
	if includeHelpInvoc {
		newinvoc(false)
		printarg("-h | -help", false)
	}

	return b.String()
}

func slicemap[S ~[]E, E, F any](s S, fn func(E) F) []F {
	var ret []F
	if len(s) > 0 {
		ret = make([]F, len(s))
	}
	for i := range s {
		ret[i] = fn(s[i])
	}
	return ret
}
