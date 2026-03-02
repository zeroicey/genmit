package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "genmit",
	Short: "Genmit generates Git commit messages using AI",
	Long: `Genmit is a CLI tool that automatically generates Git commit messages
using OpenAI's API. It analyzes your git changes and creates descriptive
commit messages following conventional commit format.`,
	// Run generate by default when no subcommand is specified
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGenerate(cmd, args)
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(configCmd)
}
