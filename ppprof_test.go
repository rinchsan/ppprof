package ppprof_test

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/rinchsan/ppprof"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.RunWithSuggestedFixes(t, testdata, ppprof.Analyzer, "a")
	analysistest.RunWithSuggestedFixes(t, testdata, ppprof.Analyzer, "b")
	analysistest.RunWithSuggestedFixes(t, testdata, ppprof.Analyzer, "c")
	analysistest.RunWithSuggestedFixes(t, testdata, ppprof.Analyzer, "d")
	analysistest.RunWithSuggestedFixes(t, testdata, ppprof.Analyzer, "e")
	analysistest.RunWithSuggestedFixes(t, testdata, ppprof.Analyzer, "main")
}
