package templates

import (
	"embed"
	"io"
	"text/template"
)

//go:embed base.go.tmpl
var piggyBaseFile embed.FS
var piggyBaseTmpl *template.Template

//go:embed piggy.go.tmpl
var piggyFile embed.FS
var piggyTmpl *template.Template

func init() {
	piggyBaseTmpl = template.Must(template.ParseFS(piggyBaseFile, "*"))
	piggyTmpl = template.Must(template.ParseFS(piggyFile, "*"))
}

func PiggyBase(writer io.Writer, data any) error {
	return piggyBaseTmpl.Execute(writer, data)
}

func Piggy(writer io.Writer, data any) error {
	return piggyTmpl.Execute(writer, data)
}
