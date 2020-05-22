// Package exhaustive provides an analyzer that checks for enum switch statements
// that are not exhaustive. The analyzer can suggest fixes to make offending switch
// statements exhaustive.
//
// Exhaustiveness
//
// An enum switch statment is exhaustive if it has cases for each of the enum's members.
// For an enum type defined in the same package as the switch statement, both
// exported and unexported enum members must be present in order to consider
// the switch exhaustive. On the other hand, for an enum type defined
// in an external package it is sufficient for just the exported enum members
// to be present in order to consider the switch exhaustive.
//
// Definition of enum
//
// For the purpose of this program, an enum type is a package-level named integer, float, or
// string type. Such a type qualifies as an enum type only if it there exist one
// or more  package-level variables of this named type in the package. These variables
// constitute the enum's members.
//
// In the code sample below, Biome is an enum type with 3 members.
//
//   type Biome int
//
//   const (
//       Tundra Biome = iota
//       Savanna
//       Desert
//   )
//
// Fixes
//
// The analyzer can suggest fixes for a switch statement if it is not exhaustive,
// and if it does not have a 'default' case. The suggested fix always adds a single
// case clause for the missing enum members. The body of the case clause consists
// of a single statement:
//
//   panic(fmt.Sprintf("unhandled value: %v", v))
//
// where v is the expression in the switch statement's tag (in other words, the
// value being switched upon). If the switch statement's tag is a function or a
// method call the analyzer does not reuse the expression in the
// panic call because such calls could be mutative.
//
// Imports will be adjusted automatically to account for the package fmt dependency.
//
// Flags
//
// The analyzer accepts a boolean flag: --default-signifies-exhaustive.
// The flag, if set, indicates to the analyzer that switch statements
// are to be considered exhaustive as long as a 'default' case is present, even
// if all enum members aren't listed in the switch statements cases.
//
// Skip checking of specific switch statements
//
// The presence of the directive comment:
//
//   //exhaustive:ignore
//
// next to a switch statement indicates to the analyzer that it should skip
// checking of the switch statement. No diagnostics are reported.
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
	// Analyzer.Flags.BoolVar(&fCheckMaps, "maps", false, "check key exhaustiveness of map literals of enum key type, in addition to checking switch statements")
	Analyzer.Flags.BoolVar(&fDefaultSuffices, "default-signifies-exhaustive", false, "switch statements are considered exhaustive if a 'default' case is present")
}

var Analyzer = &analysis.Analyzer{
	Name:      "exhaustive",
	Doc:       "check for any non-exhaustive enum switch statements",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&enumsFact{}},
}

const IgnoreDirectivePrefix = "//exhaustive:ignore"

func containsIgnoreDirective(comments []*ast.Comment) bool {
	for _, c := range comments {
		if strings.HasPrefix(c.Text, IgnoreDirectivePrefix) {
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
	if len(e) != 0 {
		pass.ExportPackageFact(&enumsFact{entries: e})
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	comments := make(map[*ast.File]ast.CommentMap) // CommentMap per package file, lazily populated by reference

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
