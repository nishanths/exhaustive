package example

type vcs int

const (
	bazaar vcs = iota
	fossil
	git
	mercurial
	subversion
	darcs
)

func f(v vcs) {
	switch v { // want "^switch not exhaustive: missing cases: mercurial, darcs$"
	case bazaar:
	case fossil:
	case git:
	case subversion:
	}
}
