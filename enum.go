package exhaustive

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
)

// constantValue is a constant.Value.ExactString().
type constantValue string

// %p fmt formatted string of an address.
type addr string

// enums contains the enum types and their members for a single package.
type enums map[enumType]*enumMembers

// enumType is represents an enum type. It is the key type for the enums type.
// It is designed to be gob-coding compatible.
type enumType struct {
	Name string
	Addr addr // *types.Named addr
}

// enumMembers is the members for a single enum type.
// The zero value is ready to use.
type enumMembers struct {
	Names        []string                   // enum member names, AST order
	NameToValue  map[string]constantValue   // enum member name -> constant value
	ValueToNames map[constantValue][]string // constant value -> enum member names
}

func (em *enumMembers) add(name string, val constantValue) {
	em.Names = append(em.Names, name)

	if em.NameToValue == nil {
		em.NameToValue = make(map[string]constantValue)
	}
	em.NameToValue[name] = val

	if em.ValueToNames == nil {
		em.ValueToNames = make(map[constantValue][]string)
	}
	em.ValueToNames[val] = append(em.ValueToNames[val], name)
}

func typesNamedAddr(t *types.Named) addr {
	return addr(fmt.Sprintf("%p", t))
}

// Find the enums for the files in a package. The files is typically obtained from
// pass.Files and typesInfo is obtained from pass.TypesInfo.
func findEnums(files []*ast.File, typesInfo *types.Info) enums {
	possibleEnumTypes := make(map[enumType]struct{})

	// Gather possible enum types.
	findPossibleEnumTypes(files, typesInfo, func(enumTyp enumType) {
		possibleEnumTypes[enumTyp] = struct{}{}
	})

	pkgEnums := make(enums)

	// Gather enum members.
	findEnumMembers(files, typesInfo, possibleEnumTypes, func(memberName string, enumTyp enumType, val constantValue) {
		if _, ok := pkgEnums[enumTyp]; !ok {
			pkgEnums[enumTyp] = &enumMembers{}
		}
		pkgEnums[enumTyp].add(memberName, val)
	})

	return pkgEnums
}

func findPossibleEnumTypes(files []*ast.File, typesInfo *types.Info, found func(enumTyp enumType)) {
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
					found(enumType{named.Obj().Name(), typesNamedAddr(named)})
				}
			}
		}
	}
}

func findEnumMembers(files []*ast.File, typesInfo *types.Info, possibleEnumTypes map[enumType]struct{}, found func(memberName string, enumTyp enumType, val constantValue)) {
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
					enumTyp := enumType{namedType.Obj().Name(), typesNamedAddr(namedType)}
					if _, ok := possibleEnumTypes[enumTyp]; !ok {
						continue
					}
					found(obj.Name(), enumTyp, determineConstVal(name, typesInfo))
				}
			}
		}
	}
}

func determineConstVal(name *ast.Ident, typesInfo *types.Info) constantValue {
	c := typesInfo.Defs[name].(*types.Const)
	return constantValue(c.Val().ExactString())
}
