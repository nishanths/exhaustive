package exhaustive

import (
	"regexp"
	"testing"
)

func TestRegexpFlag(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		var v regexpFlag
		if v.r != nil {
			t.Errorf("want nil, got %+v", v.r)
		}
		_ = v.Get().(*regexp.Regexp) // should not panic
	})

	t.Run("empty input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set(""); err != nil {
			t.Errorf("error unexpectedly non-nil: %v", err)
		}
		if v.r != nil {
			t.Errorf("want nil, got %+v", v.r)
		}
		_ = v.Get().(*regexp.Regexp) // should not panic
	})

	t.Run("bad input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set("("); err == nil {
			t.Errorf("error unexpectedly nil")
		}
		if v.r != nil {
			t.Errorf("want nil, got %+v", v.r)
		}
		_ = v.Get().(*regexp.Regexp) // should not panic
	})

	t.Run("good input", func(t *testing.T) {
		var v regexpFlag
		if err := v.Set("^foo$"); err != nil {
			t.Errorf("error unexpectedly non-nil: %v", err)
		}
		if !v.Get().(*regexp.Regexp).MatchString("foo") {
			t.Errorf("did not match")
		}
	})
}
