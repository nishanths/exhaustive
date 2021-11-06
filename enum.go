package exhaustive

import (
	"go/ast"
	"go/token"
	"go/types"
)

// constantValue is a constant.Value.ExactString().
type constantValue string

// enums contains the enum types and their members for a single package.
type enums map[enumType]*enumMembers

// Represents an enum type (or sometimes a potential enum type).
type enumType struct{ named *types.Named }

func (et enumType) String() string       { return et.named.String() }
func (et enumType) name() string         { return et.named.Obj().Name() }
func (et enumType) object() types.Object { return et.named.Obj() }

// enumMembers is the members for a single enum type.
// The zero value is ready to use.
type enumMembers struct {
	// TODO: NameToValue doesn't work correctly if there are multiple blank
	// identifiers ("_"); only one of them can be saved.
	// There may be a similar issue for same-named type alias enum members, in
	// the future, depending on how we design type alias analysis.

	Names        []string                   // enum member names, AST order
	NameToValue  map[string]constantValue   // enum member name -> constant value
	ValueToNames map[constantValue][]string // constant value -> enum member names
}

func (em *enumMembers) add(name string, val constantValue) {
	if em.NameToValue == nil {
		em.NameToValue = make(map[string]constantValue)
	}
	if em.ValueToNames == nil {
		em.ValueToNames = make(map[constantValue][]string)
	}

	em.Names = append(em.Names, name)
	em.NameToValue[name] = val
	em.ValueToNames[val] = append(em.ValueToNames[val], name)
}

// Find the enums for the files in a package. The files are typically obtained from
// pass.Files and info is obtained from pass.TypesInfo.
func findEnums(files []*ast.File, info *types.Info) enums {
	// Gather possible enum types.
	enumTypes := make(map[*types.Named]struct{})
	findPossibleEnumTypes(files, info, func(named *types.Named) {
		enumTypes[named] = struct{}{}
	})

	out := make(enums)

	// Gather enum members.
	findEnumMembers(files, info, enumTypes, func(memberName string, enumTyp enumType, val constantValue) {
		if _, ok := out[enumTyp]; !ok {
			out[enumTyp] = &enumMembers{}
		}
		out[enumTyp].add(memberName, val)
	})

	return out
}

func findPossibleEnumTypes(files []*ast.File, info *types.Info, found func(named *types.Named)) {
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
				t := s.(*ast.TypeSpec) // because gen.Tok == token.TYPE
				if t.Assign.IsValid() {
					// TypeSpec is AliasSpec; we don't support it at the moment.
					// Additionally:
					// In type T1 = T2,  info.Defs[t.Name].Type() results in the object on the right-hand side.
					// In type T1 T2,    info.Defs[t.Name].Type() results in the object of the left-hand side.
					// This needs to be resolved.
					continue
				}
				obj := info.Defs[t.Name]
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
					found(named)
				}
			}
		}
	}
}

func findEnumMembers(files []*ast.File, info *types.Info, possibleEnumTypes map[*types.Named]struct{}, found func(memberName string, enumTyp enumType, val constantValue)) {
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
				v := s.(*ast.ValueSpec) // because gen.Tok == token.CONST
				for _, name := range v.Names {
					obj := info.Defs[name]
					namedType, ok := obj.Type().(*types.Named)
					if !ok {
						continue
					}
					if _, ok := possibleEnumTypes[namedType]; !ok {
						continue
					}
					found(obj.Name(), enumType{namedType}, determineConstVal(name, info))
				}
			}
		}
	}
}

func determineConstVal(name *ast.Ident, info *types.Info) constantValue {
	c := info.Defs[name].(*types.Const)
	return constantValue(c.Val().ExactString())
}
