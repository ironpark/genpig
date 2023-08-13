package parser

import (
	"github.com/samber/lo"
	"go/ast"
	"strings"
)

type FunctionCall struct {
	Package *ImportPackage
	Name    string
	Args    []Arg
}

func (fc *FunctionCall) String() string {
	return fc.Package.Alias + "." + fc.Name + "(" + fc.ArgString() + ")"
}
func (fc *FunctionCall) ArgString() string {
	return strings.Join(lo.Map(fc.Args, func(item Arg, index int) string {
		return item.String()
	}), ",")
}
func NewFunctionCall(exprStmt *ast.ExprStmt, imports Imports) (fc *FunctionCall) {
	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return nil
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return nil
	}
	fc = &FunctionCall{
		Name:    selectorExpr.Sel.Name,
		Package: imports.GetByAlias(ident.Name),
	}
	for _, arg := range callExpr.Args {
		fc.Args = append(fc.Args, NewArg(arg))
	}
	return
}
