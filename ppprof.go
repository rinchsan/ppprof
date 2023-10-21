package ppprof

import (
	"go/ast"
	"go/types"

	"golang.org/x/exp/slices"
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

			if !isPprofSetUp(pass, n.Body.List) {
				pass.Reportf(n.Pos(), "should set up pprof at the beginning of main")
				return
			}
		}
	})

	return nil, nil
}

func objFromStmt(pass *analysis.Pass, stmt ast.Stmt) types.Object {
	expr, _ := stmt.(*ast.ExprStmt)
	if stmt == nil {
		return nil
	}
	call, _ := expr.X.(*ast.CallExpr)
	if call == nil {
		return nil
	}
	sel, _ := call.Fun.(*ast.SelectorExpr)
	if sel == nil {
		return nil
	}
	return pass.TypesInfo.ObjectOf(sel.Sel)
}

func isRuntimeSetBlockProfileRate(pass *analysis.Pass, stmt ast.Stmt) bool {
	obj := objFromStmt(pass, stmt)
	return obj.Pkg().Path() == "runtime" && obj.Name() == "SetBlockProfileRate"
}

func isRuntimeSetMutexProfileFraction(pass *analysis.Pass, stmt ast.Stmt) bool {
	obj := objFromStmt(pass, stmt)
	return obj.Pkg().Path() == "runtime" && obj.Name() == "SetMutexProfileFraction"
}

func isProfileServedInGoStmt(stmt ast.Stmt) bool {
	return true
}

func isPprofSetUp(pass *analysis.Pass, stmts []ast.Stmt) bool {
	if len(stmts) < 3 {
		return false
	}

	isPprofImported := slices.ContainsFunc(pass.Pkg.Imports(), func(pkg *types.Package) bool {
		return pkg.Path() == "net/http/pprof"
	})

	isRuntimeSet := (isRuntimeSetBlockProfileRate(pass, stmts[0]) && isRuntimeSetMutexProfileFraction(pass, stmts[1])) ||
		(isRuntimeSetBlockProfileRate(pass, stmts[1])) && isRuntimeSetMutexProfileFraction(pass, stmts[0])

	isProfileServed := isProfileServedInGoStmt(stmts[2])

	return isPprofImported && isRuntimeSet && isProfileServed
}
