package exhaustive

import (
	"go/types"
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
		group := inASTOrder(names, astPositions)
		groups = append(groups, strings.Join(group, "|"))
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
