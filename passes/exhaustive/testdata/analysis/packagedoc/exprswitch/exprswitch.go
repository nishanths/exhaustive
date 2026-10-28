package foo

import "test/analysis/packagedoc/exprswitch/a"
import "test/analysis/packagedoc/exprswitch/b"

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
	switch v { // want "^missing cases in expression switch: darcs$"
	case bazaar:
	case fossil:
	case git:
	case mercurial:
	case subversion:
	}
}

func f2(v a.A) {
	switch v {
	case a.X0:
	case a.X1:
	}

	switch v {
	case b.X0:
	case b.X1:
	}

	switch v { // want "^missing cases in expression switch: b.X1$"
	case a.X0:
	}

	switch v { // want "^missing cases in expression switch: b.X1$"
	case b.X0:
	}
}
