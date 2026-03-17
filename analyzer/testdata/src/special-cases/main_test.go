package special_cases

import (
	"strings"
	"testing"
)

func TestStringLiteral(t *testing.T) {
	structName := strings.ToLower("Hello")
	if structName != "hello" {
		t.Errorf("strings.ToLower = %q, want 'hello'", structName)
	}
}
