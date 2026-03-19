package model

import (
	"go/ast"
)

type (
	// TestFunction is the holder of a test function declaration.
	// A test function must:
	// 1. Start with "Test".
	// 2. Have exactly one parameter.
	// 3. Have that parameter be of type *testing.T.
	TestFunction struct {
		// importGroup contains the import important on this test.
		importGroup ImportGroup

		// testVar is the name given to the testing.T parameter
		testVar string

		// funcDecl the original function declaration.
		funcDecl *ast.FuncDecl

		// tableDrivenInfo table-driven test information for this test function, nil if not a table-driven test.
		tableDrivenInfo *TableDrivenInfo
	}

	// TestedCallExpr contains the actual call to the function tested.
	// got := MyFunction(in) <- TestedCallExpr
	// if got != want {
	//   t.Errorf(...)
	// }.
	TestedCallExpr struct {
		// callExpr contains the actual call to the function tested.
		callExpr *ast.CallExpr

		// params contains all the left hand side params of the assignment
		params []*ast.Ident
	}
)

// NewTestFunction returns a new TestFunction based on the funcDecl.
func NewTestFunction(importGroup ImportGroup, funcDecl *ast.FuncDecl) (TestFunction, bool) {
	ok, testVar := isTestFunction(funcDecl)
	if !ok {
		return TestFunction{}, false
	}

	tbi := newTableDrivenInfo(testVar, funcDecl)

	return TestFunction{
		importGroup:     importGroup,
		testVar:         testVar,
		funcDecl:        funcDecl,
		tableDrivenInfo: tbi,
	}, true
}

func (t TestFunction) ImportGroup() ImportGroup {
	return t.importGroup
}

// GetActualTestBlockStmt returns the actual block test logic, if it's not a table-driven test
// it returns the actual body of the function, and if it's table-driven test it returns
// the content inside the t.Run function.
func (t TestFunction) GetActualTestBlockStmt() *ast.BlockStmt {
	if t.tableDrivenInfo != nil {
		return t.tableDrivenInfo.Block
	}

	return t.funcDecl.Body
}

// TestVar returns the name of the testing.T parameter.
func (t TestFunction) TestVar() string {
	return t.testVar
}

func (t TestFunction) TableDrivenInfo() *TableDrivenInfo {
	return t.tableDrivenInfo
}

// TestPartBlocks returns all the tested blocks of the test function.
func (t TestFunction) TestPartBlocks() []TestPartBlock {
	blStmt := t.GetActualTestBlockStmt()
	testVar := t.TestVar()

	toReturn := make([]TestPartBlock, 0)

	var stmts []ast.Stmt
	if blStmt != nil {
		stmts = blStmt.List
	}

	for i, stmt := range stmts {
		if ifStmt, ok := stmt.(*ast.IfStmt); ok {
			if i == 0 {
				continue
			}

			// the statement should contain the tested function, unless the previous assignment is another if stmt
			// that may contain another testing condition.
			prev := stmts[i-1]
			if _, prevIsIfStmt := prev.(*ast.IfStmt); prevIsIfStmt && i-2 > -1 {
				prev = stmts[i-2]
			}

			testBlock, isTestBlock := NewTestPartBlock(t.ImportGroup(), testVar, prev, ifStmt)
			if !isTestBlock {
				continue
			}

			toReturn = append(toReturn, testBlock)
		}
	}

	return toReturn
}

// NewTestedCallExpr creates a testedFuncStmt after checking that the stmt is a typical function call.
// 1. Statement is an *ast.AssignStmt.
// 2. Right hand side is a *ast.CallExpr
// 3. Left hand side is a list of *ast.Ident, containing the got parameter.
// Can also handle an *ast.IfStmt where the assignment is in the Init field (inlined case).
func NewTestedCallExpr(stmt ast.Stmt) (TestedCallExpr, bool) {
	var callExpr *ast.CallExpr

	params := make([]*ast.Ident, 0)

	// Handle IfStmt with inlined assignment (e.g., if got := func(); got != want)
	if ifStmt, isIfStmt := stmt.(*ast.IfStmt); isIfStmt {
		if ifStmt.Init == nil {
			return TestedCallExpr{}, false
		}

		stmt = ifStmt.Init
	}

	assignStmt, isAssignStmt := stmt.(*ast.AssignStmt)
	if !isAssignStmt {
		return TestedCallExpr{}, false
	}

	if len(assignStmt.Rhs) != 1 {
		return TestedCallExpr{}, false
	}

	for _, expr := range assignStmt.Lhs {
		ident, ok := expr.(*ast.Ident)
		if !ok {
			return TestedCallExpr{}, false
		}

		params = append(params, ident)
	}

	ce, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok {
		return TestedCallExpr{}, false
	}

	callExpr = ce

	return TestedCallExpr{
		callExpr: callExpr,

		params: params,
	}, true
}

func (t TestedCallExpr) CallExpr() *ast.CallExpr {
	return t.callExpr
}

func (t TestedCallExpr) Params() []*ast.Ident {
	return t.params
}

func (t TestedCallExpr) FunctionName() string {
	fn, err := getFunctionName(t.callExpr.Fun)
	if err != nil {
		return ""
	}

	return fn
}

func getFunctionName(expr ast.Expr) (string, error) {
	switch fn := (expr).(type) {
	case *ast.Ident:
		return fn.Name, nil
	case *ast.SelectorExpr:
		value, err := getFunctionName(fn.X)
		if err != nil {
			return "", err
		}

		return value + "." + fn.Sel.Name, nil
	}

	return "", nil
}
