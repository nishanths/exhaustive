package exhaustive

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ast/inspector"
)

// constantValue is a constant.Value.ExactString().
type constantValue string

// enums contains the enum types and their members for a single package.
type enums map[enumType]*enumMembers

// Represents an enum type (or sometimes a potential enum type).
type enumType struct{ tn *types.TypeName }

func (et enumType) String() string       { return et.tn.String() } // for debugging
func (et enumType) object() types.Object { return et.tn }          // types.Object for fact export

// enumMembers is the members for a single enum type.
// The zero value is ready to use.
type enumMembers struct {
	// TODO: In the future, depending on how we design type alias analysis,
	// NameToValue may not work correctly if there are multiple same-named type
	// alias enum members; only one of them can be saved in the map.

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

// NOTE: This type must be usable as a map key
// (comparison by '==' must do the right thing)
type enumTypeAndScope struct {
	scope *types.Scope
	typ   enumType
}

func findEnums(pkgScopeOnly bool, pkg *types.Package, inspect *inspector.Inspector, info *types.Info) enums {
	var out enums = make(map[enumType]*enumMembers)

	// -- Find possible enum types --

	enumTypes := make(map[enumTypeAndScope]struct{})

	f := func(tn *types.TypeName, scope *types.Scope) {
		if scope != pkg.Scope() && pkgScopeOnly {
			return
		}
		e := enumTypeAndScope{scope, enumType{tn}}
		enumTypes[e] = struct{}{}
	}

	inspect.Preorder([]ast.Node{&ast.GenDecl{}}, func(n ast.Node) {
		possibleEnumTypes(n.(*ast.GenDecl), info, f)
	})

	// -- Find enum members --

	g := func(enumTyp enumType, memberName string, val constantValue) {
		if _, ok := out[enumTyp]; !ok {
			out[enumTyp] = &enumMembers{}
		}
		out[enumTyp].add(memberName, val)
	}

	inspect.Preorder([]ast.Node{&ast.GenDecl{}}, func(n ast.Node) {
		possibleEnumMembers(n.(*ast.GenDecl), info, enumTypes, g)
	})

	return out
}

// possibleEnumTypes reports types that could possibly be enum types
// in the given GenDecl. It calls found for enum type found.
func possibleEnumTypes(gen *ast.GenDecl, info *types.Info, found func(tn *types.TypeName, scope *types.Scope)) {
	if gen.Tok != token.TYPE {
		return
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
		if isBlankIdentifier(obj) {
			// These objects have a nil parent scope (so trying to match
			// a const with this enum type will fail).
			// Also, we have no real purpose to record them.
			continue
		}

		{
			_, ok := obj.(*types.TypeName)
			assert(ok, "obj must be *types.TypeName")
		}

		named, ok := obj.Type().(*types.Named)
		if !ok {
			continue
		}

		// RHS type of `named` should either be an enum type (named with
		// with underlying valid basic type) or directory
		// be a valid basic type. We can handle both cases
		// by checking `named.Underlying()`.

		basic, ok := named.Underlying().(*types.Basic)
		if !ok {
			continue
		}
		switch i := basic.Info(); {
		case i&types.IsInteger != 0, i&types.IsFloat != 0, i&types.IsString != 0:
			found(obj.(*types.TypeName), obj.Parent())
		}
	}
}

func possibleEnumMembers(gen *ast.GenDecl, info *types.Info, possibleEnumTypes map[enumTypeAndScope]struct{}, found func(enumTyp enumType, memberName string, val constantValue)) {
	if gen.Tok != token.CONST {
		return
	}

	for _, s := range gen.Specs {
		v := s.(*ast.ValueSpec) // because gen.Tok == token.CONST
		for _, name := range v.Names {
			obj := info.Defs[name]
			if obj == nil {
				continue
			}
			if isBlankIdentifier(obj) {
				// These objects have a nil parent scope (so trying to match
				// the const with its enum type will fail).
				// Also, we have no real purpose to record them.
				continue
			}

			{
				_, ok := obj.(*types.Const)
				assert(ok, "obj must be *types.Const")
			}

			named, ok := obj.Type().(*types.Named)
			if !ok {
				continue
			}
			tn := named.Obj()

			// Enum type's scope and enum member's scope must be the same.
			// If they're not, don't consider the const a member.
			e := enumTypeAndScope{obj.Parent(), enumType{tn}}
			if _, ok := possibleEnumTypes[e]; !ok {
				continue
			}

			found(enumType{tn}, obj.Name(), determineConstVal(name, info))
		}
	}
}

func determineConstVal(name *ast.Ident, info *types.Info) constantValue {
	c := info.Defs[name].(*types.Const)
	return constantValue(c.Val().ExactString())
}

func isBlankIdentifier(obj types.Object) bool {
	return obj.Name() == "_" // NOTE: go/types/decl.go does a direct comparison like this
}
