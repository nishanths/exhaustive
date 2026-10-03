package finder

import (
	"go/ast"
	"go/token"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:       "finder",
	Doc:        "find various language entities",
	Run:        run,
	Requires:   []*analysis.Analyzer{inspect.Analyzer},
	ResultType: reflect.TypeFor[*Result](),
	FactTypes:  []analysis.Fact{new(isDefinedType)},
}

type Result struct {
	// TypeDecls is a collection of type name symbols from
	// type declarations.
	// The following forms of type declarations are considered:
	// type definition, alias declaration.
	// The Type() is either a *types.Alias or *types.Named.
	TypeDecls map[*types.TypeName]struct{}
}

func run(pass *analysis.Pass) (any, error) {
	findDefinedTypes(pass)

	result := &Result{
		TypeDecls: make(map[*types.TypeName]struct{}),
	}
	for _, f := range pass.AllObjectFacts() {
		switch f.Fact.(type) {
		case *isDefinedType:
			tn := f.Object.(*types.TypeName)
			result.TypeDecls[tn] = struct{}{}
		}
	}
	return result, nil
}

var universeAny = types.Universe.Lookup("any")

func findDefinedTypes(pass *analysis.Pass) {
	in := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	in.Preorder([]ast.Node{&ast.GenDecl{}}, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)
		if genDecl.Tok != token.TYPE {
			return
		}
		for _, spec := range genDecl.Specs {
			obj := pass.TypesInfo.Defs[spec.(*ast.TypeSpec).Name]
			tn, ok := obj.(*types.TypeName)
			if !ok {
				continue
			}
			switch tn.Type().(type) {
			case *types.Alias:
				pass.ExportObjectFact(tn, &isDefinedType{1})
			case *types.Named:
				pass.ExportObjectFact(tn, &isDefinedType{1})
			}
		}
	})
}

func assert(x bool) {
	if !x {
		panic("assertion failed")
	}
}

var _ analysis.Fact = (*isDefinedType)(nil)

type isDefinedType struct{ Version int }

func (*isDefinedType) AFact()           {}
func (f *isDefinedType) String() string { return "definedtype" }
