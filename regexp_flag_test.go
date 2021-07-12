package exhaustive

import (
	"regexp"
	"testing"
)

func TestRegexpFlag(t *testing.T) {
	t.Run("bad input", func(t *testing.T) {
		var v regexpFlag
		requireError(t, v.Set("("))
		_ = v.Get().(*regexp.Regexp) // should not panic
	})

	t.Run("good input", func(t *testing.T) {
		var v regexpFlag
		requireNoError(t, v.Set("^foo$"))
		if !v.Get().(*regexp.Regexp).MatchString("foo") {
			t.Errorf("want regexp match, but did not match")
		}
	})

	t.Run("unset", func(t *testing.T) {
		var v regexpFlag
		requireNoError(t, v.Set(""))
		_ = v.Get().(*regexp.Regexp) // should not panic
	})
}
