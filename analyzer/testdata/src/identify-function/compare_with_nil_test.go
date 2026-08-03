package main

import (
	"errors"
	"testing"
)

func TestDivideNotInFailureMessage(t *testing.T) {
	t.Parallel()
	a, b := 2, 3
	_, err := divide(a, b)
	if err != nil {
		t.Errorf("unexpected error: %v", err) // want `Failure messages should include the name of the function that failed`
	}
}

func TestDivideInFailureMessage(t *testing.T) {
	t.Parallel()
	a, b := 2, 3
	_, err := divide(a, b)
	if err != nil {
		t.Errorf("divide, unexpected error: %v", err)
	}
}

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}

	return a / b, nil
}
