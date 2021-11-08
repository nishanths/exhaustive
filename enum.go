package exhaustive

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/ast/inspector"
)

// constantValue is a constant.Value.ExactString().
type constantValue string

// enums contains the enum types and their members for a single package.
type enums map[enumType]*enumMembers

// Represents an enum type (or sometimes a potential enum type).
type enumType struct{ tn *types.TypeName }

func (et enumType) String() string           { return et.tn.String() } // for debugging
func (et enumType) factObject() types.Object { return et.tn }          // types.Object for fact export

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

func (em *enumMembers) factString() string {
	var buf strings.Builder
	for j, vv := range em.Names {
		buf.WriteString(vv)
		// add comma separator between each enum member
		if j != len(em.Names)-1 {
			buf.WriteString(",")
		}
	}
	return buf.String()
}

// Traverses the AST of the inspector and returns the TypeDef, AliasDecl, and constant identifiers.
func collectSpecs(inspect *inspector.Inspector) (typedefs []*ast.TypeSpec, aliases []*ast.TypeSpec, consts []*ast.Ident) {
	inspect.Preorder([]ast.Node{&ast.GenDecl{}}, func(n ast.Node) {
		gen := n.(*ast.GenDecl)
		for _, s := range gen.Specs {
			switch gen.Tok {
			case token.TYPE:
				if t := s.(*ast.TypeSpec); t.Assign.IsValid() {
					aliases = append(aliases, t)
				} else {
					typedefs = append(typedefs, t)
				}
			case token.CONST:
				consts = append(consts, s.(*ast.ValueSpec).Names...)
			}
		}
	})
	return
}

// NOTE: This type must be usable as a map key
// (comparison by '==' must do the right thing)
type enumTypeAndScope struct {
	scope *types.Scope
	typ   enumType
}

func findEnums(pkgScopeOnly bool, pkg *types.Package, inspect *inspector.Inspector, info *types.Info) enums {
	typedefs, _, consts := collectSpecs(inspect)
	result := make(map[enumType]*enumMembers)

	// -- Find possible typedef enum types --

	typedefEnumTypes := make(map[enumTypeAndScope]struct{})

	for _, t := range typedefs {
		tn, scope, ok := possibleTypedefEnumType(t, info)
		if !ok {
			continue
		}
		if scope != pkg.Scope() && pkgScopeOnly {
			continue
		}
		e := enumTypeAndScope{scope, enumType{tn}}
		typedefEnumTypes[e] = struct{}{}
	}

	// -- Find enum members of typedef enum types --

	for _, c := range consts {
		enumTyp, memberName, val, ok := possibleTypedefEnumMembers(c, info, typedefEnumTypes)
		if !ok {
			continue
		}
		if _, ok := result[enumTyp]; !ok {
			result[enumTyp] = &enumMembers{}
		}
		result[enumTyp].add(memberName, val)
	}

	// TODO:
	// -- Find possible type alias enum types --

	// TypeSpec is AliasSpec; we don't support it at the moment.
	// Additionally:
	// In type T1 = T2,  info.Defs[t.Name].Type() results in the object on the right-hand side.
	// In type T1 T2,    info.Defs[t.Name].Type() results in the object of the left-hand side.
	// This needs to be resolved.

	// -- Find enum members of type alias enum types --

	return result
}

func possibleTypedefEnumType(typedef *ast.TypeSpec, info *types.Info) (*types.TypeName, *types.Scope, bool) {
	if typedef.Assign.IsValid() {
		panic("TypeSpec is a type alias")
	}

	obj := info.Defs[typedef.Name]
	if obj == nil {
		return nil, nil, false
	}
	if isBlankIdentifier(obj) {
		// These objects have a nil parent scope (so trying to match
		// a const with this enum type will fail).
		// Also, we have no real purpose to record them.
		return nil, nil, false
	}

	_, ok := obj.(*types.TypeName)
	assert(ok, "obj must be *types.TypeName")

	named, ok := obj.Type().(*types.Named)
	if !ok {
		return nil, nil, false
	}

	// RHS type of `named` should either be an enum type (named with
	// with underlying valid basic type) or directory
	// be a valid basic type. We can handle both cases
	// by checking `named.Underlying()`.

	basic, ok := named.Underlying().(*types.Basic)
	if !ok {
		return nil, nil, false
	}
	switch i := basic.Info(); {
	case i&types.IsInteger != 0, i&types.IsFloat != 0, i&types.IsString != 0:
		return obj.(*types.TypeName), obj.Parent(), true
	}

	return nil, nil, false
}

func possibleTypedefEnumMembers(constName *ast.Ident, info *types.Info, possibleTypes map[enumTypeAndScope]struct{}) (et enumType, memberName string, val constantValue, ok bool) {
	obj := info.Defs[constName]
	if obj == nil {
		return enumType{}, "", "", false
	}
	if isBlankIdentifier(obj) {
		// These objects have a nil parent scope (so trying to match
		// the const with its enum type will fail).
		// Also, we have no real purpose to record them.
		return enumType{}, "", "", false
	}

	_, ok = obj.(*types.Const)
	assert(ok, "obj must be *types.Const")

	named, ok := obj.Type().(*types.Named)
	if !ok {
		return enumType{}, "", "", false
	}
	tn := named.Obj()

	// Enum type's scope and enum member's scope must be the same.
	// If they're not, don't consider the const a member.
	e := enumTypeAndScope{obj.Parent(), enumType{tn}}
	if _, ok := possibleTypes[e]; !ok {
		return enumType{}, "", "", false
	}

	return enumType{tn}, obj.Name(), determineConstVal(constName, info), true
}

func determineConstVal(name *ast.Ident, info *types.Info) constantValue {
	c := info.Defs[name].(*types.Const)
	return constantValue(c.Val().ExactString())
}

func isBlankIdentifier(obj types.Object) bool {
	return obj.Name() == "_" // NOTE: go/types/decl.go does a direct comparison like this
}
