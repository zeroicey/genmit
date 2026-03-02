# genmit

> AI-powered Git commit message generator CLI tool

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

genmit is a command-line tool that automatically generates Git commit messages using OpenAI API by analyzing your code changes.

## Features

- Automatically analyzes Git changes (staged and unstaged files)
- Uses OpenAI API to generate Conventional Commits compliant messages
- Supports custom OpenAI-compatible APIs (DeepSeek, Alibaba Qwen, etc.)
- Customizable model and prompt templates
- Interactive confirmation with optional auto-commit

## Installation

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

# Set custom prompt template
genmit config prompt "Your custom prompt..."
```

### Config file example

`~/.genmit`:

```ini
baseurl=https://api.openai.com/v1
apikey=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
model=gpt-4o-mini
prompt=You are a helpful assistant that generates Git commit messages...
```

### List configuration

```bash
genmit config list
```

Output:

```
Config file: C:\Users\zeroicey\.genmit

baseurl = https://api.openai.com/v1
apikey  = sk-9****xxxx
model   = gpt-4o-mini
prompt  = You are a helpful assistant that generates Git commit messag...
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

## Default prompt template

```
You are a helpful assistant that generates Git commit messages.
Given the following git diff output, generate a concise and descriptive commit message.
Follow conventional commit format: type(scope): description

Types: feat, fix, docs, style, refactor, test, chore

Git diff:
{diff}

Commit message:
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

## License

[MIT License](LICENSE)

## Author

[@zeroicey](https://github.com/zeroicey)
