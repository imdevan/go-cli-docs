package main

import (
	"go-cli-docs/internal/workflow"
	"os"

	"github.com/spf13/cobra"
)

func newGenerateCmd() *cobra.Command {
	var genAPIDocs bool
	isProd := os.Getenv("NODE_ENV") == "production"
	defaultGenAPI := !isProd

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate all documentation from source",
		Long: `generate invokes the full docs generation pipeline:

  1. Read package metadata (TOML)
  2. Generate markdown content pages
  3. Parse Cobra commands
  4. Generate command documentation
  5. Generate API documentation (gomarkdoc)
  6. Generate config (config.mjs)
  7. Generate sidebar (sidebar.mjs)`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(genAPIDocs)
		},
	}

	cmd.Flags().BoolVar(&genAPIDocs, "gen-api-docs", defaultGenAPI, "Generate API documentation via gomarkdoc")

	return cmd
}

func runGenerate(genAPIDocs bool) error {
	return workflow.Generate(genAPIDocs)
}
