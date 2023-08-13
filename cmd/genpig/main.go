package main

import (
	"flag"
	"fmt"
	parser "github.com/ironpark/genpig/internal/parser"
	"github.com/ironpark/genpig/internal/templates"
	"github.com/samber/lo"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	currentGoFile := filepath.Join(dir, os.Getenv("GOFILE"))
	targetStruct := flag.String("struct", "", "struct name for generation")
	flag.Parse()
	modPath, file, err := parser.GetModule(currentGoFile)
	if err != nil {
		return
	}

	moduleName := file.Module.Syntax.Token[1]
	relPath, _ := filepath.Rel(filepath.Dir(modPath), filepath.Dir(currentGoFile))
	fmt.Println(relPath, moduleName)
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
	configMerge := lo.Map(configStruct.UnWarpedFields(""), func(field parser.Field, index int) string {
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
		return fmt.Sprintf("pig.cfg%s=genpig.Merge(%s)", field.Name, strings.Join(values, ","))
	})

	genPath := filepath.Join(piggyDir, piggyBaseFileName)
	err = PiggyBaseGenerate(genPath, map[string]any{
		"Imports":               configStruct.Dependencies(),
		"PackageName":           goFile.PackageName,
		"OriginalStructPackage": filepath.Join(moduleName, relPath),
		"OriginalStructName":    *targetStruct,
		"Fields":                configStruct.NotEmbeddingFields(),
		"WithThreadSafe":        true,
		"WithSingleton":         true,
		"WithSetter":            configStruct.OptionCheck("genpig.Setter"),
		"HasJsonTag":            configStruct.TagExist("json"),
		"HasYamlTag":            configStruct.TagExist("yaml"),
		"HasYamlToml":           configStruct.TagExist("toml"),
		"ConfigMerge":           configMerge,
		"ConfigPaths":           configPaths,
		"ConfigNames":           configNames,
	})
	if err == nil {
		exec.Command("go", "fmt", genPath).Run()
	}
	genPath = filepath.Join(piggyDir, "piggy", "piggy_gen.go")
	err = PiggyGenerate(genPath, map[string]any{
		"Fields": configStruct.NotEmbeddingFields(),
	})
	if err == nil {
		exec.Command("go", "fmt", genPath).Run()
	}
}

func PiggyBaseGenerate(path string, params map[string]any) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	err = templates.PiggyBase(f, params)
	if err != nil {
		return err
	}
	return nil
}

func PiggyGenerate(path string, params map[string]any) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	err = templates.Piggy(f, params)
	if err != nil {
		return err
	}
	return nil
}
