// Utilities used across test files.

package exhaustive

import (
	"testing"
)

func ptrString(s string) *string { return &s }

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
