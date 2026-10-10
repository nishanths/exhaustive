package invalidcode

func f6() {
	type n interface{}
	type s struct{ *n } // invalid: embedded field type cannot be a pointer to an interface
}
