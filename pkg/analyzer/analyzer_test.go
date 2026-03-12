package analyzer_test

import (
	"testing"

	"github.com/BORONINN/go-log-linter/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	testCases := []struct {
		name string
		pkgs []string
	}{
		{
			name: "slog cases",
			pkgs: []string{"slogcases"},
		},
		{
			name: "zap cases",
			pkgs: []string{"zapcases"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analysistest.Run(t, testdata, analyzer.Analyzer, tc.pkgs...)
		})
	}
}
