// Package exhaustive provides an analyzer that helps ensure enum switch statements
// are exhaustive. The analyzer also provides fixes to make the offending switch
// statements exhaustive (see "Fixes" section).
//
// See "cmd/exhaustive" for the related command line program.
//
// Definition of enum
//
// The Go programming language does not have a specification for enums.
// This program uses the following reasonable specification instead.
//
// An enum type is a package-level named integer, float, or
// string type. An enum type must have associated with it one or more
// package-level variables of the named type in the package. These variables
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
// Switch statement exhaustiveness
//
// An enum switch statment is exhaustive if it has cases for each of the enum's members.
//
// For an enum type defined in the same package as the switch statement, both
// exported and unexported enum members must be present in order to consider
// the switch exhaustive. On the other hand, for an enum type defined
// in an external package it is sufficient for just exported enum members
// to be present in order to consider the switch exhaustive.
//
// Flags
//
// The analyzer accepts a boolean flag: -default-signifies-exhaustive.
// The flag, if set, indicates to the analyzer that switch statements
// are to be considered exhaustive as long as a 'default' case is present, even
// if all enum members aren't listed in the switch statements cases.
//
// Skip checking of specific switch statements
//
// If the following directive comment:
//
//   //exhaustive:ignore
//
// is associated with a switch statement, the analyzer skips
// checking of the switch statement and no diagnostics are reported.
//
// Fixes
//
// The analyzer suggests fixes for a switch statement if it is not exhaustive
// and does not have a 'default' case. The suggested fix always adds a single
// case clause for the missing enum members. The body of the case clause consists
// of the statement:
//
//   panic(fmt.Sprintf("unhandled value: %v", v))
//
// where v is the expression in the switch statement's tag (in other words, the
// value being switched upon). If the switch statement's tag is a function or a
// method call the analyzer does not suggest a fix, as reusing the call expression
// in the panic/fmt.Sprintf call could be mutative.
//
// The rationale for the fix is that it might be better to panic loudly on
// existing unhandled or impossible cases than to let them slip by quietly unnoticed.
// An even better fix would, of course, be to manually inspect the sites reported
// by the package and handle the missing cases if necessary.
//
// Imports will be adjusted automatically to account for the package fmt dependency.
//
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
	Analyzer.Flags.BoolVar(&fDefaultSuffices, "default-signifies-exhaustive", false, "switch statements are considered exhaustive if a 'default' case is present")
}

var Analyzer = &analysis.Analyzer{
	Name:      "exhaustive",
	Doc:       "check for any non-exhaustive enum switch statements",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&enumsFact{}},
}

// IgnoreDirectivePrefix is used to exclude checking of specific switch statements.
// See https://godoc.org/github.com/nishanths/exhaustive#hdr-Skip_checking_of_specific_switch_statements
// for details.
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
