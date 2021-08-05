// Utilities used across test files.

package exhaustive

import (
	"fmt"
	"reflect"
	"testing"
)

func ptrString(s string) *string { return &s }

func checkEqualf(t *testing.T, want, got interface{}, format string, args ...interface{}) bool {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%swant %v, got %v", logPrefix(format, args...), want, got)
		return false
	}
	return true
}

func checkEqual(t *testing.T, want, got interface{}) bool {
	t.Helper()
	return checkEqualf(t, want, got, "")
}

// v at the call site must not be an interface type. (in other words,
// calls to checkNil must pass a pointer to a concrete type's value.)
func checkNil(t *testing.T, v interface{}) bool {
	t.Helper()
	if isNil(t) {
		t.Errorf("want nil, got %v", v)
		return false
	}
	return true
}

// poor man's copy of https://github.com/stretchr/testify/blob/acba37e5db06f0093b465a7d47822bf13644b66c/assert/assertions.go#L520
// matching our need.
func isNil(object interface{}) bool {
	value := reflect.ValueOf(object)
	return value.IsNil()
}

func checkNoError(t *testing.T, err error) bool {
	t.Helper()
	if err != nil {
		t.Errorf("want nil error, got %s", err)
		return false
	}
	return true
}

func checkError(t *testing.T, err error) bool {
	t.Helper()
	if err == nil {
		t.Errorf("want error, got nil")
		return false
	}
	return true
}

func logPrefix(format string, args ...interface{}) string {
	if format != "" {
		return fmt.Sprintf(format, args...) + ": "
	}
	return ""
}
