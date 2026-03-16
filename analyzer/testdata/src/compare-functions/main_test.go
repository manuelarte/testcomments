package compare_functions

import "testing"

type MyStruct struct {
	Name, Surname string
}

func TestCompareTwoStructs(t *testing.T) {
	a := MyStruct{Name: "John", Surname: "Doe"}
	b := MyStruct{Name: "Alice", Surname: "Doe"}

	if areEqual(a, b) {
		t.Errorf("a and b should not be equal")
	}
}

func TestAssertTwoStructs(t *testing.T) {
	a := MyStruct{Name: "John", Surname: "Doe"}
	b := MyStruct{Name: "Alice", Surname: "Doe"}

	assertEqual(t, a, b)
}

func areEqual(a, b MyStruct) bool { // want `Use cmp.Equal or cmp.Diff for equality comparison`
	return areTwoEqual(a, b)
}

func areTwoEqual(a MyStruct, b MyStruct) bool { // want `Use cmp.Equal or cmp.Diff for equality comparison`
	return a.Name == b.Name && a.Surname == b.Surname
}

func assertEqual(t *testing.T, a, b MyStruct) { // want `Use cmp.Equal or cmp.Diff for equality comparison`
	assertTwoEqual(t, a, b)
}

func assertTwoEqual(t *testing.T, a MyStruct, b MyStruct) { // want `Use cmp.Equal or cmp.Diff for equality comparison`
	if a.Name != b.Name || a.Surname != b.Surname {
		t.Errorf("a and b should be equal")
	}
}
