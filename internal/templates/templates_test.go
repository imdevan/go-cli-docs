package templates_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"go-cli-docs/internal/templates"
)

func TestEmbeddedTemplates(t *testing.T) {
	requiredFiles := []string{
		"config.mjs.tmpl",
		"sidebar.mjs.tmpl",
		"command.md.tmpl",
		"page.md.tmpl",
	}

	for _, name := range requiredFiles {
		t.Run(name, func(t *testing.T) {
			data, err := templates.FS.ReadFile(name)
			if err != nil {
				t.Fatalf("failed to read embedded file %s: %v", name, err)
			}
			if len(data) == 0 {
				t.Errorf("embedded file %s is empty", name)
			}
		})
	}
}

func TestRenderToFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "template_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tmpl, err := template.New("test").Parse("Hello {{ .Name }}!")
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	targetPath := filepath.Join(tempDir, "output.txt")
	data := struct {
		Name string
	}{
		Name: "World",
	}

	err = templates.RenderToFile(tmpl, targetPath, data)
	if err != nil {
		t.Fatalf("RenderToFile failed: %v", err)
	}

	content, err := ioutil.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("failed to read rendered file: %v", err)
	}

	expected := "Hello World!"
	if string(content) != expected {
		t.Errorf("expected %q, got %q", expected, string(content))
	}
}
