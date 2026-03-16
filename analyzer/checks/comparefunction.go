package checks

import (
	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcomments/analyzer/model"
)

// CompareFunction checks that custom struct comparison functions use cmp.Equal or cmp.Diff.
type CompareFunction struct {
	category string
}

// NewCompareFunction creates a new CompareFunction check.
func NewCompareFunction() CompareFunction {
	return CompareFunction{
		category: "Equality Comparison and Diffs",
	}
}

// Check verifies that struct comparison is using cmp.Equal or cmp.Diff instead of manual field comparison.
func (c CompareFunction) Check(pass *analysis.Pass, compareFunc model.CompareFunction) {
	node := compareFunc.FuncDecl()
	pass.Report(analysis.Diagnostic{
		Pos:      node.Pos(),
		End:      node.End(),
		Category: c.category,
		Message:  "Use cmp.Equal or cmp.Diff for equality comparison",

		URL: "https://github.com/manuelarte/testcomments/tree/main?tab=readme-ov-file#equality-comparison-and-diffs",
	})
}
