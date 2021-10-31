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
			checkEnumsFactExported(t, v)
		default:
			t.Errorf("unhandled type %T", v)
			return
		}
	})
}

// check that fields in all types references in v's definition are exported.
func checkEnumsFactExported(t *testing.T, v *enumsFact) {
	t.Helper()

	if c := reflect.TypeOf(v).Elem().NumField(); c != 1 {
		t.Errorf("unexpected number of fields: %d, wanted: 1 (test needs update?)", c)
		return
	}

	// The single field is known to named 'Enums'; find it.
	// The 'Enums' field is exported obviously as we're referring to it with uppercase.
	f, ok := reflect.TypeOf(v).Elem().FieldByName("Enums")
	if !ok {
		t.Errorf("failed to find field")
		return
	}

	// Sanity check: We know the Enums field to be a map[string]*enumMembers.
	// Check that it matches our knowledge.
	keyType, elemType := f.Type.Key(), f.Type.Elem()
	if keyType.String() != "string" {
		t.Errorf("want key type string, got %v", keyType.String())
		return
	}
	if elemType.String() != "*exhaustive.enumMembers" {
		t.Errorf("want elem type *exhaustive.enumMembers, got %v", elemType.String())
		return
	}

	enumMembers := elemType.Elem() // pointer value

	// Check that all fields are exported.
	for i := 0; i < enumMembers.NumField(); i++ {
		// TODO: Go1.17 will add StructField.IsExported(), maybe that is appropriate to use here?
		// https://github.com/golang/go/issues/41563
		if name := enumMembers.Field(i).Name; !ast.IsExported(name) {
			t.Errorf("field %q not exported", name)
		}
	}
}
