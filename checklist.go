package exhaustive

import (
	"fmt"
	"go/ast"
	"go/types"
	"regexp"
)

// A checklist is the set of enum member names that should be listed in a switch
// statement's case clauses in order for the switch to be exhaustive. The found
// method marks a member as being listed in the switch, so, in usage, a checklist
// is the set of yet unsatisfied enum members.
//
// Only interact via its methods. It is not safe for concurrent use.
type checklist struct {
	em *enumMembers
	m  map[string]struct{} // remaining unsatisfied member names
}

func makeChecklist(em *enumMembers, enumPkg *types.Package, includeUnexported bool, ignore *regexp.Regexp) *checklist {
	m := make(map[string]struct{})

	add := func(memberName string) {
		if memberName == "_" {
			// blank identifier is often used to skip entries in iota lists
			return
		}
		if ignore != nil && ignore.MatchString(enumPkg.Path()+"."+memberName) {
			return
		}
		if !ast.IsExported(memberName) && !includeUnexported {
			return
		}
		m[memberName] = struct{}{}
	}

	for _, name := range em.Names {
		add(name)
	}

	return &checklist{
		em: em,
		m:  m,
	}
}

func (c *checklist) found(memberName string, strategy checkingStrategy) {
	switch strategy {
	case strategyValue:
		// delete all of the same-valued names
		constVal := c.em.NameToValue[memberName]
		for _, n := range c.em.ValueToNames[constVal] {
			delete(c.m, n)
		}

	case strategyName:
		// delete the given name alone
		delete(c.m, memberName)

	default:
		panic(fmt.Sprintf("unknown strategy %v", strategy))
	}
}

func (c *checklist) remaining() map[string]struct{} {
	return c.m
}
