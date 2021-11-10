package exhaustive

import (
	"bytes"
	"encoding/gob"
	"go/ast"
	"reflect"
	"regexp"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRegexpFlag(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		var v regexpFlag
		if got := v.value(); got != nil {
			t.Errorf("want nil, got %+v", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set(""); err != nil {
			t.Errorf("error unexpectedly non-nil: %v", err)
		}
		if got := v.value(); got != nil {
			t.Errorf("want nil, got %+v", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("bad input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set("("); err == nil {
			t.Errorf("error unexpectedly nil")
		}
		if got := v.value(); got != nil {
			t.Errorf("want nil, got %+v", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("good input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set("^foo$"); err != nil {
			t.Errorf("error unexpectedly non-nil: %v", err)
		}
		if v.value() == nil {
			t.Errorf("unexpectedly nil")
		}
		if !v.value().MatchString("foo") {
			t.Errorf("did not match")
		}
		if got, want := v.String(), regexp.MustCompile("^foo$").String(); got != want {
			t.Errorf("want %q, got %q", got, want)
		}
	})

	// The flag.Value interface doc says: "The flag package may call the
	// String method with a zero-valued receiver, such as a nil pointer."
	t.Run("String() nil receiver", func(t *testing.T) {
		var v *regexpFlag
		// expect no panic, and ...
		if got := v.String(); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})
}

func TestExhaustive(t *testing.T) {
	run := func(t *testing.T, pattern string, setup ...func()) {
		t.Helper()
		t.Run(pattern, func(t *testing.T) {
			resetFlags()
			for _, f := range setup {
				f()
			}
			analysistest.Run(t, analysistest.TestData(), Analyzer, pattern)
		})
	}

	// Enum discovery.
	run(t, "enum...")

	// Tests for the -check-generated flag.
	run(t, "generated-file/check-generated-off...")
	run(t, "generated-file/check-generated-on...", func() { fCheckGeneratedFiles = true })

	// Tests for the -default-signifies-exhaustive flag.
	// (For tests with this flag off, see other testdata packages
	// such as "general...".)
	run(t, "default-signifies-exhaustive/default-absent...", func() { fDefaultSignifiesExhaustive = true })
	run(t, "default-signifies-exhaustive/default-present...", func() { fDefaultSignifiesExhaustive = true })

	// Tests for the -ignore-enum-member flag.
	run(t, "ignore-enum-member...", func() {
		re := regexp.MustCompile(`_UNSPECIFIED$|^general/y\.Echinodermata$|^ignore-enum-member.User$`)
		fIgnoreEnumMembers = regexpFlag{re}
	})

	// Tests for -package-scope-only flag.
	run(t, "scope/allscope...")
	run(t, "scope/pkgscope...", func() { fPackageScopeOnly = true })

	// Switch statements with ignore directive comment should not be checked.
	run(t, "ignore-comment...")

	// For satisfy exhaustiveness, it is sufficient for each unique constant
	// value of the members to be listed, not each member by name.
	run(t, "duplicate-enum-value...")

	// Type alias switch statements.
	run(t, "typealias...")

	// General tests (a mixture).
	run(t, "general...")
}

func TestEnumMembersFact(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		e := enumMembersFact{
			Members: enumMembers{
				Names: []string{"Tundra", "Savanna", "Desert"},
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
		if want := "Tundra,Savanna,Desert"; want != e.String() {
			t.Errorf("want %v, got %v", want, e.String())
		}

		e = enumMembersFact{
			Members: enumMembers{
				Names: []string{"_", "add", "sub", "mul", "quotient", "remainder"},
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
		if want := "_,add,sub,mul,quotient,remainder"; want != e.String() {
			t.Errorf("want %v, got %v", want, e.String())
		}
	})
}

// This test exists to prevent regressions where changes made to a fact type used
// by the Analyzer makes the type fail to gob-encode/decode. Particuarly:
//
//  * gob values cannot seem to have nil pointers.
//  * fields must be exported to survive the encode/decode.
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

	// Fields should all be exported. And no pointer types should be present
	// unless you're absolutely sure, since nil pointers don't work with gob.
	t.Run("fields", func(t *testing.T) {
		switch v := fact.(type) {
		// NOTE: if there are more fact types, add them here.
		case *enumMembersFact:
			checkTypeEnumMembersFact(t, reflect.TypeOf(v).Elem())
		default:
			t.Errorf("unhandled type %T", v)
			return
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

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("want nil error, got %s", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("want error, got nil")
	}
}

func assertPanic(t *testing.T, f func(), wantPanicVal interface{}) {
	t.Helper()

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic, but no panic occurred")
			return
		}
		if !reflect.DeepEqual(r, wantPanicVal) {
			t.Errorf("want panic with: %+v, got panic with: %v", wantPanicVal, r)
			return
		}
	}()

	f()
}
