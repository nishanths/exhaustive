package example

type vcs int

const (
	bzr vcs = iota
	fossil
	git
	hg
	svn
	darcs
)

func f(v vcs) {
	switch v { // want "^switch not exhaustive: missing cases: hg, darcs$"
	case bzr:
	case fossil:
	case git:
	case svn:
	}
}
