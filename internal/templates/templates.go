package templates

import (
	"embed"
	"text/template"
)

//go:embed base.go.tmpl
var piggyBaseFile embed.FS
var PiggyBaseTmpl *template.Template

//go:embed piggy.go.tmpl
var piggyFile embed.FS
var PiggyTmpl *template.Template

func init() {
	PiggyBaseTmpl = template.Must(template.ParseFS(piggyBaseFile, "*"))
	PiggyTmpl = template.Must(template.ParseFS(piggyFile, "*"))
}
