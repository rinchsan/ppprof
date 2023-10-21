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
	analysistest.Run(t, testdata, ppprof.Analyzer, "a")
	analysistest.Run(t, testdata, ppprof.Analyzer, "b")
	analysistest.Run(t, testdata, ppprof.Analyzer, "c")
	analysistest.Run(t, testdata, ppprof.Analyzer, "d")
	analysistest.Run(t, testdata, ppprof.Analyzer, "e")
}
