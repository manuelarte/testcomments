package model

import "go/ast"

// CompareFunction holds function that is used to compare two structs
// the signature can be
//   - myFunction(t *testing.T, x MyStruct, y MyStruct)
//   - myFunction(t *testing.T, x, y MyStruct)
//   - myFunction(x MyStruct, y MyStruct) bool
//   - myFunction(x, y MyStruct) bool
type CompareFunction struct {
	// funcDecl the original function declaration.
	funcDecl *ast.FuncDecl
}

// NewCompareFunction returns a new CompareFunction based on the funcDecl.
func NewCompareFunction(importGroup ImportGroup, funcDecl *ast.FuncDecl) (CompareFunction, bool) {
	// TODO: implement
	return CompareFunction{
		funcDecl: funcDecl,
	}, true
}
