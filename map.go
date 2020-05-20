package exhaustive

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

func checkMapLiterals(pass *analysis.Pass, inspect *inspector.Inspector, comments map[*ast.File]ast.CommentMap) {
	for _, f := range pass.Files {
		for _, d := range f.Decls {
			gen, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}
			if gen.Tok != token.VAR {
				continue // map literals have to be declared as "var"
			}
			for _, s := range gen.Specs {
				v := s.(*ast.ValueSpec)
				for idx, name := range v.Names {
					obj := pass.TypesInfo.Defs[name]
					if obj == nil {
						continue
					}

					mapType, ok := obj.Type().(*types.Map)
					if !ok {
						continue
					}

					keyType, ok := mapType.Key().(*types.Named)
					if !ok {
						continue
					}
					keyPkg := keyType.Obj().Pkg()
					if keyPkg == nil {
						// Doc comment: nil for labels and objects in the Universe scope.
						// This happens for the `error` type, for example.
						// Continuing would mean that ImportPackageFact panics.
						continue
					}

					var enums enumsFact
					if !pass.ImportPackageFact(keyPkg, &enums) {
						// Can't do anything further.
						continue
					}

					enumMembers, ok := enums.entries[keyType]
					if !ok {
						// Key type is not a known enum.
						continue
					}

					if containsIgnoreDirective(v.Doc.List) || containsIgnoreDirective(v.Comment.List) {
						continue
					}

					samePkg := keyPkg == pass.Pkg
					checkUnexported := samePkg

					log.Println(idx, v.Values, mapType)
					_ = enumMembers
					_ = checkUnexported
				}

			}

		}
	}
}
