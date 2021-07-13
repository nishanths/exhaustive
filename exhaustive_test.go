package exhaustive

import (
	"bytes"
	"encoding/gob"
	"go/ast"
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestEnum(t *testing.T) {
	resetFlags()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "enumvariant")
}

func TestSwitch(t *testing.T) {
	resetFlags()
	analysistest.Run(t, analysistest.TestData(), Analyzer, "switch/x", "switch/y")
}

func TestSwitch_ignorePattern(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		resetFlags()
		err := fIgnorePattern.Set("_UNSPECIFIED$|^switch/y.Echinodermata$")
		if !assertNoError(t, err) {
			return
		}
		analysistest.Run(t, analysistest.TestData(), Analyzer, "switch/ignorepattern")
	})
}

func TestSwitchFix(t *testing.T) {
	resetFlags()
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), Analyzer, "switchfix")
}

// This test exists to prevent regressions where changes made to a fact type used
// by the Analyzer makes the type fail to gob-encode/decode. Particuarly:
//  * gob values cannot seem to have nil pointers.
//  * fields must be exported to survive the encode/decode.
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

	t.Run("encode", func(t *testing.T) {
		if err := enc.Encode(factType); err != nil {
			t.Errorf("failed to gob-encode: %s", err)
			return
		}
	})

	t.Run("decode", func(t *testing.T) {
		if err := dec.Decode(factType); err != nil {
			t.Errorf("failed to gob-decode: %s", err)
			return
		}
	})

	t.Run("fields exported", func(t *testing.T) {
		switch v := factType.(type) {
		case *enumsFact:
			f, ok := reflect.TypeOf(v).Elem().FieldByName("Enums") // Enums field is exported obviously as we're referring to it with uppercase
			if !ok {
				t.Errorf("failed to find field")
			}

			enumMembers := f.Type.Elem().Elem()                            // 1st Elem(): obtain map value, 2nd Elem(): obtain pointer value
			assertEqual(t, "exhaustive.enumMembers", enumMembers.String()) // sanity check

			for i := 0; i < enumMembers.NumField(); i++ {
				if name := enumMembers.Field(i).Name; !ast.IsExported(name) {
					t.Errorf("field %q not exported", name)
				}
			}
		default:
			t.Errorf("unhandled type %T", v)
			return
		}
	})
}
