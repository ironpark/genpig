package parser

import (
	"fmt"
	"github.com/samber/lo"
	"go/ast"
	"regexp"
	"strings"
)

type Struct struct {
	Name   string
	Fields []Field
	Ast    *ast.StructType
}

var tagParser *regexp.Regexp

func init() {
	tagParser = regexp.MustCompile("([^\":` \\t]+):\"([^\":` \\t]+)\"")
}
func NewStruct(name string, astStruct *ast.StructType, imports Imports) (st *Struct) {
	st = &Struct{
		Name:   name,
		Fields: nil,
		Ast:    astStruct,
	}
	for _, astField := range astStruct.Fields.List {
		field := NewField(astField, imports)
		if !field.IsPrimitiveType() {
			field.Package = imports.GetByAlias(strings.Split(field.Type, ".")[0])
		}
		st.Fields = append(st.Fields, field)
	}
	return st
}
func (st *Struct) UnWarpedFields(prefix string) (fields []Field) {
	for _, f := range st.Fields {
		if f.Nested != nil {
			fields = append(fields, f.Nested.UnWarpedFields(prefix+"."+f.Name)...)
		} else {
			f.Name = prefix + "." + f.Name
			fields = append(fields, f)
		}
	}
	return
}
func (st *Struct) OptionCheck(option string) (ok bool) {
	_, ok = lo.Find(st.Fields, func(item Field) bool {
		return item.Type == option
	})
	return
}

func (st *Struct) Dependencies() (imports []string) {
	return lo.Keys(lo.Reduce(st.Fields, func(agg map[string]bool, item Field, index int) map[string]bool {
		if item.Package != nil && item.Name != "" {
			agg[item.Package.String()] = true
		}
		return agg
	}, map[string]bool{}))
}

func (st *Struct) NotEmbeddingFields() []Field {
	return lo.Filter(st.Fields, func(item Field, index int) bool {
		return item.Name != ""
	})
}

func (st *Struct) TagExist(tag string) (ok bool) {
	_, ok = lo.Find(st.Fields, func(item Field) bool {
		return item.Has(tag)
	})
	return
}

func (st *Struct) String() (str string) {
	str = fmt.Sprintf("type %s struct {\n", st.Name)
	if st.Name == "" {
		str = "struct {\n"
	}
	for i := range st.Fields {
		str += "  " + st.Fields[i].String() + "\n"
	}
	str += "}"
	return
}

type Field struct {
	Name    string
	Type    string
	Tags    []Tag
	Comment string
	Package *ImportPackage
	Nested  *Struct
}

type Tag struct {
	Key   string
	Value string
}

func NewField(astField *ast.Field, imports Imports) (field Field) {
	switch v := astField.Type.(type) {
	case *ast.SelectorExpr:
		if ident, ok := v.X.(*ast.Ident); ok {
			field.Type = ident.Name + "."
		}
		field.Type += v.Sel.String()
	case *ast.Ident:
		field.Type = v.Name
	case *ast.StructType:
		//Nested Struct
		field.Nested = NewStruct("", v, imports)
		field.Type = field.Nested.String()
	}
	if astField.Tag != nil {
		for _, tag := range tagParser.FindAllStringSubmatch(astField.Tag.Value, -1) {
			field.Tags = append(field.Tags, Tag{
				Key:   tag[1],
				Value: tag[2],
			})
		}
	}
	if len(astField.Names) != 0 {
		field.Name = astField.Names[0].Name
	}
	return
}

func (f Field) String() string {
	tagString := "`" + strings.Join(lo.Map(f.Tags, func(item Tag, index int) string {
		return fmt.Sprintf("%s:\"%s\"", item.Key, item.Value)
	}), " ") + "`"
	return fmt.Sprintf("%s %s %s", f.Name, f.Type, tagString)
}
func (f Field) Has(tagKey string) bool {
	for _, tag := range f.Tags {
		if tag.Key == tagKey {
			return true
		}
	}
	return false
}
func (f Field) IsNestedStructType() bool {
	return f.Nested != nil
}
func (f Field) IsPrimitiveType() bool {
	if strings.Contains(f.Type, ".") {
		return false
	}
	// Pointer
	t := strings.ReplaceAll(f.Type, "*", "")
	switch t {
	case "string", "rune":
		return true
	case "float32", "float64":
		return true
	case "boolean":
		return true
	case "int8", "int16", "int32", "int64":
		return true
	case "uint8", "uint16", "uint32", "uint64":
		return true
	default:
		return false
	}
}
