package caseinterface

type mn interface {
	m()
	n()
}
type m interface{ m() }
type n interface{ n() }
type o interface{ o() }
type mno interface {
	m()
	n()
	o
}
type mne interface {
	m()
	n()
	e()
}

type t struct{}

func (t) m()
func (t) n()

type u struct{}

func (*u) m()
func (*u) n()
func (*u) o()

func f0() {
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
	case interface {
		m()
		n()
	}:
	}

	switch x.(type) {
	case interface{ mn }:
	}

	switch x.(type) {
	case interface {
		m
		n
	}:
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

func f1[P n, Q o, R mn, S any]() {
	var x mn

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	case P:
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	case R:
	}

	switch x.(type) {
	case P:
	case m:
	}

	switch x.(type) {
	case P:
	case mn:
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	case Q:
	}

	switch x.(type) { // want "^missing cases in type switch: \\*u$"
	case t:
	case *t:
	case Q:
	}

	switch x.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	case S:
	}

	var y m

	switch y.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	}

	switch y.(type) { // want "^missing cases in type switch: t, \\*t, \\*u$"
	case R:
	}
}
