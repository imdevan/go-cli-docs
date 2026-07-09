package main

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var rootCmd = newRootCmd()

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go-cli-docs",
		Short: "Generate Astro Starlight documentation for Go CLI projects",
		Long: `go-cli-docs generates Astro Starlight documentation for Go CLI projects.

Commands:
  init      Scaffold the Astro Starlight docs directory
  generate  Generate all documentation from source
  watch     Watch source files and re-generate on change`,
		// No RunE – invoking the root command without a subcommand shows help.
		SilenceUsage: true,
	}

	var showVersion bool
	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print version and exit")
	cmd.Flags().String("config", "", "Path to config file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(resolvedVersion())
			return nil
		}
		return cmd.Help()
	}

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newGenerateCmd())
	cmd.AddCommand(newWatchCmd())
	cmd.AddCommand(newCompletionCmd())

	return cmd
}

// resolvedVersion returns the build version embedded by go build, or "dev".
func resolvedVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		v := info.Main.Version
		if v != "" && v != "(devel)" {
			return v
		}
	}
	return "dev"
}

func Execute() error {
	return rootCmd.Execute()
}
