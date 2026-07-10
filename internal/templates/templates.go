package templates

import (
	"embed"
	"os"
	"text/template"
)

// FS holds the embedded template files.
//go:embed *.tmpl
var FS embed.FS

// RenderToFile executes the given template with data and writes the result to path.
func RenderToFile(tmpl *template.Template, path string, data any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, data)
}
