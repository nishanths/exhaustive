package exhaustive

import "testing"

func TestEnumsFact(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		e := enumsFact{
			Enums: map[string]*enumMembers{
				"Biome": {OrderedNames: []string{"Tundra", "Savanna", "Desert"}},
				"op":    {OrderedNames: []string{"add", "sub", "mul", "quotient", "remainder"}},
			},
		}
		if want := "Biome:Tundra,Savanna,Desert; op:add,sub,mul,quotient,remainder"; want != e.String() {
			t.Errorf("want %v, got %v", want, e.String())
		}
	})
}
