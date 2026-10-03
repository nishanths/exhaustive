package commentassoc

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

	//exhaustive:enforce
	var (
		_ = map[t1]bool{x0: true} // want "^map literal not exhaustive: missing keys: x1$"
		_ = map[t1]bool{x1: true} // want "^map literal not exhaustive: missing keys: x0$"
	)

	var _ = map[t1]bool{x0: true}

	//exhaustive:enforce
	var _ = map[t1]bool{x0: true} // want "^map literal not exhaustive: missing keys: x1$"

	/* ValueSpec */

	var (
		_ = map[t1]bool{x0: true}
		_ = map[t1]bool{x1: true}
	)
	var (
		//exhaustive:enforce
		_ = map[t1]bool{x0: true} // want "^map literal not exhaustive: missing keys: x1$"
		_ = map[t1]bool{x1: true}
	)

	/* AssignStmt, with assignment operator */

	_ = map[t1]bool{x0: true}

	//exhaustive:enforce
	_ = map[t1]bool{x0: true} // want "^map literal not exhaustive: missing keys: x1$"

	/* AssignStmt, with DEFINE operator */

	tmp1 := map[t1]bool{x0: true}

	//exhaustive:enforce
	tmp2 := map[t1]bool{x0: true} // want "^map literal not exhaustive: missing keys: x1$"
	_ = tmp1
	_ = tmp2

	/* AssignStmt variations */

	f := func(...any) any { return struct{}{} }

	//exhaustive:enforce
	_ = map[t1]bool{x0: true}[x0] // want "^map literal not exhaustive: missing keys: x1$"

	//exhaustive:enforce
	_ = &map[t1]bool{x0: true} // want "^map literal not exhaustive: missing keys: x1$"

	//exhaustive:enforce
	_ = f(nil, []int{0}, map[t1]bool{x0: true}, nil) // want "^map literal not exhaustive: missing keys: x1$"

	/* CallExpr, not supported */

	//exhaustive:enforce
	f(nil, []int{0}, map[t1]bool{x0: true}, nil)

	/* ReturnStmt, not supported */

	tmp3 := func() any {
		//exhaustive:enforce
		return map[t1]bool{x0: true}
	}
	_ = tmp3
}
