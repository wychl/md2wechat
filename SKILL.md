---
name: md2wechat
description: 将 Markdown 文档转换为微信公众号兼容的 HTML 格式，支持主题、自定义颜色、阅读时间估算，输出 JSON 结构。
version: 1.0.0
author: wychl
---

# md2wechat Skill

## 📖 描述

`md2wechat` 读取 Markdown 文本（从文件或标准输入），应用主题样式，生成包含内联 CSS 且符合微信公众号规范的 HTML，并通过 JSON 返回结果。该技能为文章排版、知识库发布或集成到其他工具而设计。

## 🎯 触发场景

当用户提出以下需求时，应主动使用此技能：

- “将这段 Markdown 转成微信公众号 HTML”
- “生成微信图文格式”
- “帮我排版公众号文章”
- “转换这个 .md 文件为微信样式”
- “使用某个主题转换我的文档”

## ⚙️ 调用模式

| 属性       | 值                                      |
|------------|-------------------------------------------|
| 执行器     | `command`                                 |
| 命令       | `md2wechat -skill [选项]`                 |
| 输入       | 通过标准输入（stdin）传递 Markdown 内容     |
| 输出       | JSON 格式（详情见下文）                    |

> **重要**：所有调用**必须**包含 `-skill` 标志，该标志会启用 JSON 输出并关闭无关提示信息，确保程序输出纯 JSON。

## 📋 参数映射

用户输入中的参数将按以下方式映射到命令行选项（支持短名，但建议使用全名）：

| 用户参数       | 类型    | 默认值       | 命令行选项        | 说明                                 |
|----------------|---------|--------------|-------------------|--------------------------------------|
| `theme`        | string  | `"default"`  | `-theme`          | 主题名：`default`, `simple`, `grace`, `fresh`, `warm` |
| `color`        | string  | `"#0F4C81"`  | `-color`          | 主色调，十六进制，如 `#E67E22`       |
| `reading_time` | boolean | `false`      | `-reading-time`   | 是否在 HTML 中显示阅读时间估算       |

> 注意：`-skill` 标志由技能框架自动添加，用户无需指定。

## 🔧 命令模板

**标准模板（从 stdin 读取）：**

```bash
echo '{{.MarkdownContent}}' | md2wechat -skill \
  {{if .theme}}-theme "{{.theme}}"{{end}} \
  {{if .color}}-color "{{.color}}"{{end}} \
  {{if .reading_time}}-reading-time{{end}}
```

**如果用户提供了文件路径**（而不是文本内容），应先读取文件再传入：

```bash
cat "{{.FilePath}}" | md2wechat -skill \
  {{if .theme}}-theme "{{.theme}}"{{end}} \
  {{if .color}}-color "{{.color}}"{{end}} \
  {{if .reading_time}}-reading-time{{end}}
```

## 📤 输出解析

`md2wechat` 执行后会输出一行 JSON，其结构如下：

### 成功响应

```json
{
  "success": true,
  "html": "<!DOCTYPE html>...",
  "title": "从 Front Matter 提取的标题",
  "words": 1234,
  "minutes": 7,
  "theme_used": "default",
  "error": null
}
```

| 字段         | 说明                                           |
|--------------|------------------------------------------------|
| `success`    | 固定为 `true`                                 |
| `html`       | 完整的 HTML 字符串，可直接用于微信公众号        |
| `title`      | 文章标题（来自 Markdown 的 front matter 或正文第一个标题） |
| `words`      | 中文字数（不含代码块）                          |
| `minutes`    | 根据字数估算的阅读时间（分钟）                  |
| `theme_used` | 实际使用的主题名称                             |
| `error`      | 成功时为 `null`                               |

### 失败响应

```json
{
  "success": false,
  "html": "",
  "title": "",
  "words": 0,
  "minutes": 0,
  "theme_used": "",
  "error": "具体的错误原因，如文件不存在、语法错误等"
}
```

## 🤖 AI 处理指南（重要）

当命令执行成功（`success` 为 `true`）时，**不要**直接将 `html` 字段的完整内容返回给用户，而是按以下步骤处理：

