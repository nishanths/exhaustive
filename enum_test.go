package exhaustive

import (
	"reflect"
	"testing"
)

func TestEnumMembers_add(t *testing.T) {
	var v enumMembers
	v.add("foo", nil)
	v.add("z", ptrString("X"))
	v.add("bar", nil)
	v.add("y", ptrString("Y"))
	v.add("x", ptrString("X"))

	if want, got := []string{"foo", "z", "bar", "y", "x"}, v.OrderedNames; !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
	if want, got := map[string]string{
		"z": "X",
		"y": "Y",
		"x": "X",
	}, v.NameToValue; !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}

	if want, got := map[string][]string{
		"X": {"z", "x"},
		"Y": {"y"},
	}, v.ValueToNames; !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}
