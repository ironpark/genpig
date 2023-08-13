package parser

import (
	"fmt"
	"github.com/samber/lo"
	"go/ast"
	"strings"
)

type ImportPackage struct {
	Alias       string
	PackageName string
	Path        string
}

func (p *ImportPackage) String() string {
	if p.Alias == p.PackageName {
		return fmt.Sprintf(`"%s"`, p.Path)
	}
	return fmt.Sprintf(`%s "%s"`, p.Alias, p.Path)
}

func NewImportPackage(importSpec *ast.ImportSpec) (im *ImportPackage) {
	im = &ImportPackage{
		Path: strings.ReplaceAll(importSpec.Path.Value, "\"", ""),
	}
	im.PackageName, _ = lo.Last(strings.Split(im.Path, "/"))
	im.Alias = im.PackageName
	if importSpec.Name != nil {
		im.Alias = importSpec.Name.Name
	}
	return im
}

type Imports []*ImportPackage

func NewImports(astImports []*ast.ImportSpec) Imports {
	return lo.Map(astImports, func(item *ast.ImportSpec, index int) *ImportPackage {
		return NewImportPackage(item)
	})
}

func (imps Imports) Append(importSpec *ast.ImportSpec) Imports {
	return append(imps, NewImportPackage(importSpec))
}

func (imps Imports) GetByAlias(alias string) (imp *ImportPackage) {
	imp, _ = lo.Find(imps, func(item *ImportPackage) bool {
		return item.Alias == alias
	})
	return
}
