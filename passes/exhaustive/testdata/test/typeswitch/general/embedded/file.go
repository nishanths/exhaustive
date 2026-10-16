package embedded

type t struct{}
func (t) m()

type u struct{}
func (*u) n()

type s1 struct{ u }
type s2 struct{ *u }
type s3 struct{ t; u }
type s4 struct{ *t; *u }

type io interface{ o() }
type e1 struct{ io }

func f1() {
	switch (interface{ n() }(nil)).(type) { // want "^missing cases in type switch: \\*s1, s2, \\*s2, \\*s3, s4, \\*s4, \\*u$"
	}

	switch (interface{ n() }(nil)).(type) {
	case *s1:
	case s2:
	case *s2:
	case *s3:
	case s4:
	case *s4:
	case *u:
	}
}

func f2() {
	switch (interface{ m(); n() }(nil)).(type) { // want "^missing cases in type switch: \\*s3, s4, \\*s4$"
	}

	switch (interface{ m(); n() }(nil)).(type) {
	case *s3:
	case s4:
	case *s4:
	}
}

func f3() {
	switch (interface{ o() }(nil)).(type) { // want "^missing cases in type switch: e1, \\*e1$"
	}

	switch (interface{ o() }(nil)).(type) {
	case e1:
	case *e1:
	}
}
