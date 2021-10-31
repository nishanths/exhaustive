// Utilities used across test files.

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
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic, but did not")
			return
		}
		if !reflect.DeepEqual(r, wantPanicVal) {
			t.Errorf("wanted panic with: %+v, got panic with: %v", wantPanicVal, r)
			return
		}
	}()

	f()
}
