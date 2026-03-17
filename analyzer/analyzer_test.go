package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		patterns string
		options  map[string]string
	}{
		"reflect-deepequal": {
			patterns: "reflect-deepequal",
			options: map[string]string{
				IdentifyTheFunctionCHeck: "false",
			},
		},
		"compare functions": {
			patterns: "compare-functions",
		},
		"got before want": {
			patterns: "got-before-want",
			options: map[string]string{
				EqualityComparisonReflectCheckName: "false",
				GotBeforeWantCheck:                 "true",
				IdentifyTheFunctionCHeck:           "false",
			},
		},
		"identify function": {
			patterns: "identify-function",
			options: map[string]string{
				EqualityComparisonReflectCheckName: "false",
			},
		},
		"table-driven test format map-inlined": {
			patterns: "table-driven-testing-format/map-inlined",
			options: map[string]string{
				TableDrivenFormatCheckTypeName:    "map",
				TableDrivenFormatCheckInlinedName: "true",
			},
		},
		"table-driven test format map-non-inlined": {
			patterns: "table-driven-testing-format/map-non-inlined",
			options: map[string]string{
				TableDrivenFormatCheckTypeName:    "map",
				TableDrivenFormatCheckInlinedName: "false",
			},
		},
		"special cases": {
			patterns: "special-cases",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			a := New()

			for k, v := range test.options {
				err := a.Flags.Set(k, v)
				if err != nil {
					t.Fatal(err)
				}
			}

			analysistest.Run(t, analysistest.TestData(), a, test.patterns)
		})
	}
}
