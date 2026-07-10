package main

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var rootCmd = newRootCmd()

// @docs-command:root
//
//	name: go-cli-docs
//	description:
//		Generate Astro Starlight documentation for Go CLI projects.
//		The tool parses Cobra commands and flags, rendering markdown pages,
//		sidebar configs, and API docs.
//
//	example:
//		```bash
//		go-cli-docs init
//		go-cli-docs generate
//		go-cli-docs watch
//		```
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go-cli-docs",
		Short: "Generate Astro Starlight documentation for Go CLI projects",
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