1. **提取并清理标题**  
   - 使用 `title` 字段的值作为文件名基础。  
   - 移除非法字符：`/ \ : * ? " < > |`，将空格替换为下划线 `_`。  
   - 长度限制为 100 字符（超出则截断）。  
   - 若 `title` 为空或无效（如纯空白），则使用当前时间戳：`YYYYMMDD_HHMMSS`。

2. **确定输出目录**  
   - 优先使用用户提供的 `output_dir` 参数（如果存在于上下文中）。  
   - 否则，使用技能所在目录下的 `output` 子目录：  
     `<skill_base_dir>/output`  
     *技能基础目录通常为 `~/.openclaw/workspace/skills/md2wechat/`，若不确定，可执行 `pwd` 或通过脚本获取当前工作目录并推断。*  
   - 确保目录存在：使用 `mkdir -p <output_dir>`。

3. **写入文件**  
   - 完整路径 = `<output_dir>/<cleaned_title>.html`  
   - 将 `html` 字段的内容写入该文件（UTF‑8 编码，无 BOM）。

4. **返回给用户**  
   - 回复格式示例：  

     ```plaintext
     ✅ 转换成功！
     - 文件保存路径：/absolute/path/to/output/我的文章.html
     - 文章标题：我的文章
     - 字数：1234
     - 阅读时间：7 分钟
     
     可使用浏览器打开该文件预览，或直接复制内容到微信公众号编辑器。
     ```

   - 如果用户要求“不保存文件只显示预览”，可以仅返回 HTML 片段（但通常不推荐，因为 HTML 太长）。

若 `success` 为 `false`，则向用户清晰展示 `error` 字段内容，并建议检查 Markdown 格式、主题名称或颜色值。

## 📌 注意事项

- **输入大小限制**：Markdown 内容不应超过 10 MB（约 200 万字）。  
- **依赖**：需确保 `md2wechat` 二进制已安装在系统 `PATH` 中。若未找到，可提示用户安装（参考项目 README）。  
- **UTF‑8 编码**：输入和输出均使用 UTF‑8，避免乱码。  
- **管道兼容**：支持从标准输入读取，也支持使用 `-input` 指定文件，但为了统一，技能强制使用 stdin。

## 💡 示例

### 示例 1：基本转换

**用户输入：**  

```markdown
帮我将这段 Markdown 转成微信 HTML：
# 我的第一篇文章
欢迎阅读。
```

**AI 构造的命令：**  

```bash
echo '# 我的第一篇文章\n欢迎阅读。' | md2wechat -skill
```

**AI 处理流程：**  

- 解析 JSON → `title` = "我的第一篇文章"  
- 清理 → 文件名 = `我的第一篇文章.html`  
- 写入 `~/.openclaw/workspace/skills/md2wechat/output/我的第一篇文章.html`  
- 回复用户文件路径和统计信息。

### 示例 2：指定主题和阅读时间

**用户输入：**  

```markdown
使用优雅紫主题，并显示阅读时间，转换以下内容：
## 二级标题
正文内容...
```

**命令：**  

```bash
echo '## 二级标题\n正文内容...' | md2wechat -skill -theme grace
```

## 🔧 故障排查

| 问题                     | 可能原因                                   | 解决方法                                   |
|--------------------------|--------------------------------------------|--------------------------------------------|
| `command not found`      | `md2wechat` 不在 PATH 中                    | 提示用户安装：`go install github.com/wychl/md2wechat/cmd/md2wechat@latest` |
| 输出 JSON 解析失败       | 二进制版本过旧，不支持 `-skill`              | 升级到最新版（v1.0.0+）                    |
| 中文乱码                 | 终端或管道编码不是 UTF‑8                    | 确保系统环境为 UTF‑8；Windows 用户使用 `chcp 65001` |
| 生成的 HTML 在微信中样式错乱 | 主题或颜色值未正确传递                    | 检查命令行参数，确保颜色格式为 `#RRGGBB`   |
| 阅读时间不显示           | 忘记加 `-reading-time` 标志                 | 确认用户需要阅读时间时添加该标志           |
