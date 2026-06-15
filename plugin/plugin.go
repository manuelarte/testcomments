// Package plugin add supports to adding this analyzer as a golangci-lint
// plugin.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcomments/analyzer"
)

//nolint:gochecknoinits // init needed for plugin
func init() {
	register.Plugin("testcomments", New)
}

func New(_ any) (register.LinterPlugin, error) {
	return &testcommentsPlugin{}, nil
}

var _ register.LinterPlugin = new(testcommentsPlugin)

type testcommentsPlugin struct{}

func (u testcommentsPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.New(),
	}, nil
}

func (u testcommentsPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
