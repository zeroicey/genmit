# genmit

> 使用 AI 自动生成 Git 提交信息的 CLI 工具

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

[English Documentation](docs/README.md) | 中文文档

genmit 是一个命令行工具，通过分析你的 Git 更改，使用 OpenAI API 自动生成规范的提交信息。

## 功能特性

- 自动分析 Git 更改（已暂存和未暂存的文件）
- 使用 OpenAI API 生成符合 Conventional Commits 规范的多行提交信息
- 支持自定义 OpenAI 兼容的 API（如 DeepSeek、阿里通义千问等）
- 可自定义模型、提示词模板和输出语言
- 大型 Diff 自动截断功能
- 交互式确认，可选择自动执行 git commit

## 安装

### 下载预编译二进制文件

从 [release](release) 目录下载对应平台文件：

- **Windows**: `genmit-windows-amd64.exe`
- **Linux**: `genmit-linux-amd64` / `genmit-linux-arm64`
- **macOS**: `genmit-darwin-amd64` (Intel) / `genmit-darwin-arm64` (Apple Silicon)

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

# 设置提交信息语言 (en, zh, ja, es, fr 等)
genmit config lang zh

# 设置最大 diff 长度（字符数，0 = 不限制）
genmit config maxdiffsize 10000

# 设置提示词模板
genmit config prompt "Your custom prompt..."
```

### 配置文件示例

`~/.genmit` (TOML 格式):

```toml
baseurl = "https://api.openai.com/v1"
apikey = "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
model = "gpt-4o-mini"
lang = "zh"
maxdiffsize = 10000

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
```

### 查看配置

```bash
genmit config list
```

输出：

```
Config file: C:\Users\zeroicey\.genmit

baseurl     = https://api.openai.com/v1
apikey      = sk-9****xxxx
model       = gpt-4o-mini
lang        = zh
maxdiffsize = 10000
prompt      = You are an expert Git commit message generator...
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

## 提交信息格式

genmit 生成符合 Conventional Commits 规范的多行提交信息：

```
feat(dashboard): 添加Moment提交频率热力图功能
- 在Dashboard页面新增MomentHeatmap组件，展示过去一年每天的提交数量
- 使用react-calendar-heatmap库实现热力图渲染，颜色深浅表示活跃度
- 实现鼠标悬停显示日期和对应的Moment数量提示
- 支持响应式设计，适配手机、平板及桌面设备不同屏幕尺寸
- 后端新增/moments/statistics接口，按日期聚合统计用户提交数据
- 设计并实现相关DTO、VO、Service、Repository层逻辑
- 前端新增useMomentStatistics Hook用于请求和管理统计数据
- 集成深色模式配色方案，使用Tailwind CSS颜色体系
```

输出语言由 `lang` 配置项控制。

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

## 配置选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `apikey` | OpenAI API 密钥 | *必填* |
| `baseurl` | API 基础地址 | `https://api.openai.com/v1` |
| `model` | 模型名称 | `gpt-4o-mini` |
| `lang` | 提交信息语言 | `en` |
| `maxdiffsize` | 最大 diff 字符数（0 = 不限制） | `10000` |
| `prompt` | 自定义提示词模板 | (内置模板) |

## 默认提示词模板

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

### Diff 过大警告

如果看到 `⚠️ Diff too large`，说明 diff 已被截断以适应 `maxdiffsize` 限制。可以：

- 增加限制：`genmit config maxdiffsize 50000`
- 禁用截断：`genmit config maxdiffsize 0`

## 许可证

[MIT License](LICENSE)

## 作者

[@zeroicey](https://github.com/zeroicey)
