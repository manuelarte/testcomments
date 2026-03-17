package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func double(a int) int {
	return a * 2
}

func TestDouble(t *testing.T) {
	t.Parallel()

	expected := 2
	actual := double(1)
	if expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual) // want `Test outputs should output the actual value that the function returned before printing the value that was expected`
	}
}

func TestTableDrivenDouble(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want int
	}{
		{
			name: "simple case",
			want: 2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := double(1)
			if got != test.want {
				t.Errorf("want %v, got %v", test.want, got) // want `Test outputs should output the actual value that the function returned before printing the value that was expected`
			}
		})
	}
}

func TestTableDrivenInlinedDouble(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want int
	}{
		{
			name: "simple case",
			want: 2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := double(1); got != test.want {
				t.Errorf("want %v, got %v", test.want, got) // want `Test outputs should output the actual value that the function returned before printing the value that was expected`
			}
		})
	}
}

func TestShouldNotGetTriggered(t *testing.T) {
	want := []string{"John", "Doe"}
	got := splitString("John Doe")

	if diff := cmp.Diff(got, want, ""); diff != "" {
		t.Errorf("diff (-got +want):\n%s", diff)
	}
}

func splitString(fullName string) []string {
	return strings.Split(fullName, " ")
}
