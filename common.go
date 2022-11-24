package exhaustive

import (
	"flag"
	"go/types"
	"regexp"
	"sort"
	"strings"
)

// diagnosticMissingMembers constructs the list of missing enum members,
// suitable for use in a reported diagnostic message.
func diagnosticMissingMembers(missingMembers map[string]struct{}, em enumMembers) []string {
	// inASTOrder sorts the given names in AST order. The AST position
	// of each name is determined using the astPositions map. Names with
	// smaller position values appear in the AST before names with large
	// position values.
	//
	// The slice is sorted in place. It is also returned for
	// convenience.
	inASTOrder := func(names []string, astPositions map[string]int) []string {
		sort.Slice(names, func(i, j int) bool {
			return astPositions[names[i]] < astPositions[names[j]]
		})
		return names
	}

	// byConstVal groups member names by constant value.
	byConstVal := func(names map[string]struct{}, nameToValue map[string]constantValue) map[constantValue][]string {
		ret := make(map[constantValue][]string)
		for name := range names {
			val := nameToValue[name]
			ret[val] = append(ret[val], name)
		}
		return ret
	}

	// indices maps each string in the input slice to its index.
	indices := func(names []string) map[string]int {
		ret := make(map[string]int, len(names))
		for i, name := range names {
			ret[name] = i
		}
		return ret
	}

	astPositions := indices(em.Names)

	var groups []string
	for _, names := range byConstVal(missingMembers, em.NameToValue) {
		group := strings.Join(inASTOrder(names, astPositions), "|")
		groups = append(groups, group)
	}
	return inASTOrder(groups, astPositions)
}

// diagnosticEnumTypeName returns a string representation of an enum
// type for use in reported diagnostics.
func diagnosticEnumTypeName(enumType *types.TypeName, samePkg bool) string {
	if samePkg {
		return enumType.Name()
	}
	return enumType.Pkg().Name() + "." + enumType.Name()
}

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
