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
	err = GenerateDocs(false)
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
	err = GenerateDocs(false)
	if err != nil {
		t.Fatalf("Second GenerateDocs call (idempotency check) failed: %v", err)
	}
}
