package astpain

import (
	"go/ast"
)

// GetDeferredFunctionName возвращает имя функции, вызов которой был отложен через defer,
// если входящий node является *ast.DeferStmt.
func GetDeferredFunctionName(node ast.Node) string {
	deferStmt, ok := node.(*ast.DeferStmt)
	if !ok {
		return ""
	}

	deferCall := deferStmt.Call

	switch funSelector := deferCall.Fun.(type) {
	case *ast.SelectorExpr:
		return getCallableNameAndFunctionName(funSelector)
	case *ast.Ident:
		return funSelector.Name
	case *ast.FuncLit:
		return "anonymous func"
	}

	return ""
}

func getCallableNameAndFunctionName(node *ast.SelectorExpr) string {
	switch nodeX := node.X.(type) {
	case *ast.Ident:
		return nodeX.Name + "." + node.Sel.Name
	case *ast.SelectorExpr:
		return getCallableNameAndFunctionName(nodeX) + "." + node.Sel.Name
	}

	return ""
}
