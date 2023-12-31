package ppprof

import (
	"go/ast"
	"go/types"
	"strconv"

	"github.com/gostaticanalysis/analysisutil"
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

	var fileWithMain *ast.File

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if !isMainFunc(n) {
				return
			}

			fileWithMain = analysisutil.File(pass, n.Pos())

			if !isPprofSetUp(pass, n.Body.List) {
				pass.Report(analysis.Diagnostic{
					Pos:     n.Pos() + 1,
					Message: "should set up pprof at the beginning of main",
					SuggestedFixes: []analysis.SuggestedFix{
						{
							Message: "set up pprof",
							TextEdits: []analysis.TextEdit{
								{
									Pos: n.Body.Lbrace + 1,
									End: n.Body.Lbrace + 1,
									NewText: []byte(`
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	go func() {
		log.Fatal(http.ListenAndServe("localhost:6060", nil))
	}()

`),
								},
							},
						},
					},
				})
				return
			}
		}
	})

	if fileWithMain == nil {
		return nil, nil
	}

	isPprofImported := slices.ContainsFunc(fileWithMain.Imports, func(imp *ast.ImportSpec) bool {
		v, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			return false
		}
		return v == "net/http/pprof"
	})
	if !isPprofImported {
		pass.Report(analysis.Diagnostic{
			Pos:     fileWithMain.Package,
			Message: "should import net/http/pprof",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "import net/http/pprof",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     fileWithMain.Name.End(),
							End:     fileWithMain.Name.End(),
							NewText: []byte("\n" + `import _ "net/http/pprof"`),
						},
					},
				},
			},
		})
	}

	return nil, nil
}

func isMainFunc(decl *ast.FuncDecl) bool {
	return decl.Name.Name == "main" &&
		decl.Recv == nil &&
		len(decl.Type.Params.List) == 0 &&
		decl.Type.Results == nil
}

func callWithFunObjFromStmt(pass *analysis.Pass, stmt ast.Stmt) (*ast.CallExpr, types.Object) {
	expr, _ := stmt.(*ast.ExprStmt)
	if expr == nil {
		return nil, nil
	}
	call, _ := expr.X.(*ast.CallExpr)
	if call == nil {
		return nil, nil
	}
	sel, _ := call.Fun.(*ast.SelectorExpr)
	if sel == nil {
		return nil, nil
	}
	return call, pass.TypesInfo.ObjectOf(sel.Sel)
}

func isRuntimeSetBlockProfileRate(pass *analysis.Pass, stmt ast.Stmt) bool {
	_, obj := callWithFunObjFromStmt(pass, stmt)
	if obj == nil {
		return false
	}

	return obj.Pkg().Path() == "runtime" && obj.Name() == "SetBlockProfileRate"
}

func isRuntimeSetMutexProfileFraction(pass *analysis.Pass, stmt ast.Stmt) bool {
	_, obj := callWithFunObjFromStmt(pass, stmt)
	if obj == nil {
		return false
	}

	return obj.Pkg().Path() == "runtime" && obj.Name() == "SetMutexProfileFraction"
}

func isProfileServed(call *ast.CallExpr, obj types.Object) bool {
	if obj.Pkg().Path() != "net/http" || obj.Name() != "ListenAndServe" {
		return false
	}
	for _, arg := range call.Args {
		blit, _ := arg.(*ast.BasicLit)
		if blit == nil {
			continue
		}
		v, err := strconv.Unquote(blit.Value)
		if err != nil {
			continue
		}
		if v == "localhost:6060" || v == "0.0.0.0:6060" {
			return true
		}
	}
	return false
}

func isProfileServedInGoStmt(pass *analysis.Pass, stmt ast.Stmt) bool {
	goStmt, _ := stmt.(*ast.GoStmt)
	if goStmt == nil {
		return false
	}
	fun, _ := goStmt.Call.Fun.(*ast.FuncLit)
	if fun == nil {
		return false
	}
	for _, stmt := range fun.Body.List {
		call, obj := callWithFunObjFromStmt(pass, stmt)
		if call == nil || obj == nil {
			continue
		}
		if isProfileServed(call, obj) {
			return true
		}
		for _, arg := range call.Args {
			call, _ := arg.(*ast.CallExpr)
			if call == nil {
				continue
			}
			sel, _ := call.Fun.(*ast.SelectorExpr)
			if sel == nil {
				continue
			}
			obj := pass.TypesInfo.ObjectOf(sel.Sel)
			if isProfileServed(call, obj) {
				return true
			}
		}
	}
	return false
}

func isPprofSetUp(pass *analysis.Pass, stmts []ast.Stmt) bool {
	if len(stmts) < 3 {
		return false
	}

	isRuntimeSet := (isRuntimeSetBlockProfileRate(pass, stmts[0]) && isRuntimeSetMutexProfileFraction(pass, stmts[1])) ||
		(isRuntimeSetBlockProfileRate(pass, stmts[1])) && isRuntimeSetMutexProfileFraction(pass, stmts[0])

	isProfileServed := isProfileServedInGoStmt(pass, stmts[2])

	return isRuntimeSet && isProfileServed
}
