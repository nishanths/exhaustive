package exhaustive

import (
	"fmt"
	"go/ast"
	"go/types"
	"regexp"
)

// A hitlist is the set of enum member names that should be listed in a switch
// statement's case clauses in order for the switch to be exhaustive. The found
// method marks a member as being listed in the switch, so, in usage, a hitlist
// is the set of yet unsatisfied enum members.
//
// Only interact via its methods. It is not safe for concurrent use.
type hitlist struct {
	em *enumMembers
	m  map[string]struct{} // remaining unsatisfied member names
}

func makeHitlist(em *enumMembers, enumPkg *types.Package, includeUnexported bool, ignore *regexp.Regexp) *hitlist {
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

	return &hitlist{
		em: em,
		m:  m,
	}
}

func (h *hitlist) found(memberName string, strategy checkingStrategy) {
	switch strategy {
	case strategyValue:
		if constVal, ok := h.em.NameToValue[memberName]; ok {
			// delete all of the same-valued names
			for _, n := range h.em.ValueToNames[constVal] {
				delete(h.m, n)
			}
		} else {
			// delete the name given name alone
			delete(h.m, memberName)
		}

	case strategyName:
		// delete the given name alone
		delete(h.m, memberName)

	default:
		panic(fmt.Sprintf("unknown strategy %v", strategy))
	}
}

func (h *hitlist) remaining() map[string]struct{} {
	return h.m
}
