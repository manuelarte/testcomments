// Package analyzer contains the analyzer with the business logic of this linter.
package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/manuelarte/testcomments/analyzer/checks"
	"github.com/manuelarte/testcomments/analyzer/model"
)

type (
	testcomments struct {
		settings Settings
	}
)

func New() *analysis.Analyzer {
	l := testcomments{
		settings: DefaultSettings(),
	}

	a := &analysis.Analyzer{
		Name:     "testcomments",
		Doc:      "checks test follow standards",
		URL:      "https://github.com/manuelarte/testcomments",
		Run:      l.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	a.Flags.BoolVar(&l.settings.EqualityComparison.Reflect, EqualityComparisonReflectCheckName, true,
		"Checks reflect.DeepEqual can be replaced by newer cmp.Equal or cmp.Diff.")
	a.Flags.BoolVar(&l.settings.EqualityComparison.Equal, EqualityComparisonEqualCheckName, true,
		"Checks helper comparing functions can be replaced by cmp.Equal or cmp.Diff.")
	a.Flags.BoolVar(&l.settings.GotBeforeWant, GotBeforeWantCheck, true,
		"Check that output the actual value that the function returned before printing the value that was expected.")
	a.Flags.BoolVar(&l.settings.IdentifyFunction, IdentifyTheFunctionCHeck, true,
		"Check that the failure messages in t.Errorf contains the function name.")
	a.Flags.StringVar(&l.settings.TableDrivenFormat.FormatType, TableDrivenFormatCheckTypeName, "",
		"Check that the table-driven tests are either Map or Slice.")
	a.Flags.BoolVar(&l.settings.TableDrivenFormat.Inlined, TableDrivenFormatCheckInlinedName, false,
		"Check that the table-driven tests are either inline or declared before.")

	return a
}

func NewWithSettings(settings Settings) *analysis.Analyzer {
	l := testcomments{
		settings: settings,
	}

	return &analysis.Analyzer{
		Name:     "testcomments",
		Doc:      "checks test follow standards",
		URL:      "https://github.com/manuelarte/testcomments",
		Run:      l.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

//nolint:gocognit // refactor later
func (l *testcomments) run(pass *analysis.Pass) (any, error) {
	insp, found := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !found {
		//nolint:nilnil // impossible case.
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.ImportSpec)(nil),
		(*ast.FuncDecl)(nil),
	}

	var importGroup model.ImportGroup

	var (
		tbfCheck         = checks.NewTableDrivenFormat(l.settings.TableDrivenFormat.getTableDrivenFormatPredicate())
		compareFunction  = checks.NewCompareFunction()
		reflectDeepEqual = checks.NewReflectDeepEqual()
		gotBeforeWant    = checks.NewGotBeforeWant()
		identifyFunction = checks.NewIdentifyFunction()
	)

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Only process _test.go files
		if !strings.HasSuffix(pass.Fset.File(n.Pos()).Name(), "_test.go") {
			importGroup = model.ImportGroup{}

			return
		}

		switch node := n.(type) {
		case *ast.ImportSpec:
			importGroup = importGroup.NewWithImportSpec(node)
		case *ast.FuncDecl:
			if l.settings.EqualityComparison.Equal {
				if compareFunc, isCompareFunc := model.NewCompareFunction(importGroup, node); isCompareFunc {
					compareFunction.Check(pass, compareFunc)

					return
				}
			}

			if testFunc, ok := model.NewTestFunction(importGroup, node); ok {
				tbfCheck.Check(pass, testFunc)

				if l.settings.EqualityComparison.Reflect {
					reflectDeepEqual.Check(pass, testFunc)
				}

				if l.settings.GotBeforeWant {
					gotBeforeWant.Check(pass, testFunc)
				}

				if l.settings.IdentifyFunction {
					identifyFunction.Check(pass, testFunc)
				}

				return
			}
		}
	})

	//nolint:nilnil //any, error
	return nil, nil
}
