package workflow

import (
	"fmt"
	"os"
	"os/exec"
)

// InitAndGenerate ensures the docs folder exists, initializing Astro Starlight (version locked) if necessary, then runs generation.
func InitAndGenerate(genAPIDocs bool) error {
	docsDir := "docs"
	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		fmt.Println("Docs directory missing – initializing Astro Starlight...")
		// Initialize Astro Starlight using npm.
		initCmd := exec.Command("bash", "-c", fmt.Sprintf("npm create astro@latest %s -- --template starlight --yes --no-install", docsDir))
		initCmd.Stdout = os.Stdout
		initCmd.Stderr = os.Stderr
		if err := initCmd.Run(); err != nil {
			return fmt.Errorf("failed to init Astro Starlight: %w", err)
		}
		// Install dependencies.
		if err := exec.Command("bash", "-c", fmt.Sprintf("cd %s && npm install --silent", docsDir)).Run(); err != nil {
			return fmt.Errorf("failed to install Astro dependencies: %w", err)
		}
	}
	// Proceed with generation.
	return Generate(genAPIDocs)
}

// Generate runs the documentation generation workflow.
func Generate(genAPIDocs bool) error {
	fmt.Println("Running docs generation...")
	if err := GenerateDocs(genAPIDocs); err != nil {
		return fmt.Errorf("failed to generate docs: %w", err)
	}
	return nil
}
