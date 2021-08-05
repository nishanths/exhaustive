package exhaustive

import "testing"

func TestEnumMembers_add(t *testing.T) {
	var v enumMembers
	v.add("foo", nil)
	v.add("z", ptrString("X"))
	v.add("bar", nil)
	v.add("y", ptrString("Y"))
	v.add("x", ptrString("X"))

	checkEqual(t, []string{"foo", "z", "bar", "y", "x"}, v.OrderedNames)
	checkEqual(t, map[string]string{
		"z": "X",
		"y": "Y",
		"x": "X",
	}, v.NameToValue)
	checkEqual(t, map[string][]string{
		"X": []string{"z", "x"},
		"Y": []string{"y"},
	}, v.ValueToNames)
}
