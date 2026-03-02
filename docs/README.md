# genmit

> AI-powered Git commit message generator CLI tool

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

genmit is a command-line tool that automatically generates Git commit messages using OpenAI API by analyzing your code changes.

## Features

- Automatically analyzes Git changes (staged and unstaged files)
- Uses OpenAI API to generate Conventional Commits compliant messages
- Multi-line commit message format with detailed bullet points
- Supports custom OpenAI-compatible APIs (DeepSeek, Alibaba Qwen, etc.)
- Customizable model, prompt templates, and output language
- Automatic diff truncation for large changes
- Interactive confirmation with optional auto-commit

## Installation

### Download pre-built binaries

Download from the [release](../release) page:

- **Windows**: `genmit-windows-amd64.exe`
- **Linux**: `genmit-linux-amd64` / `genmit-linux-arm64`
- **macOS**: `genmit-darwin-amd64` (Intel) / `genmit-darwin-arm64` (Apple Silicon)

### Build from source

```bash
git clone https://github.com/zeroicey/genmit.git
cd genmit
go build -o genmit .
```

### Install with Go

```bash
go install github.com/zeroicey/genmit@latest
```

## Configuration

On first run, genmit creates a config file template at `~/.genmit`, or you can set it manually:

```bash
# Set API Key
genmit config apikey sk-xxxxx

# Set API base URL (supports OpenAI-compatible APIs)
genmit config baseurl https://api.openai.com/v1

# Set model
genmit config model gpt-4o-mini

# Set commit message language (en, zh, ja, es, fr, etc.)
genmit config lang en

# Set max diff size (characters, 0 = no limit)
genmit config maxdiffsize 10000

# Set custom prompt template
genmit config prompt "Your custom prompt..."
```

### Config file example

`~/.genmit`:

```ini
baseurl=https://api.openai.com/v1
apikey=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
model=gpt-4o-mini
lang=en
maxdiffsize=10000
prompt=You are an expert Git commit message generator...
```

### List configuration

```bash
genmit config list
```

Output:

```
Config file: C:\Users\zeroicey\.genmit

baseurl     = https://api.openai.com/v1
apikey      = sk-9****xxxx
model       = gpt-4o-mini
lang        = en
maxdiffsize = 10000
prompt      = You are an expert Git commit message generator...
```

## Usage

### Basic usage

Generate commit message in current Git repository:

```bash
genmit
```

Specify directory:

```bash
genmit /path/to/repo
```

### Subcommands

```bash
# Generate commit message
genmit generate

# Manage configuration
genmit config list
genmit config <key> <value>
```

### Workflow

1. genmit analyzes Git changes in the current directory
2. Calls AI API to generate commit message
3. Displays the generated commit message
4. Asks whether to auto-execute `git commit`
   - Enter `y` or `yes` → auto commit
   - Enter anything else → only show message, commit manually

## Commit Message Format

genmit generates multi-line commit messages following Conventional Commits:

```
feat(dashboard): add moment submission frequency heatmap
- Add MomentHeatmap component to Dashboard page showing daily submission count
- Use react-calendar-heatmap library for heatmap rendering with color intensity
- Implement hover tooltip showing date and corresponding moment count
- Support responsive design for mobile, tablet, and desktop screen sizes
- Add /moments/statistics backend endpoint to aggregate submission data by date
- Implement related DTO, VO, Service, Repository layer logic
- Add useMomentStatistics Hook for requesting and managing statistics data
- Integrate dark mode color scheme using Tailwind CSS color system
```

The output language is controlled by the `lang` configuration option.

## Supported API Providers

genmit supports all OpenAI-compatible APIs:

- [OpenAI](https://platform.openai.com/) - `https://api.openai.com/v1`
- [DeepSeek](https://platform.deepseek.com/) - `https://api.deepseek.com/v1`
- [Alibaba Cloud Qwen](https://dashscope.aliyuncs.com/) - `https://dashscope.aliyuncs.com/compatible-mode/v1`
- Other OpenAI-compatible services

### Configuration examples

```bash
# DeepSeek
genmit config baseurl https://api.deepseek.com
genmit config model deepseek-chat

# Alibaba Cloud Qwen
genmit config baseurl https://dashscope.aliyuncs.com/compatible-mode/v1
genmit config model qwen-plus
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `apikey` | OpenAI API key | *required* |
| `baseurl` | API base URL | `https://api.openai.com/v1` |
| `model` | Model name | `gpt-4o-mini` |
| `lang` | Commit message language | `en` |
| `maxdiffsize` | Max diff characters (0 = unlimited) | `10000` |
| `prompt` | Custom prompt template | (built-in template) |

## Default prompt template

```
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
```

## Troubleshooting

### Git commit fails

Make sure you have configured Git user information:

```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

### API call fails

- Check if API Key is correct
- Check network connection
- Verify API base URL and model name match

### Diff too large warning

If you see `⚠️ Diff too large`, the diff has been truncated to fit within `maxdiffsize`. To:

- Increase the limit: `genmit config maxdiffsize 50000`
- Disable truncation: `genmit config maxdiffsize 0`

## License

[MIT License](LICENSE)

## Author

[@zeroicey](https://github.com/zeroicey)
