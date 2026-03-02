package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeroicey/genmit/pkg/config"
	"github.com/zeroicey/genmit/pkg/git"
	"github.com/zeroicey/genmit/pkg/openai"
)

var generateCmd = &cobra.Command{
	Use:   "generate [dir]",
	Short: "Generate a commit message for the current or specified directory",
	Long: `Generate a commit message by analyzing git changes.
If no directory is specified, uses the current directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGenerate,
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Determine the directory to analyze
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Convert to absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("failed to resolve directory: %w", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		if _, ok := err.(*config.ConfigNotFoundError); ok {
			fmt.Println(err.Error())
			return nil
		}
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get git diff
	diff, err := git.GetDiff(absDir)
	if err != nil {
		return fmt.Errorf("failed to get git diff: %w", err)
	}

	// Check if diff needs to be truncated
	if cfg.MaxDiffSize > 0 && len(diff) > cfg.MaxDiffSize {
		diff, truncated := git.TruncateDiff(diff, cfg.MaxDiffSize)
		if truncated {
			fmt.Printf("⚠️  Diff too large (%d chars), truncated to %d chars\n", len(diff), cfg.MaxDiffSize)
		}
	}

	fmt.Println("Analyzing changes...")

	// Replace {lang} placeholder in prompt with actual language
	prompt := strings.ReplaceAll(cfg.Prompt, "{lang}", cfg.Lang)

	// Generate commit message
	client := openai.NewClient(cfg.BaseURL, cfg.APIKey)
	message, err := client.GenerateCommitMessage(prompt, diff, cfg.Model)
	if err != nil {
		return fmt.Errorf("failed to generate commit message: %w", err)
	}

	// Display the generated message
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Generated commit message:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println(message)
	fmt.Println(strings.Repeat("=", 50))

	// Ask if user wants to auto-commit
	fmt.Print("\nAuto-execute git commit? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response == "y" || response == "yes" {
		if err := git.Commit(absDir, message); err != nil {
			return fmt.Errorf("failed to execute git commit: %w", err)
		}
		fmt.Println("✓ Commit successful!")
	} else {
		fmt.Println("Commit cancelled. You can manually commit with:")
		fmt.Printf("  git commit -m \"%s\"\n", message)
	}

	return nil
}
