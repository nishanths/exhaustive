package exhaustive

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"

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
		fIgnorePattern = "_UNSPECIFIED$|^switch/y.Echinodermata$"
		analysistest.Run(t, analysistest.TestData(), Analyzer, "switch/ignorepattern")
	})
}

func TestSwitchFix(t *testing.T) {
	resetFlags()
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), Analyzer, "switchfix")
}

// NOTE: This test doesn't cover everything that could go wrong during gob
// encoding/decoding.
func TestGobCompatible(t *testing.T) {
	// The go/analysis package does this internally, but we need to do it
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
