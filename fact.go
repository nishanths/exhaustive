package exhaustive

import (
	"bytes"
	"encoding/gob"
	"go/token"
	"sort"

	"golang.org/x/tools/go/analysis"
)

// NOTE: Fact types must remain gob-coding compatible.
// See TestFactsGob.

var _ analysis.Fact = (*enumMembersFact)(nil)

type enumMembersFact struct{ Members enumMembers }

func (f *enumMembersFact) AFact()         {}
func (f *enumMembersFact) String() string { return f.Members.factString() }

func (f *enumMembersFact) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(newEnumMembersFactGob(f.Members))
	return buf.Bytes(), err
}

func (f *enumMembersFact) GobDecode(data []byte) error {
	var fact enumMembersFactGob
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&fact); err != nil {
		return err
	}
	f.Members = fact.Members.enumMembers()
	return nil
}

type enumMembersFactGob struct {
	Members enumMembersGob
}

type enumMembersGob struct {
	Names        []string
	NameToPos    []enumMemberNamePos
	NameToValue  []enumMemberNameValue
	ValueToNames []enumMemberValueNames
}

type enumMemberNamePos struct {
	Name string
	Pos  token.Pos
}

type enumMemberNameValue struct {
	Name  string
	Value constantValue
}

type enumMemberValueNames struct {
	Value constantValue
	Names []string
}

func newEnumMembersFactGob(members enumMembers) enumMembersFactGob {
	return enumMembersFactGob{
		Members: newEnumMembersGob(members),
	}
}

func newEnumMembersGob(members enumMembers) enumMembersGob {
	return enumMembersGob{
		Names:        append([]string(nil), members.Names...),
		NameToPos:    newEnumMemberNamePosSlice(members.NameToPos),
		NameToValue:  newEnumMemberNameValueSlice(members.NameToValue),
		ValueToNames: newEnumMemberValueNamesSlice(members.ValueToNames),
	}
}

func newEnumMemberNamePosSlice(m map[string]token.Pos) []enumMemberNamePos {
	keys := make([]string, 0, len(m))
	for name := range m {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	result := make([]enumMemberNamePos, 0, len(keys))
	for _, name := range keys {
		result = append(result, enumMemberNamePos{
			Name: name,
			Pos:  m[name],
		})
	}
	return result
}

func newEnumMemberNameValueSlice(m map[string]constantValue) []enumMemberNameValue {
	keys := make([]string, 0, len(m))
	for name := range m {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	result := make([]enumMemberNameValue, 0, len(keys))
	for _, name := range keys {
		result = append(result, enumMemberNameValue{
			Name:  name,
			Value: m[name],
		})
	}
	return result
}

func newEnumMemberValueNamesSlice(m map[constantValue][]string) []enumMemberValueNames {
	keys := make([]constantValue, 0, len(m))
	for value := range m {
		keys = append(keys, value)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	result := make([]enumMemberValueNames, 0, len(keys))
	for _, value := range keys {
		result = append(result, enumMemberValueNames{
			Value: value,
			Names: append(
				[]string(nil),
				m[value]...,
			),
		})
	}
	return result
}

func (members enumMembersGob) enumMembers() enumMembers {
	return enumMembers{
		Names:        append([]string(nil), members.Names...),
		NameToPos:    members.nameToPos(),
		NameToValue:  members.nameToValue(),
		ValueToNames: members.valueToNames(),
	}
}

func (members enumMembersGob) nameToPos() map[string]token.Pos {
	if len(members.NameToPos) == 0 {
		return nil
	}
	result := make(map[string]token.Pos, len(members.NameToPos))
	for _, entry := range members.NameToPos {
		result[entry.Name] = entry.Pos
	}
	return result
}

func (members enumMembersGob) nameToValue() map[string]constantValue {
	if len(members.NameToValue) == 0 {
		return nil
	}
	result := make(map[string]constantValue, len(members.NameToValue))
	for _, entry := range members.NameToValue {
		result[entry.Name] = entry.Value
	}
	return result
}

func (members enumMembersGob) valueToNames() map[constantValue][]string {
	if len(members.ValueToNames) == 0 {
		return nil
	}
	result := make(map[constantValue][]string, len(members.ValueToNames))
	for _, entry := range members.ValueToNames {
		result[entry.Value] = append([]string(nil), entry.Names...)
	}
	return result
}

// exportFact exports the enum members for the given enum type.
func exportFact(pass *analysis.Pass, enumTyp enumType, members enumMembers) {
	pass.ExportObjectFact(enumTyp.factObject(), &enumMembersFact{members})
}

// importFact imports the enum members for the given possible enum type.
// An (_, false) return indicates that the enum type is not a known one.
func importFact(pass *analysis.Pass, possibleEnumType enumType) (enumMembers, bool) {
	var f enumMembersFact
	if !pass.ImportObjectFact(possibleEnumType.factObject(), &f) {
		return enumMembers{}, false
	}
	return f.Members, true
}
