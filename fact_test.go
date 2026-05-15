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
		checkEnumMembersLiteral(t, "Biome", e.Members)
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
		checkEnumMembersLiteral(t, "Token", e.Members)
		if want := "_,add,sub,mul,quotient,remainder"; e.String() != want {
			t.Errorf("got %v, want %v", e.String(), want)
		}
	})
}

// This test exists to prevent regressions where changes made to a fact type used
// by the Analyzer makes the type fail to gob-encode/decode. Particularly:
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

func TestEnumMembersFactGobDeterministic(t *testing.T) {
	fact := enumMembersFact{
		Members: enumMembers{
			Names: []string{
				"Member00",
				"Member01",
				"Member02",
				"Member03",
				"Member04",
				"Member05",
				"Member06",
				"Member07",
				"Member08",
				"Member09",
				"Member10",
				"Member11",
				"Member12",
				"Member13",
				"Member14",
				"Member15",
			},
			NameToPos: map[string]token.Pos{
				"Member00": 100,
				"Member01": 101,
				"Member02": 102,
				"Member03": 103,
				"Member04": 104,
				"Member05": 105,
				"Member06": 106,
				"Member07": 107,
				"Member08": 108,
				"Member09": 109,
				"Member10": 110,
				"Member11": 111,
				"Member12": 112,
				"Member13": 113,
				"Member14": 114,
				"Member15": 115,
			},
			NameToValue: map[string]constantValue{
				"Member00": "0",
				"Member01": "1",
				"Member02": "2",
				"Member03": "3",
				"Member04": "4",
				"Member05": "5",
				"Member06": "6",
				"Member07": "7",
				"Member08": "8",
				"Member09": "9",
				"Member10": "10",
				"Member11": "11",
				"Member12": "12",
				"Member13": "13",
				"Member14": "14",
				"Member15": "15",
			},
			ValueToNames: map[constantValue][]string{
				"0":  {"Member00"},
				"1":  {"Member01"},
				"2":  {"Member02"},
				"3":  {"Member03"},
				"4":  {"Member04"},
				"5":  {"Member05"},
				"6":  {"Member06"},
				"7":  {"Member07"},
				"8":  {"Member08"},
				"9":  {"Member09"},
				"10": {"Member10"},
				"11": {"Member11"},
				"12": {"Member12"},
				"13": {"Member13"},
				"14": {"Member14"},
				"15": {"Member15"},
			},
		},
	}

	first := encodeFact(t, &fact)
	var decoded enumMembersFact
	if err := gob.NewDecoder(bytes.NewReader(first)).Decode(&decoded); err != nil {
		t.Fatalf("failed to gob-decode: %s", err)
	}
	if !reflect.DeepEqual(decoded, fact) {
		t.Fatalf("decoded fact mismatch:\ngot  %#v\nwant %#v", decoded, fact)
	}

	for i := 0; i < 100; i++ {
		got := encodeFact(t, &fact)
		if !bytes.Equal(got, first) {
			t.Fatalf("gob encoding changed on iteration %d: first length %d, got length %d", i, len(first), len(got))
		}
	}
}

func encodeFact(t *testing.T, fact *enumMembersFact) []byte {
	t.Helper()

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(fact); err != nil {
		t.Fatalf("failed to gob-encode: %s", err)
	}
	return buf.Bytes()
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
