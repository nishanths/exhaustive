package exhaustive

import (
	"go/ast"
	"go/token"
	"go/types"
)

// enums holds the enum types and their members defined in a single package.
type enums map[string]*enumMembers // enum type name -> enum members

// enumMembers is the members for a single enum type.
// The zero value is ready to use.
type enumMembers struct {
	// Names is the enum member names,
	// in the order encountered in the AST.
	Names []string

	// NameToValue maps member name -> constVal.
	NameToValue map[string]string

	// ValueToNames maps constVal -> member names.
	// Note the use of []string for the element type of the map: Multiple
	// names can have the same value.
	ValueToNames map[string][]string
}

func (em *enumMembers) add(name string, constVal string) {
	em.Names = append(em.Names, name)

	if em.NameToValue == nil {
		em.NameToValue = make(map[string]string)
	}
	em.NameToValue[name] = constVal

	if em.ValueToNames == nil {
		em.ValueToNames = make(map[string][]string)
	}
	em.ValueToNames[constVal] = append(em.ValueToNames[constVal], name)
}

// Find the enums for the files in a package. The files is typically obtained from
// pass.Files and typesInfo is obtained from pass.TypesInfo.
func findEnums(files []*ast.File, typesInfo *types.Info) enums {
	possibleEnumTypes := make(map[string]struct{})

	// Gather possible enum types.
	findPossibleEnumTypes(files, typesInfo, func(name string) {
		possibleEnumTypes[name] = struct{}{}
	})

	pkgEnums := make(enums)

	// Gather enum members.
	findEnumMembers(files, typesInfo, possibleEnumTypes, func(memberName, typeName string, constVal string) {
		if _, ok := pkgEnums[typeName]; !ok {
			pkgEnums[typeName] = &enumMembers{}
		}
		pkgEnums[typeName].add(memberName, constVal)
	})

	return pkgEnums
}

func findPossibleEnumTypes(files []*ast.File, typesInfo *types.Info, found func(name string)) {
	for _, f := range files {
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
				t := s.(*ast.TypeSpec)
				obj := typesInfo.Defs[t.Name]
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
				case i&types.IsInteger != 0, i&types.IsFloat != 0, i&types.IsString != 0:
					found(named.Obj().Name())
				}
			}
		}
	}
}

func findEnumMembers(files []*ast.File, typesInfo *types.Info, knownEnumTypes map[string]struct{}, found func(memberName, typeName string, constVal string)) {
	for _, f := range files {
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			if gen.Tok != token.CONST {
				continue
			}
			for _, s := range gen.Specs {
				// Must be ValueSpec since we've filtered on token.CONST.
				v := s.(*ast.ValueSpec)
				for _, name := range v.Names {
					obj := typesInfo.Defs[name]
					namedType, ok := obj.Type().(*types.Named)
					if !ok {
						continue
					}
					if _, ok := knownEnumTypes[namedType.Obj().Name()]; !ok {
						continue
					}
					found(obj.Name(), namedType.Obj().Name(), determineConstVal(name, typesInfo))
				}
			}
		}
	}
}

func determineConstVal(name *ast.Ident, typesInfo *types.Info) string {
	c := typesInfo.Defs[name].(*types.Const)
	return c.Val().ExactString()
}
