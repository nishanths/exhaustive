package exhaustive

import (
	"errors"
	"reflect"
	"regexp"
	"testing"
)

func TestRegexpFlag(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		var v regexpFlag
		if got := v.re; got != nil {
			t.Errorf("got %+v, want nil", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set(""); err != nil {
			t.Errorf("error unexpectedly non-nil: %v", err)
		}
		if got := v.re; got != nil {
			t.Errorf("got %+v, want nil", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})

	t.Run("bad input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set("("); err == nil {
			t.Errorf("error unexpectedly nil")
		}
		if got := v.re; got != nil {
			t.Errorf("got %+v, want nil", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})

	t.Run("good input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set("^foo$"); err != nil {
			t.Errorf("error unexpectedly non-nil: %v", err)
		}
		if v.re == nil {
			t.Errorf("unexpectedly nil")
		}
		if !v.re.MatchString("foo") {
			t.Errorf("did not match")
		}
		if got, want := v.String(), regexp.MustCompile("^foo$").String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	// The flag.Value interface doc says: "The flag package may call the
	// String method with a zero-valued receiver, such as a nil pointer."
	t.Run("String nil receiver", func(t *testing.T) {
		var v *regexpFlag
		// expect no panic, and ...
		if got := v.String(); got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})
}

func TestStringsFlag(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var v stringsFlag
		if err := v.Set(""); err != nil {
			t.Errorf("error unexpectedly non-nil: %v", err)
		}
		if got := len(v.elements); got != 0 {
			t.Errorf("got %d, want 0 length", got)
		}
		if got := v.String(); got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})

	t.Run("happy path", func(t *testing.T) {
		var v stringsFlag
		if err := v.Set("a, b,bb, c   ,d "); err != nil {
			t.Errorf("error unexpectedly non-nil: %v", err)
		}
		if got, want := v.elements, []string{"a", "b", "bb", "c", "d"}; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := v.String(), "a,b,bb,c,d"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("filter error", func(t *testing.T) {
		errBoom := errors.New("boom")

		var v stringsFlag
		v.filter = func(e string) error {
			if e == "bb" {
				return errBoom
			}
			return nil
		}

		err := v.Set("a, b,bb, c   ,d ")
		if err == nil {
			t.Errorf("error unexpectedly nil: %v", err)
		}
		if err != errBoom {
			t.Errorf("got %v, want %v", err, errBoom)
		}
	})

	// The flag.Value interface doc says: "The flag package may call the
	// String method with a zero-valued receiver, such as a nil pointer."
	t.Run("String nil receiver", func(t *testing.T) {
		var v *stringsFlag
		// expect no panic, and ...
		if got := v.String(); got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})
}
