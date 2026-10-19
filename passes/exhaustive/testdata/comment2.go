// Note: This file is a modified copy of the testdata file in
// package mapliteral/commentassoc.

package a

func f1() {
	type t1 int

	const (
		x0 t1 = iota
		x1
	)

	/* GenDecl */

	var (
		_ = map[t1]bool{x0: true}
		_ = map[t1]bool{x1: true}
	)

	// a0
	var (
		_ = map[t1]bool{x0: true}
		_ = map[t1]bool{x1: true}
	)

	var _ = map[t1]bool{x0: true}

	// a1
	var _ = map[t1]bool{x0: true}

	/* ValueSpec */

	var (
		_ = map[t1]bool{x0: true}
		_ = map[t1]bool{x1: true}
	)
	var (
		// a2
		_ = map[t1]bool{x0: true}
		_ = map[t1]bool{x1: true}
	)

	/* AssignStmt, with assignment operator */

	_ = map[t1]bool{x0: true}

	// a3
	_ = map[t1]bool{x0: true}

	/* AssignStmt, with DEFINE operator */

	tmp1 := map[t1]bool{x0: true}

	// a4
	tmp2 := map[t1]bool{x0: true}
	_ = tmp1
	_ = tmp2

	/* AssignStmt variations */

	f := func(...any) any { return struct{}{} }

	// a5
	_ = map[t1]bool{x0: true}[x0]

	// a6
	_ = &map[t1]bool{x0: true}

	// a7
	_ = f(nil, []int{0}, map[t1]bool{x0: true}, nil)

	/* CallExpr, not supported */

	// foo
	f(nil, []int{0}, map[t1]bool{x0: true}, nil)

	/* ReturnStmt, not supported */

	tmp3 := func() any {
		// foo
		return map[t1]bool{x0: true}
	}
	_ = tmp3
}
