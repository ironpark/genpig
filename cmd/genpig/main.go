package main

import (
	"flag"
	"fmt"
	parser "github.com/ironpark/genpig/internal/parser"
	"github.com/ironpark/genpig/internal/templates"
	"github.com/samber/lo"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}

	currentGoFile := filepath.Join(dir, os.Getenv("GOFILE"))
	targetStruct := flag.String("struct", "", "struct name for generation")
	flag.Parse()
	currentGoFile = "/Users/ironpark/Documents/Project/Personal/genpig-example/conf/config.go"
	*targetStruct = "Config"
	if currentGoFile == "" {
		log.Println("GOFILE environment value is not set")
		return
	}
	if *targetStruct == "" {
		log.Println("-struct is required flag")
	}
	modPath, file, err := parser.GetModule(currentGoFile)
	if err != nil {
		return
	}

	moduleName := file.Module.Syntax.Token[1]
	relPath, _ := filepath.Rel(filepath.Dir(modPath), filepath.Dir(currentGoFile))
	basePackagePath := filepath.Join(moduleName, relPath)
	piggyDir := filepath.Dir(currentGoFile)
	fileName := filepath.Base(currentGoFile)
	piggyBaseFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "_gen.go"

	goFile := parser.ParseGoFile(currentGoFile)

	configStruct, ok := lo.Find(goFile.Structs, func(item *parser.Struct) bool {
		return item.Name == *targetStruct
	})
	if !ok { // Not Founded
		return
	}
	var configPaths string
	var configNames string
	for _, fc := range goFile.Init.FuncCalls {
		switch fc.Name {
		case "SetConfigPaths":
			configPaths = fc.ArgString()
		case "SetConfigNames":
			configNames = fc.ArgString()
		}
	}

	configMerge := lo.Map(configStruct.UnWarpedFields(), func(field *parser.Field, index int) string {
		values := lo.Map(field.Tags, func(tag parser.Tag, index int) string {
			if tag.Key == "env" {
				envFunction := "os.Getenv"
				switch field.Type {
				case "float64", "float32":
					envFunction = "genpig.EnvFloat"
				case "int", "int8", "int16", "int32", "int64":
					envFunction = "genpig.EnvInt"
				}
				return fmt.Sprintf("%s(%s(\"%s\"))", field.Type, envFunction, tag.Value)
			}
			return fmt.Sprintf("cfg%s%s", strings.ToTitle(tag.Key), field.Name)
		})

		merge := lo.Reduce(values, func(agg string, item string, index int) string {
			return agg + item + ",\n"
		}, "")
		return fmt.Sprintf("pig.cfg%s=genpig.Merge(\n%s)", field.Name, merge)
	})

	genPath := filepath.Join(piggyDir, piggyBaseFileName)
	err = TemplateGenerate(genPath, map[string]any{
		"Imports":               configStruct.Dependencies(),
		"PackageName":           goFile.PackageName,
		"OriginalStructPackage": filepath.Join(moduleName, relPath),
		"OriginalStructName":    *targetStruct,
		"Fields":                configStruct.Fields,
		"WithThreadSafe":        true,
		"WithSingleton":         true,
		"HasJsonTag":            configStruct.TagExist("json"),
		"HasYamlTag":            configStruct.TagExist("yaml"),
		"HasYamlToml":           configStruct.TagExist("toml"),
		"ConfigMerge":           configMerge,
		"ConfigPaths":           configPaths,
		"ConfigNames":           configNames,
	}, templates.PiggyBaseTmpl)
	if err == nil {
		exec.Command("go", "fmt", genPath).Run()
	}
	genPath = filepath.Join(piggyDir, "piggy", "piggy_gen.go")
	err = TemplateGenerate(genPath, map[string]any{
		"Fields":          configStruct.Fields,
		"BasePackage":     basePackagePath,
		"BasePackageName": relPath,
	}, templates.PiggyTmpl)
	//if err == nil {
	//	exec.Command("go", "fmt", genPath).Run()
	//}
}

func TemplateGenerate(path string, params map[string]any, tmpl *template.Template) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	err = tmpl.Execute(f, params)
	if err != nil {
		return err
	}
	return nil
}
