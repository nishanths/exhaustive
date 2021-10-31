package exhaustive

import (
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var _ analysis.Fact = (*enumsFact)(nil)

type enumsFact struct {
	Enums enums
}

func (e *enumsFact) AFact() {}

func (e *enumsFact) String() string {
	// sort for stability (required for testing)
	var sortedKeys []string
	for k := range e.Enums {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var buf strings.Builder
	for i, k := range sortedKeys {
		v := e.Enums[k]
		buf.WriteString(k)
		buf.WriteString(":")

		for j, vv := range v.Names {
			buf.WriteString(vv)
			// add comma separator between each enum member in an enum type
			if j != len(v.Names)-1 {
				buf.WriteString(",")
			}
		}
		// add semicolon separator between each enum type
		if i != len(sortedKeys)-1 {
			buf.WriteString("; ")
		}
	}
	return buf.String()
}
