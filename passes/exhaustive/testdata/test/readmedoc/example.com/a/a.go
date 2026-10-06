package a

type vcs int

const (
	bazaar vcs = iota
	fossil
	git
	mercurial
	subversion
	darcs
)

func f1(v vcs) {
	switch v { // want "^expression switch not exhaustive: missing cases: mercurial, darcs$"
	case bazaar:
	case fossil:
	case git:
	case subversion:
	}
}
