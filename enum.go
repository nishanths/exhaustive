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
				// Must be TypeSpec since we've filtered on token.TYPE.
				t, ok := s.(*ast.TypeSpec)
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
				// Must be ValueSpec since we've filtered on token.CONST, token.VAR.
				v := s.(*ast.ValueSpec)
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
