package exhaustive

import (
	"go/ast"
	"go/types"
	"regexp"
)

// A checklist holds a set of enum member names that have to be
// accounted for in order to satisfy exhaustiveness.
//
// The found method checks off member names from the set, based on
// constant value. The remaining method returns the member names not
// accounted for.
type checklist struct {
	em     enumMembers
	checkl map[string]struct{}
}

func makeChecklist(em enumMembers, enumPkg *types.Package, includeUnexported bool, ignore *regexp.Regexp) *checklist {
	checkl := make(map[string]struct{})
	add := func(memberName string) {
		if memberName == "_" {
			// Blank identifier is often used to skip entries in iota lists.
			// Also, it can't be referenced anywhere (including in a switch
			// statement's cases), so it doesn't make sense to include it
			// as required member to satisfy exhaustiveness.
			return
		}
		if !ast.IsExported(memberName) && !includeUnexported {
			return
		}
		if ignore != nil && ignore.MatchString(enumPkg.Path()+"."+memberName) {
			return
		}
		checkl[memberName] = struct{}{}
	}

	for _, name := range em.Names {
		add(name)
	}

	return &checklist{
		em:     em,
		checkl: checkl,
	}
}

func (c *checklist) found(val constantValue) {
	// Delete all of the same-valued names.
	for _, name := range c.em.ValueToNames[val] {
		delete(c.checkl, name)
	}
}

func (c *checklist) remaining() map[string]struct{} {
	return c.checkl
}
