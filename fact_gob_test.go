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
// The test likely doesn't cover everything that could go wrong during gob
// encoding/decoding.
func TestFactsGobCompatible(t *testing.T) {
	// The go/analysis package does this internally, but we need to do it
	// manually here for the test.
	for _, typ := range Analyzer.FactTypes {
		gob.Register(typ)
	}

	for _, typ := range Analyzer.FactTypes {
		t.Run("fact type "+reflect.TypeOf(typ).String(), func(t *testing.T) {
			checkOneFactType(t, typ)
		})
	}
}

func checkOneFactType(t *testing.T, fact analysis.Fact) {
	t.Helper()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	// Should be able to gob-encode.
	t.Run("gob encode", func(t *testing.T) {
		if err := enc.Encode(fact); err != nil {
			t.Errorf("failed to gob-encode: %s", err)
			return
		}
	})

	// Should be able to gob-decode.
	t.Run("gob decode", func(t *testing.T) {
		if err := dec.Decode(fact); err != nil {
			t.Errorf("failed to gob-decode: %s", err)
			return
		}
	})

	// Fields should all be exported. And no pointer types should be present
	// unless you're absolutely sure, since nil pointers don't work with gob.
	t.Run("fields", func(t *testing.T) {
		switch v := fact.(type) {
		// NOTE: if there are more fact types, add them here.
		case *enumMembersFact:
			checkEnumMembersFact(t, reflect.TypeOf(v).Elem())
		default:
			t.Errorf("unhandled type %T", v)
			return
		}
	})
}

func checkEnumMembersFact(t *testing.T, factType reflect.Type) {
	t.Helper()

	assertTypeFields(t, factType, []wantField{
		{"Members", "exhaustive.enumMembers"},
	})

	field, ok := factType.FieldByName("Members")
	if !ok {
		t.Errorf("failed to find field")
		return
	}
	enumMembersType := field.Type
	checkTypeEnumMembers(t, enumMembersType)
}

func checkTypeEnumMembers(t *testing.T, enumMembersType reflect.Type) {
	t.Helper()

	assertTypeFields(t, enumMembersType, []wantField{
		{"Names", "[]string"},
		{"NameToValue", "map[string]exhaustive.constantValue"},
		{"ValueToNames", "map[exhaustive.constantValue][]string"},
	})

	field, ok := enumMembersType.FieldByName("NameToValue")
	if !ok {
		t.Errorf("failed to find field")
		return
	}
	cvType := field.Type.Elem()
	checkTypeConstantValue(t, cvType)
}

func checkTypeConstantValue(t *testing.T, cvType reflect.Type) {
	t.Helper()
	if cvType.Kind() != reflect.String {
		t.Errorf("unexpected kind %v", cvType.Kind())
	}
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
		if field.Type.Kind() == reflect.Ptr {
			t.Errorf("field %q is pointer", field.Name)
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
