// Package plugin add supports to adding this analyzer as a golangci-lint
// plugin.
package plugin

import (
	"fmt"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcomments/analyzer"
)

//nolint:gochecknoinits // init needed for plugin
func init() {
	register.Plugin("testcomments", New)
}

func New(settings any) (register.LinterPlugin, error) {
	castedSettings, err := analyzer.SettingsFrom(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to cast settings: %w", err)
	}

	return &testcommentsPlugin{settings: castedSettings}, nil
}

var _ register.LinterPlugin = new(testcommentsPlugin)

type testcommentsPlugin struct {
	settings analyzer.Settings
}

func (p testcommentsPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.NewWithSettings(p.settings),
	}, nil
}

func (p testcommentsPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
