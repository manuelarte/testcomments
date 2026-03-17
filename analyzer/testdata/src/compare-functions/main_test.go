package compare_functions

import (
	"errors"
	"testing"
)

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

func assertEqual(t *testing.T, a, b MyStruct) {
	t.Helper()

	assertTwoEqual(t, a, b)
}

func assertTwoEqual(t *testing.T, a MyStruct, b MyStruct) {
	t.Helper()

	if a.Name != b.Name || a.Surname != b.Surname {
		t.Errorf("a and b should be equal")
	}
}

func assertErrorIs(err error, expected error) bool {
	return errors.Is(err, expected)
}

func requireErrorIs(t *testing.T, err error, expected error) {
	t.Helper()

	if !errors.Is(err, expected) {
		t.Fatalf("unexpected error got %v, expected %v", err, expected)
	}
}
