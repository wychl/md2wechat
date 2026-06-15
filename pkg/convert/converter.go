package convert

import (
	"bytes"
	"fmt"

	"github.com/wychl/md2wechat/internal/frontmatter"
	"github.com/wychl/md2wechat/internal/readingtime"
	"github.com/wychl/md2wechat/internal/theme"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// Result 定义转换结果（用于 JSON 输出）
type Result struct {
	Success   bool   `json:"success"`
	HTML      string `json:"html,omitempty"`
	Title     string `json:"title,omitempty"`
	Words     int    `json:"words,omitempty"`
	Minutes   int    `json:"minutes,omitempty"`
	ThemeUsed string `json:"theme_used,omitempty"`
	Error     string `json:"error,omitempty"`
}

// goldmarkConverter 内部转换器（复用原有实现）
type goldmarkConverter struct {
	md goldmark.Markdown
}

func newGoldmarkConverter() *goldmarkConverter {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.TaskList,
			extension.DefinitionList,
			extension.Typographer,
			extension.CJK,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
			),
		),
		goldmark.WithRendererOptions(
			goldhtml.WithUnsafe(),
		),
	)
	return &goldmarkConverter{md: md}
}

func (c *goldmarkConverter) convert(content []byte) (string, error) {
	var buf bytes.Buffer
	reader := text.NewReader(content)
	ctx := parser.NewContext()
	doc := c.md.Parser().Parse(reader, parser.WithContext(ctx))

	if err := c.md.Renderer().Render(&buf, content, doc); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Convert 是公开的核心转换函数
func Convert(input []byte, opts Options) (*Result, error) {
	// 1. 解析 Front Matter
	fm, body := frontmatter.Parse(input)
	title := fm.Title
	if title == "" {
		title = "Markdown 文档"
	}

	// 2. Markdown → HTML
	conv := newGoldmarkConverter()
	htmlBody, err := conv.convert(body)
	if err != nil {
		return &Result{Success: false, Error: err.Error()}, err
	}

	// 3. 阅读时间统计
	words, minutes := 0, 0
	if opts.ReadingTime {
		words, minutes = readingtime.Estimate(string(body))
	}

	// 4. 生成 CSS 样式
	cssStyle := theme.GenerateStyle(opts.Theme, opts.PrimaryColor)

	// 5. 阅读时间 HTML 片段
	readingHTML := ""
	if opts.ReadingTime {
		readingHTML = fmt.Sprintf(`<div class="reading-time">📖 阅读时间：约 %d 分钟 (%d 字)</div>`, minutes, words)
	}

	// 6. 组装完整 HTML
	fullHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=yes">
    <title>%s</title>
    %s
</head>
<body>
    <div id="output">
        <section class="container">
            %s
            %s
        </section>
    </div>
</body>
</html>`, title, cssStyle, readingHTML, htmlBody)

	return &Result{
		Success:   true,
		HTML:      fullHTML,
		Title:     title,
		Words:     words,
		Minutes:   minutes,
		ThemeUsed: opts.Theme,
	}, nil
}
