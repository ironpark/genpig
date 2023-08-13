package parser

import (
	"fmt"
	"github.com/samber/lo"
	"go/ast"
	"strings"
)

type Arg interface {
	Type() string
	Value() string
	String() string
}

type ArgValue struct {
	ast *ast.BasicLit
}

func (a ArgValue) Type() string {
	return a.ast.Kind.String()
}

func (a ArgValue) Value() string {
	return a.ast.Value
}

func (a ArgValue) String() string {
	return a.ast.Value
}

type ArgComposite struct {
	ast       *ast.CompositeLit
	keyValues []CompositeKeyValue
}

func (a ArgComposite) Type() string {
	switch t := a.ast.Type.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)
	}
	return ""
}

func (a ArgComposite) Value() string {
	return strings.Join(lo.Map(a.keyValues, func(item CompositeKeyValue, index int) string {
		return item.key + ":" + item.value.String()
	}), ",")
}

func (a ArgComposite) String() string {
	return a.Type() + "{" + a.Value() + "}"
}

type CompositeKeyValue struct {
	key   string
	value Arg
}

func NewComposite(compositeLit *ast.CompositeLit) (composit ArgComposite) {
	composit = ArgComposite{
		ast:       compositeLit,
		keyValues: nil,
	}
	for _, elt := range compositeLit.Elts {
		kvExpr, ok := elt.(*ast.KeyValueExpr)
		if ok {
			ckv := CompositeKeyValue{
				key: kvExpr.Key.(*ast.Ident).Name,
			}
			switch v := kvExpr.Value.(type) {
			case *ast.BasicLit:
				ckv.value = ArgValue{ast: v}
			case *ast.CompositeLit:
				ckv.value = NewComposite(v)
			}
			composit.keyValues = append(composit.keyValues, ckv)
		}
	}
	return
}

func NewArg(arg ast.Expr) (av Arg) {
	switch typedArg := arg.(type) {
	case *ast.BasicLit:
		return ArgValue{typedArg}
	case *ast.CompositeLit:
		return NewComposite(typedArg)
	}
	return
}
