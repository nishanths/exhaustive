package exhaustive

import (
	"regexp"
	"testing"
)

func TestRegexpFlag(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		var v regexpFlag
		checkNil(t, v.r)
		_ = v.Get().(*regexp.Regexp) // should not panic
	})

	t.Run("empty input", func(t *testing.T) {
		var v regexpFlag
		checkNoError(t, v.Set(""))
		checkNil(t, v.r)
		_ = v.Get().(*regexp.Regexp) // should not panic
	})

	t.Run("bad input", func(t *testing.T) {
		var v regexpFlag
		checkError(t, v.Set("("))
		checkNil(t, v.r)
		_ = v.Get().(*regexp.Regexp) // should not panic
	})

	t.Run("good input", func(t *testing.T) {
		var v regexpFlag
		checkNoError(t, v.Set("^foo$"))
		if !v.Get().(*regexp.Regexp).MatchString("foo") {
			t.Errorf("want regexp match, but did not match")
		}
	})

}
