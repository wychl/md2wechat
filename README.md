# md2wechat

将 Markdown 转换为微信公众号 HTML 格式 —— Go 语言实现

[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## ✨ 特性

- ✅ **开箱即用** - 单一二进制文件，无需任何依赖
- ✅ **微信公众号兼容** - 生成的 HTML 完全符合微信平台规范
- ✅ **多种内置主题** - 5 种精美主题，支持自定义主色调
- ✅ **代码语法高亮** - 基于 goldmark-highlighting，支持 190+ 种语言
- ✅ **Front Matter 解析** - 支持 YAML 元数据（文章标题等）
- ✅ **阅读时间估算** - 自动统计字数并显示阅读时长
- ✅ **完全内联样式** - CSS 嵌入 HTML，无需外部文件
- ✅ **嵌入主题文件** - 使用 `embed` 将 CSS 打包进二进制，便于分发
- ✅ **灵活配置** - 支持命令行参数、环境变量、配置文件

## 📦 安装

### 从源码构建

```bash
git clone https://github.com/wychl/md2wechat
cd md2wechat
go build -o md2wechat ./cmd/md2wechat/main.go
```

### 使用 go install

```bash
go install github.com/wychl/md2wechat/cmd/md2wechat@latest
```

### 下载预编译二进制

访问 [Releases](https://github.com/wychl/md2wechat/releases) 页面下载对应平台的二进制文件。

## 🚀 使用

### 基础用法

```bash
# 转换文件并输出到标准输出
md2wechat README.md

# 输出到指定文件
md2wechat -output wechat.html article.md

# 使用特定主题和颜色
md2wechat -theme grace -color #92617E doc.md -output out.html

# 从标准输入读取（管道）
cat doc.md | md2wechat -theme simple
```

### 命令行选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `-theme` | 主题名称：`default`, `simple`, `grace`, `fresh`, `warm` | `default` |
| `-color` | 主色调（十六进制，如 `#0F4C81`） | `#0F4C81` |
| `-reading-time` | 显示阅读时间预估 | `false` |
| `-output` | 输出 HTML 文件路径（不指定则输出到 stdout） | - |
| `-json` | 输出 JSON 格式（包含 HTML 和元数据） | `false` |
| `-quiet` | 静默模式，不输出任何提示信息 | `false` |
| `-skill` | 技能模式（等价于 `-json -quiet`） | `false` |
| `-input` | 指定输入文件（也可作为位置参数） | - |
| `-help` | 显示帮助信息 | - |

### 环境变量

所有选项均可通过环境变量设置，前缀为 `MD2WECHAT_`：

| 环境变量 | 对应选项 |
|----------|----------|
| `MD2WECHAT_THEME` | `-theme` |
| `MD2WECHAT_COLOR` | `-color` |
| `MD2WECHAT_READING_TIME` | `-reading-time` |
| `MD2WECHAT_OUTPUT` | `-output` |

示例：

```bash
export MD2WECHAT_THEME=simple
export MD2WECHAT_COLOR="#2C3E50"
md2wechat doc.md
```

### 配置文件

支持配置文件 `.md2wechat.yaml` / `.md2wechat.json` / `.md2wechat.toml`，放在当前目录或用户目录。

示例 `.md2wechat.yaml`：

```yaml
theme: grace
color: "#92617E"
reading-time: true
output: "article.html"
quiet: false
```

**优先级**：命令行参数 > 环境变量 > 配置文件 > 默认值

### JSON 输出模式

使用 `-json` 或 `-skill` 时，程序输出 JSON 格式，便于其他工具集成：

```json
{
  "success": true,
  "html": "<div>...</div>",
  "meta": {
    "title": "文章标题",
    "reading_time": 3,
    "word_count": 456
  }
}
```

错误时输出：

```json
{
  "success": false,
  "error": "错误描述"
}
```

## 🎨 主题

| 主题 | 名称 | 默认主色调 | 风格 |
|------|------|------------|------|
| 经典蓝 | `default` | `#0F4C81` | 正式、专业 |
| 优雅紫 | `grace` | `#92617E` | 柔和、典雅 |
| 简洁灰 | `simple` | `#2C3E50` | 现代、极简 |
| 暖橙 | `warm` | `#E67E22` | 活力、温暖 |
| 清新绿 | `fresh` | `#27AE60` | 自然、护眼 |

## 📝 Front Matter 支持

在 Markdown 文件开头使用 YAML 格式定义元数据：

```markdown
---
title: 我的第一篇文章
author: 作者名
date: 2025-01-01
---

# 正文开始...
```

当前支持的字段：

- `title` - 文章标题（会显示在生成的 HTML 中）

## 📄 输出 HTML 规范

生成的 HTML 完全符合微信公众号平台要求：

- ✅ 所有样式内联（`<style>` 标签嵌入，微信支持）
- ✅ 使用微信兼容的 HTML 标签（`div`, `section`, `pre`, `code` 等）
- ✅ 图片自适应：`max-width: 100%` + `height: auto`
- ✅ 代码块横向滚动：`overflow-x: auto`
- ✅ 响应式布局：移动端字体和间距自动调整
- ✅ 无外部依赖：所有资源（CSS、字体）均内嵌或使用系统默认

## 🛠 开发

### 项目结构

```
md2wechat/
├── cmd/
│   └── md2wechat/              # 主入口
│       └── main.go
├── internal/                   # 私有包
│   ├── config/                 # 配置加载（支持 viper）
│   ├── frontmatter/            # YAML 元数据解析
│   ├── input/                  # 输入读取（文件/标准输入）
│   ├── output/                 # 输出处理（文件/stdout/JSON）
│   ├── readingtime/            # 阅读时间计算
│   └── theme/                  # 主题加载与嵌入
│       └── themes/             # CSS 主题文件（内嵌）
├── pkg/                        # 可公开导入的库
│   └── convert/                # 核心转换逻辑（导出 Convert 函数）
├── go.mod
├── go.sum
└── README.md
```

### 本地运行

```bash
# 克隆仓库
git clone https://github.com/wychl/md2wechat
cd md2wechat

# 直接运行示例
go run ./cmd/md2wechat/main.go -theme fresh README.md

# 构建
go build -o md2wechat ./cmd/md2wechat/main.go
```

### 添加新主题

1. 在 `internal/theme/themes/` 目录下创建 `your-theme.css`
2. 主题需配合 `base.css` 使用（仅覆盖变量或特定样式）
3. 重新构建：`go build ./cmd/md2wechat/main.go`

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License © 2026 wychl
