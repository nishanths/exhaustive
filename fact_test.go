package exhaustive

import "testing"

func TestEnumMembersFact(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		e := enumMembersFact{
			Members: enumMembers{
				Names: []string{"Tundra", "Savanna", "Desert"},
				NameToValue: map[string]constantValue{
					"Tundra":  "1",
					"Savanna": "2",
					"Desert":  "3",
				},
				ValueToNames: map[constantValue][]string{
					"1": {"Tundra"},
					"2": {"Savanna"},
					"3": {"Desert"},
				},
			},
		}
		checkEnumMembersLiteral(t, "Biome", e.Members)
		if want := "Tundra,Savanna,Desert"; want != e.String() {
			t.Errorf("want %v, got %v", want, e.String())
		}

		e = enumMembersFact{
			Members: enumMembers{
				Names: []string{"_", "add", "sub", "mul", "quotient", "remainder"},
				NameToValue: map[string]constantValue{
					"_":         "0",
					"add":       "1",
					"sub":       "2",
					"mul":       "3",
					"quotient":  "3",
					"remainder": "3",
				},
				ValueToNames: map[constantValue][]string{
					"0": {"_"},
					"1": {"add"},
					"2": {"sub"},
					"3": {"mul", "quotient", "remainder"},
				},
			},
		}
		checkEnumMembersLiteral(t, "Token", e.Members)
		if want := "_,add,sub,mul,quotient,remainder"; want != e.String() {
			t.Errorf("want %v, got %v", want, e.String())
		}
	})
}
