package exhaustive

import "testing"

func TestEnumsFact(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		e := enumsFact{
			Enums: map[string]*enumMembers{
				"Biome": {Names: []string{"Tundra", "Savanna", "Desert"}},
				"op":    {Names: []string{"_", "add", "sub", "mul", "quotient", "remainder"}},
			},
		}
		if want := "Biome:Tundra,Savanna,Desert; op:_,add,sub,mul,quotient,remainder"; want != e.String() {
			t.Errorf("want %v, got %v", want, e.String())
		}
	})
}
