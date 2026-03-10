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

func areEqual(a, b MyStruct) bool {
	return a.Name == b.Name && a.Surname == b.Surname // want `Use cmp.Equal or cmp.Diff for equality comparison`
}
