package checks

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcomments/analyzer/model"
)

// ReflectDeepEqual checks that reflect.DeepEqual can be replaced by newer cmp.Equal.
type ReflectDeepEqual struct {
	category string
}

// NewReflectDeepEqual creates a new EqualityComparison.
func NewReflectDeepEqual() ReflectDeepEqual {
	return ReflectDeepEqual{
		category: "Equality Comparison and Diffs",
	}
}

//nolint:gocritic // still under development
func (c ReflectDeepEqual) Check(pass *analysis.Pass, testFunc model.TestFunction) {
	reflectImportName, ok := testFunc.ImportGroup().ReflectImportName()
	if !ok {
		return
	}

	blStmt := testFunc.GetActualTestBlockStmt()

	var stmts []ast.Stmt
	if blStmt != nil {
		stmts = blStmt.List
	}

	for _, stmt := range stmts {
		switch node := stmt.(type) {
		case *ast.IfStmt:
			// check reflect.DeepEqual calls
			diag := c.checkCond(node.Cond, reflectImportName)
			if diag != nil {
				pass.Report(*diag)
			}
		}
	}
}

func (c ReflectDeepEqual) checkCond(cond ast.Expr, reflectImportName string) *analysis.Diagnostic {
	switch node := cond.(type) {
	case *ast.CallExpr:
		return c.checkCallExpr(node, reflectImportName)
	case *ast.UnaryExpr:
		return c.checkUnaryExpr(node, reflectImportName)
	}

	return nil
}

//nolint:gocritic // still under development
func (c ReflectDeepEqual) checkUnaryExpr(unary *ast.UnaryExpr, reflectImportName string) *analysis.Diagnostic {
	switch node := unary.X.(type) {
	case *ast.CallExpr:
		// check reflect.DeepEqual
		return c.checkCallExpr(node, reflectImportName)
	}

	return nil
}

//nolint:gocritic // still under development
func (c ReflectDeepEqual) checkCallExpr(call *ast.CallExpr, reflectImportName string) *analysis.Diagnostic {
	switch node := call.Fun.(type) {
	case *ast.SelectorExpr:
		if ident, ok := node.X.(*ast.Ident); ok && ident.Name == reflectImportName && node.Sel.Name == "DeepEqual" {
			return &analysis.Diagnostic{
				Pos:      node.Pos(),
				End:      node.End(),
				Category: c.category,
				Message:  "Use cmp.Equal or cmp.Diff for equality comparison",

				URL: "https://github.com/manuelarte/testcomments/tree/main?tab=readme-ov-file#equality-comparison-and-diffs",
			}
		}
	}

	return nil
}
