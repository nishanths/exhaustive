package exhaustive

import (
	"golang.org/x/tools/go/analysis"
)

// NOTE: Fact types must remain gob-coding compatible.
// See fact_gob_test.go.

var _ analysis.Fact = (*enumMembersFact)(nil)

type enumMembersFact struct{ Members enumMembers }

func (f *enumMembersFact) AFact()         {}
func (f *enumMembersFact) String() string { return f.Members.factString() }

// exportFact exports the enum members for the given enum type.
func exportFact(pass *analysis.Pass, enumTyp enumType, members *enumMembers) {
	pass.ExportObjectFact(enumTyp.factObject(), &enumMembersFact{*members})
}

// importFact imports the enum members for the given possible enum type. An
// (_, false) return indicates that no members exist for the given type, and by
// definition that the given type is not an enum type.
func importFact(pass *analysis.Pass, possibleEnumType enumType) (*enumMembers, bool) {
	var f enumMembersFact
	ok := pass.ImportObjectFact(possibleEnumType.factObject(), &f)
	if !ok {
		return nil, false
	}
	return &f.Members, true
}
