package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeroicey/genmit/pkg/config"
)

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "Manage configuration",
	Long: `Manage configuration in ~/.genmit.

Commands:
  genmit config list              Show all configuration values
  genmit config <key> <value>     Set a configuration value

Supported keys: apikey, baseurl, model, prompt`,
	RunE: runConfig,
}

func runConfig(cmd *cobra.Command, args []string) error {
	// No arguments - show help
	if len(args) == 0 {
		return cmd.Help()
	}

	// list command
	if args[0] == "list" {
		return runConfigList()
	}

	// set command (requires key and value)
	if len(args) != 2 {
		return fmt.Errorf("invalid arguments. Use 'genmit config list' or 'genmit config <key> <value>'")
	}

	return runConfigSet(args[0], args[1])
}

func runConfigList() error {
	cfg, err := config.Load()
	if err != nil {
		// If config file doesn't exist, show default values
		if _, ok := err.(*config.ConfigNotFoundError); ok {
			fmt.Println("Configuration file not found. Create one with 'genmit config <key> <value>'")
			return nil
		}
		return fmt.Errorf("failed to load config: %w", err)
	}

	configPath, _ := config.ConfigPath()
	fmt.Printf("Config file: %s\n\n", configPath)

	// Mask API key for security
	apiKey := cfg.APIKey
	if apiKey != "" && len(apiKey) > 8 {
		apiKey = apiKey[:4] + "****" + apiKey[len(apiKey)-4:]
	}

	fmt.Printf("baseurl = %s\n", cfg.BaseURL)
	fmt.Printf("apikey  = %s\n", apiKey)
	fmt.Printf("model   = %s\n", cfg.Model)
	fmt.Printf("prompt  = %s\n", truncateString(cfg.Prompt, 60))

	return nil
}

func runConfigSet(key, value string) error {
	key = strings.ToLower(key)

	// Validate key
	validKeys := map[string]bool{
		"apikey":  true,
		"baseurl": true,
		"model":   true,
		"prompt":  true,
	}

	if !validKeys[key] {
		return fmt.Errorf("invalid key '%s'. Valid keys are: apikey, baseurl, model, prompt", key)
	}

	// Set the configuration value
	if err := config.Set(key, value); err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}

	fmt.Printf("✓ Successfully set %s = %s\n", key, value)
	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
