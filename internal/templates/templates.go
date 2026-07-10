package templates

import (
	"embed"
	"os"
	"path/filepath"
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

// Load resolves a template by name, preferring a user-supplied override over
// the embedded default. Each entry in overrides is either a path to a
// specific template file or a directory containing a file named exactly
// name. Entries are checked in order; the first match wins. If no override
// matches, the embedded template is used.
func Load(name string, overrides []string) (*template.Template, error) {
	for _, o := range overrides {
		info, err := os.Stat(o)
		if err != nil {
			continue
		}

		var candidate string
		if info.IsDir() {
			candidate = filepath.Join(o, name)
			if _, err := os.Stat(candidate); err != nil {
				continue
			}
		} else {
			if filepath.Base(o) != name {
				continue
			}
			candidate = o
		}

		data, err := os.ReadFile(candidate)
		if err != nil {
			return nil, err
		}
		return template.New(name).Parse(string(data))
	}

	return template.ParseFS(FS, name)
}
