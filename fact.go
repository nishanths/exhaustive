package exhaustive

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

// NOTE: Fact types must remain gob-coding compatible.
// See fact_gob_test.go.

var _ analysis.Fact = (*enumMembersFact)(nil)

type enumMembersFact struct{ Members enumMembers }

func (f *enumMembersFact) AFact() {}

func (f *enumMembersFact) String() string {
	var buf strings.Builder
	for j, vv := range f.Members.Names {
		buf.WriteString(vv)
		// add comma separator between each enum member in an enum type
		if j != len(f.Members.Names)-1 {
			buf.WriteString(",")
		}
	}
	return buf.String()
}

// exportFact exports the enum members for the given enum type.
func exportFact(pass *analysis.Pass, enumTyp enumType, members *enumMembers) {
	pass.ExportObjectFact(enumTyp.object(), &enumMembersFact{*members})
}

// importFact imports the enum members for the given (possible) enum type.
// An (_, false) return indicates that no members exist for the given type.
func importFact(pass *analysis.Pass, possibleEnumType enumType) (*enumMembers, bool) {
	var f enumMembersFact
	ok := pass.ImportObjectFact(possibleEnumType.object(), &f)
	if !ok {
		return nil, false
	}
	return &f.Members, true
}
