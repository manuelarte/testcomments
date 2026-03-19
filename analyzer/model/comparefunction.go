package model

import (
	"go/ast"
	"go/token"
	"strings"
)

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
func NewCompareFunction(importGroup ImportGroup, funcDecl *ast.FuncDecl) (CompareFunction, bool) {
	booleanCompareFunction, isBooleanCompareFunction := newBooleanCompareFunction(importGroup, funcDecl)
	if isBooleanCompareFunction {
		return booleanCompareFunction, true
	}

	testingsTCompareFunction, isTestingsTCompareFunction := newTestingsTCompareFunction(importGroup, funcDecl)
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

func newBooleanCompareFunction(importGroup ImportGroup, funcDecl *ast.FuncDecl) (BooleanCompareFunction, bool) {
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

		structType, isSame := sameStructType(funcDecl.Type.Params.List[0].Type, funcDecl.Type.Params.List[1].Type)
		if !isSame || structType == "error" {
			return BooleanCompareFunction{}, false
		}

		param1 = funcDecl.Type.Params.List[0].Names[0].Name
		param2 = funcDecl.Type.Params.List[1].Names[0].Name
	default:
		return BooleanCompareFunction{}, false
	}

	if !isComparing(importGroup, funcDecl.Body, param1, param2) {
		return BooleanCompareFunction{}, false
	}

	return BooleanCompareFunction{
		funcDecl: funcDecl,
		param1:   param1,
		param2:   param2,
	}, true
}

func newTestingsTCompareFunction(importGroup ImportGroup, funcDecl *ast.FuncDecl) (TestingsTCompareFunction, bool) {
	if funcDecl.Type.Results != nil {
		return TestingsTCompareFunction{}, false
	}

	params := funcDecl.Type.Params
	if params == nil || len(params.List) < 2 {
		return TestingsTCompareFunction{}, false
	}

	// Check that the first parameter is *testing.T
	if !isTestingTField(params.List[0]) {
		return TestingsTCompareFunction{}, false
	}

	var param1, param2 string

	// Check the remaining parameters (should be 2 more parameters total, or 1 parameter with 2 names)
	switch len(params.List) {
	case 2:
		// Case: t *testing.T, a, b MyStruct (one field with two names)
		param := params.List[1]
		if param.Type == nil {
			return TestingsTCompareFunction{}, false
		}

		if len(param.Names) != 2 {
			return TestingsTCompareFunction{}, false
		}

		param1 = param.Names[0].Name
		param2 = param.Names[1].Name
	case 3:
		// Case: t *testing.T, a MyStruct, b MyStruct (two fields with one name each)
		if len(params.List[1].Names) != 1 || len(params.List[2].Names) != 1 {
			return TestingsTCompareFunction{}, false
		}

		structType, isSame := sameStructType(params.List[1].Type, params.List[2].Type)
		if !isSame || structType == "error" {
			return TestingsTCompareFunction{}, false
		}

		param1 = params.List[1].Names[0].Name
		param2 = params.List[2].Names[0].Name
	default:
		return TestingsTCompareFunction{}, false
	}

	if !isComparing(importGroup, funcDecl.Body, param1, param2) {
		return TestingsTCompareFunction{}, false
	}

	return TestingsTCompareFunction{
		funcDecl: funcDecl,
		param1:   param1,
		param2:   param2,
	}, true
}

//nolint:gocognit
func isComparing(importGroup ImportGroup, block *ast.BlockStmt, param1, param2 string) bool {
	var isLintableComparison, usesCmp bool

	ast.Inspect(block, func(n ast.Node) bool {
		// If we already found a cmp usage, we can stop inspecting.
		if usesCmp {
			return false
		}

		switch node := n.(type) {
		case *ast.CallExpr:
			se, ok := node.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			// Check for cmp.Equal or cmp.Diff
			if importGroup.GoCmp != nil {
				if isGoCmpEqual(importName(importGroup.GoCmp), se) || isGoCmpDiff(importName(importGroup.GoCmp), se) {
					usesCmp = true

					return false // Stop inspection
				}
			}

			// Check for reflect.DeepEqual
			if importGroup.Reflect != nil && isReflectEqual(importName(importGroup.Reflect), se) {
				if len(node.Args) == 2 {
					arg1 := astExprToString(node.Args[0])

					arg2 := astExprToString(node.Args[1])
					if (isParamOrFieldOfParam(arg1, param1) && isParamOrFieldOfParam(arg2, param2)) ||
						(isParamOrFieldOfParam(arg1, param2) && isParamOrFieldOfParam(arg2, param1)) {
						isLintableComparison = true
					}
				}
			}

		case *ast.BinaryExpr:
			// Check for direct comparison == or !=
			if node.Op == token.EQL || node.Op == token.NEQ {
				arg1 := astExprToString(node.X)

				arg2 := astExprToString(node.Y)
				if (isParamOrFieldOfParam(arg1, param1) && isParamOrFieldOfParam(arg2, param2)) ||
					(isParamOrFieldOfParam(arg1, param2) && isParamOrFieldOfParam(arg2, param1)) {
					isLintableComparison = true
				}
			}
		}

		return true
	})

	return isLintableComparison && !usesCmp
}

func isParamOrFieldOfParam(arg, param string) bool {
	return arg == param || strings.HasPrefix(arg, param+".")
}

func sameStructType(type1, type2 ast.Expr) (string, bool) {
	// Convert types to string representation for comparison
	// This is a simple approach - in a real scenario, we'd need to resolve type names
	t1Str := astExprToString(type1)
	t2Str := astExprToString(type2)

	// Both should be non-empty and equal
	if t1Str == "" || t2Str == "" {
		return "", false
	}

	return t1Str, t1Str == t2Str
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
