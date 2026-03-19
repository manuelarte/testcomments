package model

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestNewCompareFunction(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		code      string
		isCompare bool
	}{
		"Simple bool compare": {
			code: `
package main
func Compare(a, b MyStruct) bool {
	return a == b
}
`,
			isCompare: true,
		},
		"Simple bool compare with !=": {
			code: `
package main
func Compare(a, b MyStruct) bool {
	return a != b
}
`,
			isCompare: true,
		},
		"Not compare": {
			code: `
package main
func Compare(a, b MyStruct) bool {
	return true
}
`,
			isCompare: false,
		},
		"Compare with reflect": {
			code: `
package main
import "reflect"
func Compare(a, b MyStruct) bool {
	return reflect.DeepEqual(a, b)
}
`,
			isCompare: true,
		},
		"Compare with cmp.Diff": {
			code: `
package main
import "github.com/google/go-cmp/cmp"
func Compare(a, b MyStruct) bool {
	return cmp.Diff(a, b) == ""
}
`,
			isCompare: false,
		},
		"Compare with cmp.Equal": {
			code: `
package main
import "github.com/google/go-cmp/cmp"
func Compare(a, b MyStruct) bool {
	return cmp.Equal(a, b)
}
`,
			isCompare: false,
		},
		"TestingsT compare": {
			code: `
package main
import "testing"
func Compare(t *testing.T, a, b MyStruct) {
	if a != b {
		t.Errorf("error")
	}
}
`,
			isCompare: true,
		},
		"TestingsT not compare": {
			code: `
package main
import "testing"
func Compare(t *testing.T, a, b MyStruct) {
	t.Log("logging")
}
`,
			isCompare: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fset := token.NewFileSet()

			file, err := parser.ParseFile(fset, "", test.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("parser.ParseFile() error = %v", err)
			}

			var (
				funcDecl    *ast.FuncDecl
				importGroup ImportGroup
			)

			ast.Inspect(file, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.ImportSpec:
					importGroup = importGroup.NewWithImportSpec(x)
				case *ast.FuncDecl:
					funcDecl = x
				}

				return true
			})

			if funcDecl == nil {
				t.Fatal("funcDecl is nil")
			}

			_, got := NewCompareFunction(importGroup, funcDecl)
			if got != test.isCompare {
				t.Errorf("NewCompareFunction() = %v, want %v", got, test.isCompare)
			}
		})
	}
}
