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
	// If a name is missing in the map, it means that the name does not have a
	// known constVal.
	NameToValue map[string]string

	// ValueToNames maps constVal -> member names.
	// Note the use of []string for the element type of the map: Multiple
	// names can have the same value.
	// Names that don't have a constVal will not have a corresponding entry in this map.
	ValueToNames map[string][]string
}

func (em *enumMembers) addWithConstVal(name string, constVal string) {
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

// add adds an encountered member name. If the constant value for the enum member
// is known, use addWithConstVal instead.
func (em *enumMembers) add(name string) {
	em.Names = append(em.Names, name)
}

// Find the enums for the files in a package. The files is typically obtained from
// pass.Files and typesInfo is obtained from pass.TypesInfo.
func findEnums(files []*ast.File, typesInfo *types.Info) enums {
	knownEnumTypes := make(map[string]struct{})

	// Gather possible enum types.
	findPossibleEnumTypes(files, typesInfo, func(name string) {
		knownEnumTypes[name] = struct{}{}
	})

	pkgEnums := make(enums)

	// Gather enum members.
	findEnumMembers(files, typesInfo, knownEnumTypes, func(memberName, typeName string, constVal string, constValOk bool) {
		if _, ok := pkgEnums[typeName]; !ok {
			pkgEnums[typeName] = &enumMembers{}
		}
		if constValOk {
			pkgEnums[typeName].addWithConstVal(memberName, constVal)
		} else {
			pkgEnums[typeName].add(memberName)
		}
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

func findEnumMembers(files []*ast.File, typesInfo *types.Info, knownEnumTypes map[string]struct{}, found func(memberName, typeName string, constVal string, ok bool)) {
	for _, f := range files {
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
				for i, name := range v.Names {
					obj := typesInfo.Defs[name]
					if obj == nil {
						continue
					}

					named, ok := obj.Type().(*types.Named)
					if !ok {
						continue
					}

					if _, ok := knownEnumTypes[named.Obj().Name()]; !ok {
						continue
					}

					var value ast.Expr
					if len(v.Values) > i {
						value = v.Values[i]
					}

					cv, ok := determineConstVal(name, value, typesInfo)
					found(obj.Name(), named.Obj().Name(), cv, ok)
				}
			}
		}
	}
}

func determineConstVal(name *ast.Ident, value ast.Expr, typesInfo *types.Info) (string, bool) {
	if s, ok := determineConstValFromName(name, typesInfo); ok {
		return s, true
	}
	if value != nil {
		return determineConstValFromValue(value, typesInfo)
	}
	return "", false
}

func determineConstValFromName(name *ast.Ident, typesInfo *types.Info) (string, bool) {
	nameObj := typesInfo.Defs[name]
	if nameObj == nil {
		return "", false
	}

	c, ok := nameObj.(*types.Const)
	if !ok {
		return "", false
	}
	if c.Val() == nil {
		return "", false
	}
	return c.Val().ExactString(), true
}

func determineConstValFromValue(value ast.Expr, typesInfo *types.Info) (string, bool) {
	con, ok := typesInfo.Types[value]
	if !ok {
		return "", false
	}
	if con.Value == nil {
		return "", false
	}
	return con.Value.ExactString(), true
}
