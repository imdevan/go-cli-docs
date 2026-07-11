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
			if err := os.WriteFile(srcFile, []byte(tt.srcContent), 0644); err != nil {
				t.Fatalf("failed to write temp source file: %v", err)
			}
			defer os.Remove(srcFile)

			dstFile := "test_temp_dst.md"
			dstPath := filepath.Join("docs", "src", "content", "docs", dstFile)
			defer os.Remove(dstPath)

			// Run generation
			if err := generateContentPage(srcFile, dstFile, tt.defaultTitle, "Test desc", nil); err != nil {
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

