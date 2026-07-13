package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateDocsIntegration(t *testing.T) {
	// Change working directory to the project root (up two levels from internal/workflow)
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(orig) //nolint:errcheck

	if err := os.Chdir("../.."); err != nil {
		t.Fatalf("failed to change directory to project root: %v", err)
	}

	// Run documentation generation with API docs disabled for testing speed
	err = GenerateDocs(false, nil)
	if err != nil {
		t.Fatalf("GenerateDocs failed: %v", err)
	}

	// Verify expected configuration and sidebar files are created
	configPath := filepath.Join("docs", "config.mjs")
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("Expected config.mjs to exist: %v", err)
	}

	sidebarPath := filepath.Join("docs", "sidebar.mjs")
	if _, err := os.Stat(sidebarPath); err != nil {
		t.Errorf("Expected sidebar.mjs to exist: %v", err)
	}

	astroConfigPath := filepath.Join("docs", "astro.config.mjs")
	if _, err := os.Stat(astroConfigPath); err != nil {
		t.Errorf("Expected astro.config.mjs to exist: %v", err)
	}

	// Verify index page is created
	indexPath := filepath.Join("docs", "src", "content", "docs", "index.md")
	if _, err := os.Stat(indexPath); err != nil {
		t.Errorf("Expected index.md to exist: %v", err)
	}

	// Verify command page is created
	generateCmdPath := filepath.Join("docs", "src", "content", "docs", "commands", "generate.md")
	if _, err := os.Stat(generateCmdPath); err != nil {
		t.Errorf("Expected generate.md to exist: %v", err)
	}

	// Verify landing page index.mdx is removed
	indexMdxPath := filepath.Join("docs", "src", "content", "docs", "index.mdx")
	if _, err := os.Stat(indexMdxPath); err == nil {
		t.Errorf("Expected index.mdx to be deleted/stripped, but it still exists")
	}

	// Verify config.mjs contents
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config.mjs: %v", err)
	}
	configContent := string(configBytes)
	if !strings.Contains(configContent, `title: "go-cli-docs"`) {
		t.Errorf("Expected config.mjs to contain correct title, got:\n%s", configContent)
	}

	// Verify astro.config.mjs contents
	astroConfigBytes, err := os.ReadFile(astroConfigPath)
	if err != nil {
		t.Fatalf("Failed to read astro.config.mjs: %v", err)
	}
	astroConfigContent := string(astroConfigBytes)
	if !strings.Contains(astroConfigContent, "import config from './config.mjs';") {
		t.Errorf("Expected astro.config.mjs to import config.mjs, got:\n%s", astroConfigContent)
	}

	// Test idempotency by running generation again
	err = GenerateDocs(false, nil)
	if err != nil {
		t.Fatalf("Second GenerateDocs call (idempotency check) failed: %v", err)
	}
}

func TestGenerateContentPage(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(orig) //nolint:errcheck

	if err := os.Chdir("../.."); err != nil {
		t.Fatalf("failed to change directory to project root: %v", err)
	}

	tests := []struct {
		name          string
		srcContent    string
		defaultTitle  string
		expectedTitle string
		expectedBody  string
	}{
		{
			name: "h1 header first line",
			srcContent: `# Hello World
This is a test.`,
			defaultTitle:  "Default",
			expectedTitle: "Hello World",
			expectedBody:  "\nThis is a test.",
		},
		{
			name: "h1 header with empty lines before it",
			srcContent: `

# Custom Title

Body goes here.`,
			defaultTitle:  "Default",
			expectedTitle: "Custom Title",
			expectedBody:  "\n\n\nBody goes here.",
		},
		{
			name: "h1 header with trailing hash",
			srcContent: `# Trailing Hash #
Body`,
			defaultTitle:  "Default",
			expectedTitle: "Trailing Hash",
			expectedBody:  "\nBody",
		},
		{
			name: "no h1 header (starts with h2)",
			srcContent: `## Subheading
Body`,
			defaultTitle:  "Default",
			expectedTitle: "Default",
			expectedBody:  "## Subheading\nBody",
		},
		{
			name: "no h1 header (starts with text)",
			srcContent: `Hello
## Subheading
Body`,
			defaultTitle:  "Default",
			expectedTitle: "Default",
			expectedBody:  "Hello\n## Subheading\nBody",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write temp source file
			srcFile := "test_temp_src.md"
			defer os.Remove(srcFile)
			if err := os.WriteFile(srcFile, []byte(tt.srcContent), 0644); err != nil {
				t.Fatalf("failed to write temp source file: %v", err)
			}

			dstFile := "test_temp_dst.md"
			dstPath := filepath.Join("docs", "src", "content", "docs", dstFile)
			defer os.Remove(dstPath)

			// Run generation
			if _, err := generateContentPage(srcFile, dstFile, tt.defaultTitle, "Test desc", nil); err != nil {
				t.Fatalf("generateContentPage failed: %v", err)
			}

			// Read generated file
			gotBytes, err := os.ReadFile(dstPath)
			if err != nil {
				t.Fatalf("failed to read generated file: %v", err)
			}
			gotContent := string(gotBytes)

			// Verify title in frontmatter
			expectedFrontmatterTitle := "title: " + tt.expectedTitle
			if !strings.Contains(gotContent, expectedFrontmatterTitle) {
				t.Errorf("expected frontmatter to contain %q, got:\n%s", expectedFrontmatterTitle, gotContent)
			}

			// Verify that the body content is as expected (usually after the frontmatter --- block)
			parts := strings.Split(gotContent, "---")
			if len(parts) < 3 {
				t.Fatalf("invalid generated markdown frontmatter format:\n%s", gotContent)
			}
			gotBody := strings.TrimSpace(strings.Join(parts[2:], "---"))
			expectedBodyTrimmed := strings.TrimSpace(tt.expectedBody)
			if gotBody != expectedBodyTrimmed {
				t.Errorf("expected body %q, got %q", expectedBodyTrimmed, gotBody)
			}
		})
	}
}

