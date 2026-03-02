# genmit

> 使用 AI 自动生成 Git 提交信息的 CLI 工具

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

genmit 是一个命令行工具，通过分析你的 Git 更改，使用 OpenAI API 自动生成规范的提交信息。

## 功能特性

- 自动分析 Git 更改（已暂存和未暂存的文件）
- 使用 OpenAI API 生成符合 Conventional Commits 规范的提交信息
- 支持自定义 OpenAI 兼容的 API（如 DeepSeek、阿里通义千问等）
- 可自定义模型和提示词模板
- 交互式确认，可选择自动执行 git commit

## 安装

### 从源码编译

```bash
git clone https://github.com/zeroicey/genmit.git
cd genmit
go build -o genmit.exe .
```

### 使用 Go 安装

```bash
go install github.com/zeroicey/genmit@latest
```

## 配置

首次运行时，genmit 会在 `~/.genmit` 创建配置文件模板，你也可以手动设置：

```bash
# 设置 API Key
genmit config apikey sk-xxxxx

# 设置 API 地址（支持 OpenAI 兼容接口）
genmit config baseurl https://api.openai.com/v1

# 设置模型
genmit config model gpt-4o-mini

# 设置提示词模板
genmit config prompt "Your custom prompt..."
```

### 配置文件示例

`~/.genmit`:

```ini
baseurl=https://api.openai.com/v1
apikey=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
model=gpt-4o-mini
prompt=You are a helpful assistant that generates Git commit messages...
```

### 查看配置

```bash
genmit config list
```

输出：

```
Config file: C:\Users\zeroicey\.genmit

baseurl = https://api.openai.com/v1
apikey  = sk-9****xxxx
model   = gpt-4o-mini
prompt  = You are a helpful assistant that generates Git commit messag...
```

## 使用

### 基本用法

在当前 Git 仓库中生成提交信息：

```bash
genmit
```

指定目录：

```bash
genmit /path/to/repo
```

### 使用子命令

```bash
# 生成提交信息
genmit generate

# 管理配置
genmit config list
genmit config <key> <value>
```

### 工作流程

1. genmit 分析当前目录的 Git 更改
2. 调用 AI API 生成提交信息
3. 显示生成的提交信息
4. 询问是否自动执行 `git commit`
   - 输入 `y` 或 `yes` → 自动提交
   - 输入其他 → 仅显示信息，手动提交

## 支持的 API 提供商

genmit 支持所有 OpenAI 兼容的 API：

- [OpenAI](https://platform.openai.com/) - `https://api.openai.com/v1`
- [DeepSeek](https://platform.deepseek.com/) - `https://api.deepseek.com/v1`
- [阿里云通义千问](https://dashscope.aliyuncs.com/) - `https://dashscope.aliyuncs.com/compatible-mode/v1`
- 其他 OpenAI 兼容服务

### 配置示例

```bash
# DeepSeek
genmit config baseurl https://api.deepseek.com
genmit config model deepseek-chat

# 阿里云通义千问
genmit config baseurl https://dashscope.aliyuncs.com/compatible-mode/v1
genmit config model qwen-plus
```

## 默认提示词模板

```
You are a helpful assistant that generates Git commit messages.
Given the following git diff output, generate a concise and descriptive commit message.
Follow conventional commit format: type(scope): description

Types: feat, fix, docs, style, refactor, test, chore

Git diff:
{diff}

Commit message:
```

## 常见问题

### Git commit 失败

如果提交失败，请确保已配置 Git 用户信息：

```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

### API 调用失败

- 检查 API Key 是否正确
- 检查网络连接
- 确认 API 地址和模型名称匹配

## 许可证

[MIT License](LICENSE)

## 作者

[@zeroicey](https://github.com/zeroicey)
