package ppprof

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "Set up pprof by ppprof"

// Analyzer analyzes the usage of pprof
var Analyzer = &analysis.Analyzer{
	Name: "ppprof",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if n.Name.Name != "main" ||
				n.Recv != nil ||
				len(n.Type.Params.List) != 0 ||
				n.Type.Results != nil {

				return
			}

			if len(n.Body.List) < 3 {
				pass.Reportf(n.Pos(), "should set up pprof at the beginning of main")
				return
			}
		}
	})

	return nil, nil
}
