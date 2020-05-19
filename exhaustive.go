package exhaustive

import (
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var flagset = flag.NewFlagSet("exhaustive", flag.ExitOnError)

var (
	fMap             = flagset.Bool("maps", true, "include maps in analysis")
	fDefaultSuffices = flagset.Bool("default-suffices", false, "consider switch to be exhaustive even if all enum members don't have individual cases, but 'default' case is present")
)

var Analyzer = &analysis.Analyzer{
	Name:      "exhaustive",
	Flags:     *flagset,
	Doc:       "check for non-exhaustive enum switch statements",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&Enums{}},
}

type Enums struct {
	Entries map[*types.Named][]types.Object
}

var _ analysis.Fact = (*Enums)(nil)

func (e *Enums) AFact() {}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	pkgEnums := make(map[*types.Named][]types.Object) // enum type -> enum members

	// Gather enum types.
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			if gen.Tok != token.TYPE {
				continue
			}
			for _, s := range gen.Specs {
				// must be a TypeSpec since we've filtered on token.TYPE,
				// but be defensive anyway.
				t, ok := s.(*ast.TypeSpec)
				if !ok {
					continue
				}
				obj := pass.TypesInfo.Defs[t.Name]
				if obj == nil {
					continue
				}

				named, ok := obj.Type().(*types.Named)
				if !ok {
					continue
				}
				basic, ok := named.Underlying().(*types.Basic)
				if !ok {
					continue
				}
				switch i := basic.Info(); {
				case i&types.IsInteger != 0:
					pkgEnums[named] = nil
				case i&types.IsFloat != 0:
					pkgEnums[named] = nil
				case i&types.IsString != 0:
					pkgEnums[named] = nil
				}
			}
		}
	}

	// Gather enum members.
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			if gen.Tok != token.CONST && gen.Tok != token.VAR {
				continue
			}
			for _, s := range gen.Specs {
				v, ok := s.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for idx, name := range v.Names {
					obj := pass.TypesInfo.Defs[name]
					if obj == nil {
						continue
					}
					named, ok := obj.Type().(*types.Named)
					if !ok {
						continue
					}
					// fmt.Println(pass.TypesInfo.Types[name].Value)

					members, ok := pkgEnums[named]
					if !ok {
						continue
					}
					fmt.Println("obj data", name.Name, name.Obj.Data)
					var cVal constant.Value
					// TODO finding out constant.Value this way won't always work
					// for iota
					// if val := v.Values[idx]; val == nil {
					// cVal = zeroValue(named.Underlying().(*types.Basic))
					// } else {
					cVal = pass.TypesInfo.Types[v.Values[idx]].Value
					// }
					fmt.Println(cVal)
					members = append(members, EnumMember{obj, cVal})
					pkgEnums[named] = members
				}
			}
		}
	}

	pass.ExportPackageFact(&Enums{
		Entries: pkgEnums,
	})

	filter := []ast.Node{&ast.SwitchStmt{}}
	inspect.Nodes(filter, func(n ast.Node, _ bool) bool {
		sw := n.(*ast.SwitchStmt)
		if sw.Tag == nil {
			return false
		}
		t := pass.TypesInfo.Types[sw.Tag]
		if !t.IsValue() {
			return false
		}
		named, ok := t.Type.(*types.Named)
		if !ok {
			return false
		}

		targetPkg := named.Obj().Pkg()
		if targetPkg == nil {
			// doc comment: nil for labels and objects in the Universe scope
			//
			// happens for type error, which has nil package.
			// continuing would mean that ImportPackageFact(targetPkg, ...) panics.
			return false
		}
		var enums Enums
		if !pass.ImportPackageFact(targetPkg, &enums) {
			// can't do anything further
			return false
		}

		// TODO
		// if it's the pass.Pkg that we're checking that the enum type resides in,
		// then exhaustiveness should check all enum members (exported and unexported).
		// if it's an external package, only check exported enum members.

		enumMembers, isEnum := enums.Entries[named]
		if !isEnum {
			return false
		}

		samePkg := named.Obj().Pkg() == pass.Pkg
		includeUnexported := samePkg

		if sw.Body == nil {
			return false
		}

		hitlist := make(map[string]bool)
		for _, m := range enumMembers {
			if m.obj.Exported() || includeUnexported {
				hitlist[m.obj.Name()] = false
			}
		}

		for _, stmt := range sw.Body.List {
			caseCl := stmt.(*ast.CaseClause) // based on doc comment on SwitchStmt.Body field
			if isDefaultCaseClause(caseCl) && *fDefaultSuffices {
				return false
			}
			for _, e := range caseCl.List {
				// TODO ensure it's a named type (not naked 3)
				// TODO ensure it's constant value that can be used in the hitlist
				// TODO: best way may be check only for Ident, selectorExpr, or things like
				// unary expr and paren expr that whittle down to ident or selectorexpr
				// then compre the string names against hitlist
				fmt.Println(pass.TypesInfo.Types[e].IsValue())
			}
		}
		fmt.Println(enumMembers)

		return false
	})

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

func isDefaultCaseClause(c *ast.CaseClause) bool {
	return c.List == nil
}
