package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the application configuration
type Config struct {
	BaseURL string
	APIKey  string
	Prompt  string
	Model   string
}

// ConfigPath returns the path to the config file (~/.genmit)
func ConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".genmit"), nil
}

// Load reads the configuration from ~/.genmit
// If the file doesn't exist, it creates a template and prompts the user
func Load() (*Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create template config file
		if err := createTemplate(configPath); err != nil {
			return nil, fmt.Errorf("failed to create config template: %w", err)
		}
		return nil, &ConfigNotFoundError{
			Msg: fmt.Sprintf("Config file created at %s\nPlease edit it and fill in your apikey before running genmit again.", configPath),
		}
	}

	// Read and parse the config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := &Config{
		BaseURL: "https://api.openai.com/v1",
		Prompt:  getDefaultPrompt(),
		Model:   "gpt-4o-mini",
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])

		switch key {
		case "baseurl":
			config.BaseURL = value
		case "apikey":
			config.APIKey = value
		case "prompt":
			config.Prompt = value
		case "model":
			config.Model = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Validate required fields
	if config.APIKey == "" {
		return nil, fmt.Errorf("apikey is not set in %s\nPlease edit the file and add your OpenAI API key", configPath)
	}

	return config, nil
}

// createTemplate creates a template config file
func createTemplate(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	template := `# genmit configuration file
# Generate your API key at: https://platform.openai.com/api-keys

baseurl=https://api.openai.com/v1
apikey=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
model=gpt-4o-mini
prompt=You are a helpful assistant that generates Git commit messages. Given the following git diff output, generate a concise and descriptive commit message. Follow conventional commit format: type(scope): description. Types: feat, fix, docs, style, refactor, test, chore. Git diff:
{diff}

Commit message:
`

	_, err = file.WriteString(template)
	return err
}

// Set updates a single configuration value
func Set(key, value string) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	// Read existing config or create new one
	var lines []string
	if file, err := os.Open(configPath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	}

	// Find and update the key, or add it if it doesn't exist
	key = strings.ToLower(key)
	found := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 && strings.TrimSpace(strings.ToLower(parts[0])) == key {
			lines[i] = fmt.Sprintf("%s=%s", parts[0], value)
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	// Write back to file
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// getDefaultPrompt returns the default prompt template
func getDefaultPrompt() string {
	return `You are a helpful assistant that generates Git commit messages. Given the following git diff output, generate a concise and descriptive commit message. Follow conventional commit format: type(scope): description. Types: feat, fix, docs, style, refactor, test, chore. Git diff:
{diff}

Commit message:
`
}

// ConfigNotFoundError is returned when the config file doesn't exist
type ConfigNotFoundError struct {
	Msg string
}

func (e *ConfigNotFoundError) Error() string {
	return e.Msg
}
