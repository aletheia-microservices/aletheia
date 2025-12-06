package utils

import (
	"go/ast"

	"analyzer/pkg/logger"
)

func GetExprType(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.UnaryExpr:
		return e.Op.String()
	default:
		logger.Logger.Fatalf("unknown type in GetExprType for %s: %v", GetType(e), e)
	}
	// FIXME: cover more use cases e.g.
	// 1. type from other package using selector
	// 2. type from current package
	return "unknown"
}
