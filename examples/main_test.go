package examples

import "testing"

func TestSomething(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b, want int
	}{
		"simple case": {
			a:    1,
			b:    2,
			want: 3,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := sum(test.a, test.b)
			if got != test.want {
				t.Errorf("got %d, want %d", got, test.want)
			}
		})
	}
}

func sum(a, b int) int {
	return a + b
}
