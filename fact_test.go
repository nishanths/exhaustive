package exhaustive

import "testing"

func TestEnumsFact(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		e := enumsFact{
			Enums: map[enumType]*enumMembers{
				{"Biome", "whatever addr"}: {[]string{"Tundra", "Savanna", "Desert"}, nil, nil},
				{"op", "whatever addr"}:    {[]string{"_", "add", "sub", "mul", "quotient", "remainder"}, nil, nil},
			},
		}
		if want := "Biome:Tundra,Savanna,Desert; op:_,add,sub,mul,quotient,remainder"; want != e.String() {
			t.Errorf("want %v, got %v", want, e.String())
		}
	})
}
