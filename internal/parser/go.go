package parser

import (
	"github.com/samber/lo"
	"go/ast"
	"go/parser"
	"go/token"
)

type GoFile struct {
	PackageName string
	Structs     []*Struct
	Init        struct {
		FuncCalls []*FunctionCall
	}
}

func ParseGoFile(path string) (file GoFile) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return
	}
	file.PackageName = node.Name.Name
	imports := NewImports(node.Imports)
	//top-level declarations parsing
	for i := range node.Decls {
		switch typedNode := node.Decls[i].(type) {
		case *ast.GenDecl:
			for _, spec := range typedNode.Specs {
				switch typedSpec := spec.(type) {
				case *ast.TypeSpec:
					// type declaration
					astStruct, isStructType := typedSpec.Type.(*ast.StructType)
					if !isStructType {
						continue
					}
					file.Structs = append(file.Structs, NewStruct(typedSpec.Name.Name, astStruct, imports))

				case *ast.ValueSpec:
					// const or var declaration
				}
			}
		case *ast.FuncDecl:
			isMethod := typedNode.Recv != nil
			// package init function
			if !isMethod && typedNode.Name.Name == "init" {
				for _, list := range typedNode.Body.List {
					exprStmt, ok := list.(*ast.ExprStmt)
					if ok {
						file.Init.FuncCalls = append(file.Init.FuncCalls, NewFunctionCall(exprStmt, imports))
					}
				}
			}
		}
	}

	for _, st := range file.Structs {
		for _, field := range st.Fields {
			if field.IsNotPrimitiveType() && field.Package == nil {
				typeStruct, ok := lo.Find(file.Structs, func(item *Struct) bool {
					return item.Name == field.Type
				})
				if ok {
					field.Struct = typeStruct
					if field.Name == "" {
						field.Name = field.Type
					}
				}
			}
		}
	}
	return
}
