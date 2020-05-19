package exhaustive

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var flagset = flag.NewFlagSet("exhaustive", flag.ExitOnError)

var (
	fMap = flagset.Bool("maps", true, "include maps in analysis")
)

var Analyzer = &analysis.Analyzer{
	Name:      "exhaustive",
	Flags:     *flagset,
	Doc:       "check for non-exhaustive enum switch statements",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&E{}},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	var enumTypes []*types.Basic

	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok {
				// don't care about declarations
				continue
			}
			if gen.Tok != token.TYPE {
				// don't care about others such as import/const declarations
				continue
			}
			for _, s := range gen.Specs {
				// must be a TypeSpec since we've filtered on token.TYPE,
				// but be defensive anyway
				t, ok := s.(*ast.TypeSpec)
				if !ok {
					continue
				}
				basic, ok := pass.TypesInfo.Types[t.Type].Type.(*types.Basic)
				if !ok {
					continue
				}
				switch i := basic.Info(); {
				case i&types.IsInteger != 0:
					fmt.Println(basic.Name(), f.Name, pass.Pkg.Name())
					enumTypes = append(enumTypes, basic)
				case i&types.IsFloat != 0:
					enumTypes = append(enumTypes, basic)
				case i&types.IsString != 0:
					enumTypes = append(enumTypes, basic)
				}
			}
		}
	}

	for _, e := range enumTypes {
		fmt.Printf("%#v\n", e)
	}

	// gather top-level type declarations from ast.File.Decls of these underlying
	// types
	//
	// 	string
	//  byte
	//  rune
	//  numerical
	//
	// gather top-level const and var declarations from ast.File.Decls
	// that have concrete type of the types from the previous step.
	//
	// consider the types a as enum types and the variables as relevant
	// enum members.

	// explore every switch statement's tag, if the tag is present.
	// see if the tag's type is an enum type for the package. if it is,
	// look through the cases of the switch to ensure that all of the enum
	// members are listed.

	filter := []ast.Node{&ast.SwitchStmt{}}
	// inspect.WithStack([]ast.Node{}, f)
	inspect.WithStack(filter, func(n ast.Node, push bool, _ []ast.Node) bool {
		sw := n.(*ast.SwitchStmt)
		if sw.Tag != nil {
			// fmt.Println(pass.Pkg.Path(), sw.Tag, pass.TypesInfo.Types[sw.Tag].Type)
		}
		return false
	})
	return nil, nil
}

type E struct{}

var _ analysis.Fact = (*E)(nil)

func (e *E) AFact() {}
