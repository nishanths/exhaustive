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
	var sorted []enumType
	for enumTyp := range e.Enums {
		sorted = append(sorted, enumTyp)
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })

	var buf strings.Builder
	for i, enumTyp := range sorted {
		v := e.Enums[enumTyp]
		buf.WriteString(enumTyp.Name)
		buf.WriteString(":")

		for j, vv := range v.Names {
			buf.WriteString(vv)
			// add comma separator between each enum member in an enum type
			if j != len(v.Names)-1 {
				buf.WriteString(",")
			}
		}
		// add semicolon separator between each enum type
		if i != len(sorted)-1 {
			buf.WriteString("; ")
		}
	}
	return buf.String()
}
