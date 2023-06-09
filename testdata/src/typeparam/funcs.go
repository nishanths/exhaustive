package typeparam

import y "general/y"

func _a[T y.Phylum | M](v T) {
	switch v { // want `^missing cases in switch of type bar.Phylum\|typeparam.M: bar.Chordata, bar.Mollusca, typeparam.B$`
	case T(A):
	case T(y.Echinodermata):
	}

	switch M(v) { // want "^missing cases in switch of type typeparam.M: typeparam.A, typeparam.B$"
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type bar.Phylum\|typeparam.M: bar.Chordata, bar.Mollusca, typeparam.B$`
		T(A):               struct{}{},
		T(y.Echinodermata): struct{}{},
	}
}

func _b[T N | MM](v T) {
	switch v { // want `^missing cases in switch of type typeparam.N\|typeparam.M: typeparam.D\|typeparam.B$`
	case T(A):
	}

	switch M(v) { // want "^missing cases in switch of type typeparam.M: typeparam.A, typeparam.B$"
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.N\|typeparam.M: typeparam.D\|typeparam.B$`
		T(A): struct{}{},
	}
}

func _c0[T O | M | N](v T) {
	switch v { // want `^missing cases in switch of type typeparam.O\|typeparam.M\|typeparam.N: typeparam.E1, typeparam.E2, typeparam.A\|typeparam.C$`
	case T(B):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.O\|typeparam.M\|typeparam.N: typeparam.E1, typeparam.E2, typeparam.A\|typeparam.C$`
		T(B): struct{}{},
	}
}

func _c2[T interface{ O } | M | interface{ N }](v T) {
	switch v { // want `^missing cases in switch of type typeparam.O\|typeparam.M\|typeparam.N: typeparam.E1, typeparam.E2, typeparam.A\|typeparam.C$`
	case T(B):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.O\|typeparam.M\|typeparam.N: typeparam.E1, typeparam.E2, typeparam.A\|typeparam.C$`
		T(B): struct{}{},
	}
}

func _d[T y.Phylum | II | M](v T) {
	switch v { // want `^missing cases in switch of type bar.Phylum\|typeparam.N\|typeparam.O\|typeparam.M: bar.Chordata, bar.Mollusca, typeparam.D\|typeparam.B, typeparam.E1, typeparam.E2$`
	case T(A):
	case T(y.Echinodermata):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type bar.Phylum\|typeparam.N\|typeparam.O\|typeparam.M: bar.Chordata, bar.Mollusca, typeparam.D\|typeparam.B, typeparam.E1, typeparam.E2$`
		T(A):               struct{}{},
		T(y.Echinodermata): struct{}{},
	}
}

func _e[T M](v T) {
	switch v { // want `^missing cases in switch of type typeparam.M: typeparam.A$`
	case T(M(B)):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.M: typeparam.A$`
		T(M(B)): struct{}{},
	}
}

func _f[T Anon](v T) {
	switch v { // want `^missing cases in switch of type typeparam.M\|typeparam.N: typeparam.B\|typeparam.D$`
	case T(C):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.M\|typeparam.N: typeparam.B\|typeparam.D$`
		T(C): struct{}{},
	}
}

func repeat0[T II | O](v T) {
	switch v { // want `^missing cases in switch of type typeparam.N\|typeparam.O: typeparam.C, typeparam.D, typeparam.E2$`
	case T(E1):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.N\|typeparam.O: typeparam.C, typeparam.D, typeparam.E2$`
		T(E1): struct{}{},
	}
}

func repeat1[T MM | M](v T) {
	switch v { // want `^missing cases in switch of type typeparam.M: typeparam.A$`
	case T(B):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.M: typeparam.A$`
		T(B): struct{}{},
	}
}

func repeat2[T interface{ M } | interface{ M }](v T) {
	switch v { // want `^missing cases in switch of type typeparam.M: typeparam.A$`
	case T(B):
	}

	_ = map[T]struct{}{ // want `^missing keys in map of key type typeparam.M: typeparam.A$`
		T(B): struct{}{},
	}
}

func _mixedTypes0[T M | QQ](v T) {
	// expect no diagnostic because underlying basic kinds are not same:
	// uint8 vs. string
	switch v {
	case T(A):
	}

	_ = map[T]struct{}{
		T(A): struct{}{},
	}
}

func _mixedTypes1[T MM | QQ](v T) {
	switch v {
	case T(A):
	}

	_ = map[T]struct{}{
		T(A): struct{}{},
	}
}

func _mixedTypes2[T interface{ M } | interface{ Q }](v T) {
	switch v {
	case T(A):
	}

	_ = map[T]struct{}{
		T(A): struct{}{},
	}
}

func _mixedTypes3[T interface{ M | Q }](v T) {
	switch v {
	case T(A):
	}

	_ = map[T]struct{}{
		T(A): struct{}{},
	}
}

func _noCheck0[T M | NotEnumType](v T) {
	// expect no diagnostic because not type elements are enum types.
	switch v {
	case T(B):
	}

	_ = map[T]struct{}{
		T(B): struct{}{},
	}
}

func _noCheck1[T LL](v T) {
	switch v {
	case T(A):
	}

	_ = map[T]struct{}{
		T(A): struct{}{},
	}
}

func _noCheck2[T KK](v T) {
	switch v {
	case T(A):
	}

	_ = map[T]struct{}{
		T(A): struct{}{},
	}
}
