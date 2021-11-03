// Shared testing utilities for use in test files.

package exhaustive

import (
	"reflect"
	"testing"
)

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
