package model

import "go/ast"

var (
	_ CompareFunction = new(BooleanCompareFunction)
	_ CompareFunction = new(TestingsTCompareFunction)
)

type (
	// CompareFunction holds function that is used to compare two structs
	// the signature can be
	//   - myFunction(t *testing.T, x MyStruct, y MyStruct)
	//   - myFunction(t *testing.T, x, y MyStruct)
	//   - myFunction(x MyStruct, y MyStruct) bool
	//   - myFunction(x, y MyStruct) bool
	CompareFunction interface {
		FuncDecl() *ast.FuncDecl
		Param1() string
		Param2() string
	}

	// BooleanCompareFunction holds compare functions that return bool
	//   - myFunction(x MyStruct, y MyStruct) bool
	//   - myFunction(x, y MyStruct) bool
	BooleanCompareFunction struct {
		// funcDecl the original function declaration.
		funcDecl *ast.FuncDecl
		param1   string
		param2   string
	}

	// TestingsTCompareFunction holds compare functions that have testing.T as first parameter
	//   - myFunction(t *testing.T, x MyStruct, y MyStruct)
	//   - myFunction(t *testing.T, x, y MyStruct)
	TestingsTCompareFunction struct {
		// funcDecl the original function declaration.
		funcDecl *ast.FuncDecl
		param1   string
		param2   string
	}
)

// NewCompareFunction returns a new CompareFunction based on the funcDecl.
// It detects functions that compare two structs by checking the signature.
func NewCompareFunction(_ ImportGroup, funcDecl *ast.FuncDecl) (CompareFunction, bool) {
	booleanCompareFunction, isBooleanCompareFunction := newBooleanCompareFunction(funcDecl)
	if isBooleanCompareFunction {
		return booleanCompareFunction, true
	}

	testingsTCompareFunction, isTestingsTCompareFunction := newTestingsTCompareFunction(funcDecl)
	if isTestingsTCompareFunction {
		return testingsTCompareFunction, true
	}

	return nil, false
}

func (b BooleanCompareFunction) FuncDecl() *ast.FuncDecl {
	return b.funcDecl
}

func (b BooleanCompareFunction) Param1() string {
	return b.param1
}

func (b BooleanCompareFunction) Param2() string {
	return b.param2
}

func (t TestingsTCompareFunction) FuncDecl() *ast.FuncDecl {
	return t.funcDecl
}

func (t TestingsTCompareFunction) Param1() string {
	return t.param1
}

func (t TestingsTCompareFunction) Param2() string {
	return t.param2
}

func newBooleanCompareFunction(funcDecl *ast.FuncDecl) (BooleanCompareFunction, bool) {
	if funcDecl.Type.Results == nil {
		return BooleanCompareFunction{}, false
	}

	outputs := funcDecl.Type.Results.List
	if len(outputs) != 1 || !isBoolType(outputs[0].Type) {
		return BooleanCompareFunction{}, false
	}

	var param1, param2 string

	switch len(funcDecl.Type.Params.List) {
	case 1:
		// expecting a, b MyStruct
		param := funcDecl.Type.Params.List[0]
		if param.Type == nil {
			return BooleanCompareFunction{}, false
		}

		if len(param.Names) != 2 {
			return BooleanCompareFunction{}, false
		}

		param1 = param.Names[0].Name
		param2 = param.Names[1].Name
	case 2:
		if len(funcDecl.Type.Params.List[0].Names) != 1 || len(funcDecl.Type.Params.List[1].Names) != 1 {
			return BooleanCompareFunction{}, false
		}

		if !isSameStructType(funcDecl.Type.Params.List[0].Type, funcDecl.Type.Params.List[1].Type) {
			return BooleanCompareFunction{}, false
		}

		param1 = funcDecl.Type.Params.List[0].Names[0].Name
		param2 = funcDecl.Type.Params.List[1].Names[0].Name
	default:
		return BooleanCompareFunction{}, false
	}

	return BooleanCompareFunction{
		funcDecl: funcDecl,
		param1:   param1,
		param2:   param2,
	}, true
}

func newTestingsTCompareFunction(funcDecl *ast.FuncDecl) (TestingsTCompareFunction, bool) {
	// TODO
	return TestingsTCompareFunction{}, false
}

// isSameStructType checks if two types are the same and are struct types (or named types that could be structs).
func isSameStructType(type1, type2 ast.Expr) bool {
	// Convert types to string representation for comparison
	// This is a simple approach - in a real scenario, we'd need to resolve type names
	t1Str := astExprToString(type1)
	t2Str := astExprToString(type2)

	// Both should be non-empty and equal
	if t1Str == "" || t2Str == "" {
		return false
	}

	return t1Str == t2Str
}

// isBoolType checks if an expression represents the bool type.
func isBoolType(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "bool"
	}

	return false
}

// astExprToString converts an AST expression to its string representation.
func astExprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return astExprToString(e.X) + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + astExprToString(e.X)
	case *ast.ArrayType:
		return "[" + astExprToString(e.Len) + "]" + astExprToString(e.Elt)
	case *ast.MapType:
		return "map[" + astExprToString(e.Key) + "]" + astExprToString(e.Value)
	default:
		return ""
	}
}
