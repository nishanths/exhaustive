package exhaustive

import (
	"bytes"
	"encoding/gob"
	"go/ast"
	"go/token"
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestEnumMembersFact(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		e := enumMembersFact{
			Members: enumMembers{
				Names: []string{"Tundra", "Savanna", "Desert"},
				NameToPos: map[string]token.Pos{
					"Tundra":  100,
					"Savanna": 200,
					"Desert":  300,
				},
				NameToValue: map[string]constantValue{
					"Tundra":  "1",
					"Savanna": "2",
					"Desert":  "3",
				},
				ValueToNames: map[constantValue][]string{
					"1": {"Tundra"},
					"2": {"Savanna"},
					"3": {"Desert"},
				},
			},
		}
		checkEnumMembersLiteral("Biome", e.Members)
		if want := "Tundra,Savanna,Desert"; e.String() != want {
			t.Errorf("got %v, want %v", e.String(), want)
		}

		e = enumMembersFact{
			Members: enumMembers{
				Names: []string{"_", "add", "sub", "mul", "quotient", "remainder"},
				NameToPos: map[string]token.Pos{
					"_":         1,
					"add":       11,
					"sub":       12,
					"mul":       33,
					"quotient":  34,
					"remainder": 35,
				},
				NameToValue: map[string]constantValue{
					"_":         "0",
					"add":       "1",
					"sub":       "2",
					"mul":       "3",
					"quotient":  "3",
					"remainder": "3",
				},
				ValueToNames: map[constantValue][]string{
					"0": {"_"},
					"1": {"add"},
					"2": {"sub"},
					"3": {"mul", "quotient", "remainder"},
				},
			},
		}
		checkEnumMembersLiteral("Token", e.Members)
		if want := "_,add,sub,mul,quotient,remainder"; e.String() != want {
			t.Errorf("got %v, want %v", e.String(), want)
		}
	})
}

// This test exists to prevent regressions where changes made to a fact type used
// by the Analyzer makes the type fail to gob-encode/decode. Particuarly:
//
//   - gob values cannot seem to have nil pointers.
//   - fields must be exported to survive the encode/decode.
//
// The test likely doesn't cover everything that could go wrong during gob
// encoding/decoding.
func TestFactsGob(t *testing.T) {
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

	// Ensure that all all fields all exported, and there are no pointer
	// types. Nil pointer values don't work with gob. We can't guarantee
	// non-nil values here, so we just disallow all pointer types.
	t.Run("fields", func(t *testing.T) {
		switch v := fact.(type) {
		// NOTE: if there are more fact types, add them here.
		case *enumMembersFact:
			checkTypeEnumMembersFact(t, reflect.TypeOf(v).Elem())
		default:
			t.Errorf("unhandled type %T", v)
		}
	})
}

func checkTypeEnumMembersFact(t *testing.T, factType reflect.Type) {
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
		{"NameToPos", "map[string]token.Pos"},
		{"NameToValue", "map[string]exhaustive.constantValue"},
		{"ValueToNames", "map[exhaustive.constantValue][]string"},
	})

	// Check that types such as token.Pos and constantValue have basic
	// underlying types (e.g. int, string).

	// check token.Pos.
	field, ok := enumMembersType.FieldByName("NameToPos")
	if !ok {
		t.Errorf("failed to find field")
		return
	}
	cvType := field.Type.Elem()
	if cvType.Kind() != reflect.Int {
		t.Errorf("unexpected kind %v", cvType.Kind())
	}

	// check constantValue.
	field, ok = enumMembersType.FieldByName("NameToValue")
	if !ok {
		t.Errorf("failed to find field")
		return
	}
	cvType = field.Type.Elem()
	if cvType.Kind() != reflect.String {
		t.Errorf("unexpected kind %v", cvType.Kind())
	}
}

func assertTypeFields(t *testing.T, typ reflect.Type, wantFields []wantField) {
	t.Helper()

	if got := typ.NumField(); got != len(wantFields) {
		t.Errorf("got %d, got %d", got, len(wantFields))
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
		if field.Name != wantFields[i].name {
			t.Errorf("got %q, want %q", field.Name, wantFields[i].name)
		}
		if field.Type.String() != wantFields[i].typ {
			t.Errorf("got %q, want %q", field.Type.String(), wantFields[i].typ)
		}
	}
}

type wantField struct {
	name string
	typ  string
}
