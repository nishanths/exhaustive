package exhaustive

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type enums map[enumType][]enumMember

type enumType *types.Named
type enumMember types.Object

func gatherEnums(pass *analysis.Pass) enums {
	pkgEnums := make(enums)

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
				for _, name := range v.Names {
					obj := pass.TypesInfo.Defs[name]
					if obj == nil {
						continue
					}
					named, ok := obj.Type().(*types.Named)
					if !ok {
						continue
					}

					members, ok := pkgEnums[named]
					if !ok {
						continue
					}
					members = append(members, obj)
					pkgEnums[named] = members
				}
			}
		}
	}

	return pkgEnums
}
