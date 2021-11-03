package exhaustive

import (
	"bytes"
	"encoding/gob"
	"go/ast"
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// This test exists to prevent regressions where changes made to a fact type used
// by the Analyzer makes the type fail to gob-encode/decode. Particuarly:
//
//  * gob values cannot seem to have nil pointers.
//  * fields must be exported to survive the encode/decode.
//
// The test doesn't cover everything that could go wrong during gob
// encoding/decoding.
func TestFactsGobCompatible(t *testing.T) {
	// The go/analysis package does this internally, but we need to do it
	// manually here for the test.
	for _, typ := range Analyzer.FactTypes {
		gob.Register(typ)
	}

	for _, typ := range Analyzer.FactTypes {
		t.Run("fact type: "+reflect.TypeOf(typ).String(), func(t *testing.T) {
			checkOneFactType(t, typ)
		})
	}
}

func checkOneFactType(t *testing.T, factType analysis.Fact) {
	t.Helper()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	// Should be able to gob-encode.
	t.Run("encode", func(t *testing.T) {
		if err := enc.Encode(factType); err != nil {
			t.Errorf("failed to gob-encode: %s", err)
			return
		}
	})

	// Should be able to gob-decode.
	t.Run("decode", func(t *testing.T) {
		if err := dec.Decode(factType); err != nil {
			t.Errorf("failed to gob-decode: %s", err)
			return
		}
	})

	// Fields should all be exported.
	t.Run("fields exported", func(t *testing.T) {
		switch v := factType.(type) {
		// NOTE: if there are more fact types, add them here.
		case *enumsFact:
			checkTypeEnumsFact(t, reflect.TypeOf(v).Elem())
		default:
			t.Errorf("unhandled type %T", v)
			return
		}
	})
}

func checkTypeEnumsFact(t *testing.T, enumsFactType reflect.Type) {
	t.Helper()

	assertTypeFields(t, enumsFactType, []wantField{
		{"Enums", "exhaustive.enums"},
	})

	// Check underlying type of the "Enums" field.
	f, ok := enumsFactType.FieldByName("Enums")
	if !ok {
		t.Errorf("failed to find field")
		return
	}
	if f.Type.Kind() != reflect.Map {
		t.Errorf("want reflect.Map, got %v (%v)", f.Type.Kind(), f.Type.Kind().String())
		return
	}
	keyType, elemType := f.Type.Key(), f.Type.Elem()
	if keyType.String() != "string" {
		t.Errorf("want key type string, got %v", keyType.String())
		return
	}
	if elemType.String() != "*exhaustive.enumMembers" {
		t.Errorf("want elem type *exhaustive.enumMembers, got %v", elemType.String())
		return
	}

	enumMembersType := elemType.Elem() // call Elem() on pointer type to get value type
	checkTypeEnumMembers(t, enumMembersType)
}

func checkTypeEnumMembers(t *testing.T, enumMembersType reflect.Type) {
	t.Helper()
	assertTypeFields(t, enumMembersType, []wantField{
		{"Names", "[]string"},
		{"NameToValue", "map[string]string"},
		{"ValueToNames", "map[string][]string"},
	})
}

func assertTypeFields(t *testing.T, typ reflect.Type, wantFields []wantField) {
	t.Helper()

	if got := typ.NumField(); got != len(wantFields) {
		t.Errorf("want %d, got %d", len(wantFields), got)
		return
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !ast.IsExported(field.Name) {
			t.Errorf("field %q not exported", field.Name)
		}
		if wantFields[i].name != field.Name {
			t.Errorf("want %q, got %q", wantFields[i].name, field.Name)
		}
		if wantFields[i].typ != field.Type.String() {
			t.Errorf("want %q, got %q", wantFields[i].typ, field.Type.String())
		}
	}
}

type wantField struct {
	name string
	typ  string
}
