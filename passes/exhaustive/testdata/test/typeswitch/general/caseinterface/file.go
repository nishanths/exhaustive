package caseinterface

type mn interface{ m(); n() }
type m interface{ m() }
type n interface{ n() }
type o interface{ o() }
type mno interface{ m(); n(); o }
type mne interface{ m(); n(); e() }

type t struct{}
func (t) m()
func (t) n()

type u struct{}
func (*u) m()
func (*u) n()
func (*u) o()

func f1() {
	var x mn

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	}

	switch x.(type) {
	case t:
	case *t:
	case *u:
	}

	switch x.(type) {
	case mn:
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*t$"
	case mno:
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	case mne:
	}

	switch x.(type) {
	case m:
	case n:
	}

	switch x.(type) {
	case interface{ m(); n() }:
	}

	switch x.(type) {
	case interface{ mn }:
	}

	switch x.(type) {
	case interface{ m; n }:
	}

	switch x.(type) {
	case m:
	}

	switch x.(type) {
	case n:
	}

	switch x.(type) {
	case mn:
	case m:
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*t$"
	case o:
	}

	switch x.(type) {
	case t:
	case *t:
	case o:
	}

	switch x.(type) {
	case interface{}:
	}

	switch x.(type) {
	case any:
	}
}

// originally copied from func f5 in testdata package typeswitch/general.
func f2[P n, Q o, R mn, S any]() {
	var x mn

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	}

	switch x.(type) {
	case P:
	}

	switch x.(type) {
	case R:
	}

	switch x.(type) {
	case P:
	case m:
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*t$"
	case Q:
	}

	switch x.(type) {
	case t:
	case *t:
	case Q:
	}

	switch x.(type) {
	case S:
	}

	var y m

	switch y.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	}

	switch y.(type) {
	case R:
	}
}
