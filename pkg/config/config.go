package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration
type Config struct {
	BaseURL     string `toml:"baseurl"`
	APIKey      string `toml:"apikey"`
	Prompt      string `toml:"prompt"`
	Model       string `toml:"model"`
	MaxDiffSize int    `toml:"maxdiffsize"`
	Lang        string `toml:"lang"`
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
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{
		BaseURL:     "https://api.openai.com/v1",
		Prompt:      getDefaultPrompt(),
		Model:       "gpt-4o-mini",
		MaxDiffSize: 10000,
		Lang:        "en",
	}

	if _, err := toml.Decode(string(data), config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
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

baseurl = "https://api.openai.com/v1"
apikey = "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
model = "gpt-4o-mini"
maxdiffsize = 10000
lang = "en"

prompt = """
You are an expert Git commit message generator. Analyze the following git diff and generate a clear, detailed commit message.

## Commit Message Format

### First Line: Title
- Use conventional commit format: type(scope): description
- Types: feat(new feature), fix(bug fix), docs(documentation), style(formatting), refactor(code restructuring), test(tests), chore(build/tooling)
- Keep title concise and descriptive (under 50 characters)

### Following Lines: Detailed Changes
- Use "- " prefix for bullet point list
- Each bullet point describes a specific change
- Cover all aspects: frontend, backend, config, docs, etc.
- Use professional and accurate technical terminology
- Organize logically (core functionality first, then auxiliary; backend before frontend)

### Language
- Generate the entire commit message in {lang} language

## Git Diff Content

{diff}

## Generate Commit Message
"""
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

	// Load existing config
	config := &Config{
		BaseURL:     "https://api.openai.com/v1",
		Prompt:      getDefaultPrompt(),
		Model:       "gpt-4o-mini",
		MaxDiffSize: 10000,
		Lang:        "en",
	}

	if data, err := os.ReadFile(configPath); err == nil {
		if _, err := toml.Decode(string(data), config); err != nil {
			return fmt.Errorf("failed to parse existing config: %w", err)
		}
	}

	// Update the specific key
	switch key {
	case "baseurl":
		config.BaseURL = value
	case "apikey":
		config.APIKey = value
	case "model":
		config.Model = value
	case "prompt":
		config.Prompt = value
	case "maxdiffsize":
		var size int
		if _, err := fmt.Sscanf(value, "%d", &size); err == nil {
			config.MaxDiffSize = size
		}
	case "lang":
		config.Lang = value
	default:
		return fmt.Errorf("invalid key '%s'. Valid keys are: baseurl, apikey, model, prompt, maxdiffsize, lang", key)
	}

	// Write back to file
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// getDefaultPrompt returns the default prompt template
func getDefaultPrompt() string {
	return `You are an expert Git commit message generator. Analyze the following git diff and generate a clear, detailed commit message.

## Commit Message Format

### First Line: Title
- Use conventional commit format: type(scope): description
- Types: feat(new feature), fix(bug fix), docs(documentation), style(formatting), refactor(code restructuring), test(tests), chore(build/tooling)
- Keep title concise and descriptive (under 50 characters)

### Following Lines: Detailed Changes
- Use "- " prefix for bullet point list
- Each bullet point describes a specific change
- Cover all aspects: frontend, backend, config, docs, etc.
- Use professional and accurate technical terminology
- Organize logically (core functionality first, then auxiliary; backend before frontend)

### Language
- Generate the entire commit message in {lang} language

## Git Diff Content

{diff}

## Generate Commit Message
`
}

// ConfigNotFoundError is returned when the config file doesn't exist
type ConfigNotFoundError struct {
	Msg string
}

func (e *ConfigNotFoundError) Error() string {
	return e.Msg
}
