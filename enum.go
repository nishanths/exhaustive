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

type enumData struct {
	typ     enumType
	members enumMembers
}

// Represents an enum type (or a potential enum type).
// It is a defined/named type's name.
type enumType struct{ *types.TypeName }

func (et enumType) String() string           { return et.TypeName.String() } // for debugging
func (et enumType) scope() *types.Scope      { return et.TypeName.Parent() } // scope that the type is declared in
func (et enumType) factObject() types.Object { return et.TypeName }          // types.Object for fact export

// enumMembers is the members for a single enum type.
// The zero value is ready to use.
type enumMembers struct {
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

func (em *enumMembers) String() string { return em.factString() } // for debugging

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
				t := s.(*ast.TypeSpec)
				if t.Assign.IsValid() {
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

func findEnums(pkgScopeOnly, excludeTypeAlias bool, pkg *types.Package, inspect *inspector.Inspector, info *types.Info) map[enumType]enumMembers {
	// toSlice := func(v map[enumType]enumData) []enumData {
	// 	var ret []enumData
	// 	for _, vv := range v {
	// 		ret = append(ret, vv)
	// 	}
	// 	return ret
	// }

	debug.Println("--", pkg, "--")

	// _, aliases, consts := collectConsts(inspect)

	// possTypes := make(map[*types.TypeName]struct{})
	// aliasLhses := make(map[*types.TypeName]struct{})
	// aliasRhses := make(map[*types.TypeName]struct{})
	result := make(map[enumType]enumMembers)
	// claimedConsts := make(map[int]struct{}) // consts indexes claimed by typedef enum members.

	// -- Find possible enum types --

	// for _, t := range typedefs {
	// 	tn, scope, ok := possibleEnumType(t, info)
	// 	if !ok {
	// 		continue
	// 	}
	// 	if scope != pkg.Scope() && pkgScopeOnly {
	// 		continue
	// 	}
	// 	possTypes[tn] = struct{}{}
	// }

	// -- Find members of enum types --

	inspect.Preorder([]ast.Node{&ast.GenDecl{}}, func(n ast.Node) {
		gen := n.(*ast.GenDecl)
		if gen.Tok != token.CONST {
			return
		}
		for _, s := range gen.Specs {
			for _, name := range s.(*ast.ValueSpec).Names {
				enumTyp, memberName, val, ok := typedefEnumMember(name, pkg, info)
				if !ok {
					continue
				}
				if enumTyp.scope() != pkg.Scope() && pkgScopeOnly {
					continue
				}
				v, ok := result[enumTyp]
				v.add(memberName, val)
				result[enumTyp] = v
			}
		}
	})

	// for i, c := range consts {
	// 	enumTyp, memberName, val, ok := typedefEnumMember(c, pkg, info)
	// 	if !ok {
	// 		continue
	// 	}

	// 	v, ok := result[enumTyp]
	// 	if !ok {
	// 		v = enumData{typ: enumTyp}
	// 	}
	// 	v.members.add(memberName, val)
	// 	result[enumTyp] = v

	// 	claimedConsts[i] = struct{}{}
	// }

	// -- Exit early if done --

	// if excludeTypeAlias {
	// 	return toSlice(result)
	// }

	// -- Find possible type alias enum types --

	// for _, a := range aliases {
	// 	aliasTN, aliasScope, enumTN, _, ok := aliasToPossibleEnumType(a, info)
	// 	if !ok {
	// 		continue
	// 	}
	// 	// Alias to enum type allowed only in package scope.
	// 	if aliasScope != pkg.Scope() {
	// 		continue
	// 	}
	// 	aliasLhses[aliasTN] = struct{}{}
	// 	aliasRhses[enumTN] = struct{}{}
	// }

	// debug.Println(possTypes)
	// debug.Println(aliasLhses)
	// debug.Println(aliasRhses)

	// -- Find enum members of type alias enum types --

	// Exclude consts claimed by previous steps.
	// var remainingConsts []*ast.Ident
	// for i, c := range consts {
	// 	if _, ok := claimedConsts[i]; ok {
	// 		continue
	// 	}
	// 	remainingConsts = append(remainingConsts, c)
	// }

	// for _, c := range remainingConsts {
	// 	_, _, _, ok := possibleAliasEnumMember(c, info, aliasLhses)
	// 	if !ok {
	// 		continue
	// 	}
	// 	// Members of alias enum type only allowed in package scope.
	// }

	return result

	// TODO: Delete this comment when done./
	//
	// TypeSpec is AliasSpec; we don't support it at the moment.
	// Additionally:
	// In type T1 = T2,  info.Defs[t.Name].Type() results in the object on the right-hand side.
	// In type T1 T2,    info.Defs[t.Name].Type() results in the object of the left-hand side.
	// This needs to be resolved.
}

// func possibleEnumType(typedef *ast.TypeSpec, info *types.Info) (*types.TypeName, *types.Scope, bool) {
// 	if typedef.Assign.IsValid() {
// 		panic("TypeSpec is alias type")
// 	}

// 	obj := info.Defs[typedef.Name]
// 	if obj == nil {
// 		return nil, nil, false
// 	}
// 	if isBlankIdentifier(obj) {
// 		// These objects have a nil parent scope (so trying to match
// 		// a const with this enum type will fail).
// 		// Also, we have no real purpose to record them.
// 		return nil, nil, false
// 	}

// 	_, ok := obj.(*types.TypeName)
// 	assert(ok, "obj must be *types.TypeName")

// 	if !enumNamedBasic(obj.Type()) {
// 		return nil, nil, false
// 	}
// 	return obj.(*types.TypeName), obj.Parent(), true
// }

// func aliasToPossibleEnumType(alias *ast.TypeSpec, info *types.Info) (*types.TypeName, *types.Scope, *types.TypeName, *types.Scope, bool) {
// 	if !alias.Assign.IsValid() {
// 		panic("TypeSpec is not alias type")
// 	}

// 	obj := info.Defs[alias.Name]
// 	if obj == nil {
// 		return nil, nil, nil, nil, false
// 	}
// 	if isBlankIdentifier(obj) {
// 		return nil, nil, nil, nil, false
// 	}

// 	_, ok := obj.(*types.TypeName)
// 	assert(ok, "obj must be *types.TypeName")

// 	rhsTV, ok := info.Types[alias.Type]
// 	if !ok {
// 		return nil, nil, nil, nil, false
// 	}
// 	rhs := rhsTV.Type

/*
	type T1 = int  // not allowed (alias -> valid basic type)
	type T2 = T3   // not allowed (alias -> alias -> valid basic type)
	type T9 = T8   // not allowed (alias -> alias -> ... -> alias -> valid basic type)
	type T4 = T5   // possible    (alias -> named type -> valid basic type)
	type T6 = T7   // possible    (alias -> alias -> ... -> alias -> named type -> valid basic type)

	type T3 = int
	type T8 = T3
	type T5 int // NOTE: does not matter right now that T5 has no known members
	type T7 = T5

	// The above holds also true if
	// T3, T5, etc. are in different packages from T2, T4, etc. respectively.
*/
// 	if !enumNamedBasic(rhs) {
// 		return nil, nil, nil, nil, false
// 	}
// 	return obj.(*types.TypeName),
// 		obj.Parent(),
// 		rhs.(*types.Named).Obj(),
// 		rhs.(*types.Named).Obj().Parent(),
// 		true
// }

func typedefEnumMember(constName *ast.Ident, pkg *types.Package, info *types.Info) (et enumType, memberName string, val constantValue, ok bool) {
	obj := info.Defs[constName]
	if obj == nil {
		return enumType{}, "", "", false
	}
	if isBlankIdentifier(obj) {
		// These objects have a nil parent scope.
		// Also, we have no real purpose to record them.
		return enumType{}, "", "", false
	}

	_, ok = obj.(*types.Const)
	assert(ok, "obj must be *types.Const")

	if !enumNamedBasic(obj.Type()) {
		return enumType{}, "", "", false
	}

	named, ok := obj.Type().(*types.Named)
	if !ok {
		return enumType{}, "", "", false
	}
	tn := named.Obj()

	// The constant and its type must be in the same package.
	if tn.Pkg() != obj.Pkg() {
		return enumType{}, "", "", false
	}
	// Enum type's scope and enum member's scope must be the same.
	// If they're not, don't consider the const a member.
	if tn.Parent() != obj.Parent() {
		return enumType{}, "", "", false
	}

	return enumType{tn}, obj.Name(), determineConstVal(constName, info), true
}

// func possibleAliasEnumMember(constName *ast.Ident, info *types.Info, aliasLhses, aliasRhses map[*types.TypeName]struct{}) (et enumType, memberName string, val constantValue, ok bool) {
// 	obj := info.Defs[constName]
// 	if obj == nil {
// 		return enumType{}, "", "", false
// 	}
// 	if isBlankIdentifier(obj) {
// 		return enumType{}, "", "", false
// 	}

// 	_, ok = obj.(*types.Const)
// 	assert(ok, "obj must be *types.Const")

// 	constType := obj.Type()

// 	debug.Printf("%v %T %v %T", obj, obj, obj.Type(), obj.Type())

// 	if !enumNamedBasic(constType) {
// 		return enumType{}, "", "", false
// 	}

// 	// Is the constant's type a known (typedef) enum type?
// 	if named, ok := constType.(*types.Named); ok {
// 		// Enum type's scope and enum member's scope must be the same.
// 		// If they're not, don't consider the const a member.
// 		tn := named.Obj()

// 		if tn.Parent() != obj.Parent() {
// 			// TODO:
// 		}
// 		// if _, ok := possibleTypes[tn]; ok {
// 		// return enumType{tn}, obj.Name(), determineConstVal(constName, info), true
// 		// }
// 	}

// 	// Is it a known alias RHS type?
// 	// if _, ok := aliasRhses[constType]; ok {
// 	// TODO
// 	// return enumType{"TODO"}, obj.Name(), determineConstVal(constName, info), true
// 	// }

// 	return enumType{}, "", "", false
// }

func determineConstVal(name *ast.Ident, info *types.Info) constantValue {
	c := info.Defs[name].(*types.Const)
	return constantValue(c.Val().ExactString())
}

func isBlankIdentifier(obj types.Object) bool {
	return obj.Name() == "_" // NOTE: go/types/decl.go does a direct comparison like this
}

func enumValidBasic(basic *types.Basic) bool {
	switch i := basic.Info(); {
	case i&types.IsInteger != 0, i&types.IsFloat != 0, i&types.IsString != 0:
		return true
	}
	return false
}

// enumNamedBasic returns whether the type t is a named type whose underlying
// type is a valid basic type to form an enum.
// A type that passes this check meets the definition of an enum type.
// Note that
//   enumNamedBasic(t) == true => t.(*types.Named)
func enumNamedBasic(t types.Type) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	basic, ok := named.Underlying().(*types.Basic)
	if !ok || !enumValidBasic(basic) {
		return false
	}
	return true
}
