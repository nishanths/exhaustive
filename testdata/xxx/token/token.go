package token

type Token int

const (
	Add Token = iota
	Subtract
	Multiply
	Quotient
	Remainder
)
