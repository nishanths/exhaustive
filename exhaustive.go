package exhaustive

import (
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	fCheckMaps       bool
	fDefaultSuffices bool
)

func init() {
	Analyzer.Flags.BoolVar(&fCheckMaps, "maps", false, "check key exhaustiveness of map literals of enum key type, in addition to checking switch statements")
	Analyzer.Flags.BoolVar(&fDefaultSuffices, "default-means-exhaustive", false, "switch statements are considered exhaustive if a 'default' case is present")
}

var Analyzer = &analysis.Analyzer{
	Name:      "exhaustive",
	Doc:       "check for any non-exhaustive enum switch statements",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&enumsFact{}},
}

const IgnoreDirective = "//exhaustive:ignore"

func containsIgnoreDirective(comments []*ast.Comment) bool {
	for _, c := range comments {
		if strings.HasPrefix(c.Text, IgnoreDirective) {
			return true
		}
	}
	return false
}

type enumsFact struct {
	entries enums
}

var _ analysis.Fact = (*enumsFact)(nil)

func (e *enumsFact) AFact() {}

func (e *enumsFact) String() string {
	// sort for stability (required for testing)
	var sortedKeys []*types.Named
	for k := range e.entries {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i].Obj().Name() < sortedKeys[j].Obj().Name()
	})

	var buf strings.Builder
	for i, k := range sortedKeys {
		v := e.entries[k]
		buf.WriteString(k.Obj().Name())
		buf.WriteString(":")
		for j, vv := range v {
			buf.WriteString(vv.Name())
			// add comma separator between each enum member in an enum type
			if j != len(v)-1 {
				buf.WriteString(",")
			}
		}
		// add semicolon separator between each enum type
		if i != len(sortedKeys)-1 {
			buf.WriteString("; ")
		}
	}
	return buf.String()
}

func run(pass *analysis.Pass) (interface{}, error) {
	e := findEnums(pass)
	pass.ExportPackageFact(&enumsFact{entries: e})

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	comments := make(map[*ast.File]ast.CommentMap) // CommentMap per package file, lazily populated

	checkSwitchStatements(pass, inspect, comments)
	if fCheckMaps {
		checkMapLiterals(pass, inspect, comments)
	}
	return nil, nil
}

func enumTypeName(e *types.Named, samePkg bool) string {
	if samePkg {
		return e.Obj().Name()
	}
	return e.Obj().Pkg().Name() + "." + e.Obj().Name()
}
