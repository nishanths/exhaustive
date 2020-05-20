package exhaustive

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const IgnoreDirective = "//exhaustive:ignore"

func containsIgnoreDirective(comments []*ast.Comment) bool {
	for _, c := range comments {
		return strings.HasPrefix(c.Text, IgnoreDirective)
	}
	return false
}

var Analyzer = &analysis.Analyzer{
	Name:      "exhaustive",
	Doc:       "check for non-exhaustive enum switch statements",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&enumsFact{}},
}

var (
	fCheckMaps       bool
	fDefaultSuffices bool
)

func init() {
	Analyzer.Flags.BoolVar(&fCheckMaps, "maps", false, "check map literals of enum key type, in addition to switch statements")
	Analyzer.Flags.BoolVar(&fDefaultSuffices, "default-is-exhaustive", false, "switch statements are considered exhaustive as long as 'default' case exists")
}

type enumsFact struct {
	entries enums
}

var _ analysis.Fact = (*enumsFact)(nil)

func (e *enumsFact) AFact() {}

func run(pass *analysis.Pass) (interface{}, error) {
	e := gatherEnums(pass)
	pass.ExportPackageFact(&enumsFact{entries: e})

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	comments := make(map[*ast.File]ast.CommentMap) // CommentMap per package file, lazily populated

	checkSwitchStatements(pass, inspect, comments)
	checkMapLiterals(pass, inspect, comments)

	return nil, nil
}

func assert(cond bool, format string, args ...interface{}) {
	if !cond {
		panicf(format, args...)
	}
}

func panicf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

func removeParens(e ast.Expr) ast.Expr {
	for {
		parenExpr, ok := e.(*ast.ParenExpr)
		if !ok {
			break
		}
		e = parenExpr.X
	}
	return e
}
