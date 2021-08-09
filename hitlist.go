package exhaustive

import (
	"go/ast"
	"go/types"
	"regexp"
)

// A hitlist is the set of enum member names that should be listed in a switch
// statement's case clauses in order for the switch to be exhaustive. In usage,
// it is the set of yet unsatisfied enum members.
type hitlist struct {
	em *enumMembers
	m  map[string]struct{}
}

func (h *hitlist) found(foundName string) {
	if h.len() == 0 {
		return
	}

	constVal, ok := h.em.NameToValue[foundName]
	if !ok {
		// only delete the name alone from hitlist
		delete(h.m, foundName)
		return
	}

	// delete all of the same-valued names from hitlist
	namesToDelete := h.em.ValueToNames[constVal]
	for _, n := range namesToDelete {
		delete(h.m, n)
	}
}

func (h *hitlist) len() int {
	return len(h.m)
}

func (h *hitlist) remaining() map[string]struct{} {
	return h.m
}

func makeHitlist(em *enumMembers, enumPkg *types.Package, checkUnexported bool, ignore *regexp.Regexp) *hitlist {
	ret := hitlist{
		em: em,
		m:  make(map[string]struct{}),
	}
	for _, name := range em.OrderedNames {
		if name == "_" {
			// blank identifier is often used to skip entries in iota lists
			continue
		}
		if ignore != nil && ignore.MatchString(enumPkg.Path()+"."+name) {
			continue
		}
		if !ast.IsExported(name) && !checkUnexported {
			continue
		}
		ret.m[name] = struct{}{}
	}
	return &ret
}