func TestParseCommandRelations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cmd_relations_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Write root.go
	rootContent := `package main
import "github.com/spf13/cobra"
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "app"}
	cmd.AddCommand(newConfigCmd())
	return cmd
}`
	if err := os.WriteFile(filepath.Join(tmpDir, "root.go"), []byte(rootContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Write config.go
	configContent := `package main
import "github.com/spf13/cobra"
func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "config"}
	cmd.AddCommand(newConfigInitCmd())
	return cmd
}`
	if err := os.WriteFile(filepath.Join(tmpDir, "config.go"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Write config_init.go
	configInitContent := `package main
import "github.com/spf13/cobra"
func newConfigInitCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "init"}
	return cmd
}`
	if err := os.WriteFile(filepath.Join(tmpDir, "config_init.go"), []byte(configInitContent), 0644); err != nil {
		t.Fatal(err)
	}

	relations, err := parseCommandRelations(tmpDir)
	if err != nil {
		t.Fatalf("parseCommandRelations failed: %v", err)
	}

	expected := map[string]string{
		"newConfigCmd":     "newRootCmd",
		"newConfigInitCmd": "newConfigCmd",
	}

	if len(relations) != len(expected) {
		t.Errorf("expected %d relations, got %d: %v", len(expected), len(relations), relations)
	}

	for k, v := range expected {
		if gotVal, ok := relations[k]; !ok || gotVal != v {
			t.Errorf("expected relations[%q] = %q, got %q", k, v, gotVal)
		}
	}
}

func TestBuildSidebarTree(t *testing.T) {
	commands := []CommandInfo{
		{CmdName: "completion", Path: []string{"completion"}},
		{CmdName: "config", Path: []string{"config"}},
		{CmdName: "config-init", Path: []string{"config", "init"}},
		{CmdName: "some-command", Path: []string{"some-command"}},
		{CmdName: "some-command-init", Path: []string{"some-command", "init"}},
		{CmdName: "go-cli-template", Path: []string{"go-cli-template"}}, // root command, should be skipped
	}

	tree := buildSidebarTree(commands, "go-cli-template")
	if len(tree) != 3 {
		t.Fatalf("expected 3 root items, got %d", len(tree))
	}

	// First item: completion
	if tree[0].Label != "completion" || tree[0].Link != "/commands/completion" || len(tree[0].Items) != 0 {
		t.Errorf("unexpected structure for completion node: %+v", tree[0])
	}

	// Second item: config (should be a group containing itself and config-init)
	configNode := tree[1]
	if configNode.Label != "config" || configNode.Link != "" || len(configNode.Items) != 2 {
		t.Fatalf("unexpected structure for config group: %+v", configNode)
	}

	if configNode.Items[0].Label != "config" || configNode.Items[0].Link != "/commands/config" {
		t.Errorf("unexpected first item of config: %+v", configNode.Items[0])
	}

	if configNode.Items[1].Label != "config init" || configNode.Items[1].Link != "/commands/config-init" {
		t.Errorf("unexpected second item of config: %+v", configNode.Items[1])
	}

	// Third item: some-command (hyphen preserved, group containing itself and init)
	someCmdNode := tree[2]
	if someCmdNode.Label != "some-command" || someCmdNode.Link != "" || len(someCmdNode.Items) != 2 {
		t.Fatalf("unexpected structure for some-command group: %+v", someCmdNode)
	}

	if someCmdNode.Items[0].Label != "some-command" || someCmdNode.Items[0].Link != "/commands/some-command" {
		t.Errorf("unexpected first item of some-command: %+v", someCmdNode.Items[0])
	}

	if someCmdNode.Items[1].Label != "some-command init" || someCmdNode.Items[1].Link != "/commands/some-command-init" {
		t.Errorf("unexpected second item of some-command: %+v", someCmdNode.Items[1])
	}

	// Verify formatting output
	formatted := formatSidebarItems(tree, "    ")
	expectedJS := `    { label: 'completion', link: '/commands/completion' },
    {
      label: 'config',
      items: [
        { label: 'config', link: '/commands/config' },
        { label: 'config init', link: '/commands/config-init' },
      ],
    },
    {
      label: 'some-command',
      items: [
        { label: 'some-command', link: '/commands/some-command' },
        { label: 'some-command init', link: '/commands/some-command-init' },
      ],
    },
`
	if formatted != expectedJS {
		t.Errorf("unexpected formatted JS:\n%s\nExpected:\n%s", formatted, expectedJS)
	}
}

