package exhaustive

import (
	"flag"
	"regexp"
	"strings"
)

var _ flag.Value = (*regexpFlag)(nil)
var _ flag.Value = (*stringsFlag)(nil)

// regexpFlag implements flag.Value for parsing
// regular expression flag inputs.
type regexpFlag struct{ r *regexp.Regexp }

func (v *regexpFlag) String() string {
	if v == nil || v.r == nil {
		return ""
	}
	return v.r.String()
}

func (v *regexpFlag) Set(expr string) error {
	if expr == "" {
		v.r = nil
		return nil
	}

	r, err := regexp.Compile(expr)
	if err != nil {
		return err
	}

	v.r = r
	return nil
}

func (v *regexpFlag) regexp() *regexp.Regexp { return v.r }

// stringsFlag implements flag.Value for parsing a comma-separated
// string list.  Surrounding space is stripped from each element of the
// list. If filter is non-nil it is called for each element in the
// input.
type stringsFlag struct {
	elements []string
	filter   func(string) error
}

func (v *stringsFlag) String() string {
	if v == nil {
		return ""
	}
	return strings.Join(v.elements, ",")
}

func (v *stringsFlag) filterFunc() func(string) error {
	if v.filter != nil {
		return v.filter
	}
	return func(_ string) error { return nil }
}

func (v *stringsFlag) Set(input string) error {
	for _, el := range strings.Split(input, ",") {
		el = strings.TrimSpace(el)
		if err := v.filter(el); err != nil {
			return err
		}
		v.elements = append(v.elements, el)
	}
	return nil
}
