package convert

import (
	"bytes"
	"fmt"

	"github.com/wychl/md2wechat/internal/frontmatter"
	"github.com/wychl/md2wechat/internal/theme"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	goldhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Result 定义转换结果（用于 JSON 输出）
type Result struct {
	Success   bool   `json:"success"`
	HTML      string `json:"html,omitempty"`
	Title     string `json:"title,omitempty"`
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
			emoji.Emoji,
			extension.GFM,
			extension.Footnote,
			extension.TaskList,
			extension.DefinitionList,
			extension.Typographer,
			extension.CJK,
			extension.Table,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
			),
			NewCustomLiExtension(),
		),
		goldmark.WithRendererOptions(
			goldhtml.WithUnsafe(),
			renderer.WithNodeRenderers(
				util.Prioritized(&customLiRenderer{}, 1000),
			),
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

	// 3. 生成 CSS 样式
	cssStyle := theme.GenerateStyle(opts.Theme, opts.PrimaryColor)

	// 4. 组装完整 HTML
	fullHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=yes">
    <title>%s</title>
    %s
</head>
<body>
	<div>
	    %s
	</div>
    </div>
</body>
</html>`, title, cssStyle, htmlBody)

	return &Result{
		Success:   true,
		HTML:      fullHTML,
		Title:     title,
		ThemeUsed: opts.Theme,
	}, nil
}
