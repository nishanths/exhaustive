package exhaustive

import "testing"

func TestEnumMembersFact(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		e := enumMembersFact{
			Members: enumMembers{
				[]string{"Tundra", "Savanna", "Desert"},
				nil,
				nil,
			},
		}
		if want := "Tundra,Savanna,Desert"; want != e.String() {
			t.Errorf("want %v, got %v", want, e.String())
		}

		e = enumMembersFact{
			Members: enumMembers{
				[]string{"_", "add", "sub", "mul", "quotient", "remainder"},
				nil,
				nil,
			},
		}
		if want := "_,add,sub,mul,quotient,remainder"; want != e.String() {
			t.Errorf("want %v, got %v", want, e.String())
		}
	})
}
