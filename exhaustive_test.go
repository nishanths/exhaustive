package exhaustive

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestEnum(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "enumvariant/...")
}

func TestSwitch(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "switch/...")
}

func TestSwitchFix(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), Analyzer, "switchfix/...")
}

func TestGobCompatible(t *testing.T) {
	// The analysis package does this internally, but we need to do it
	// manually here for the test.
	for _, typ := range Analyzer.FactTypes {
		gob.Register(typ)
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	for _, typ := range Analyzer.FactTypes {
		t.Run(reflect.TypeOf(typ).String(), func(t *testing.T) {
			buf.Reset()
			if err := enc.Encode(typ); err != nil {
				t.Errorf("failed to encode: %s", err)
				return
			}
			if err := dec.Decode(typ); err != nil {
				t.Errorf("failed to decode: %s", err)
				return
			}
		})
	}
}
