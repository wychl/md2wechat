package convert

import (
	"log"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// 你的配置结构体
type Options struct {
	Theme        string
	PrimaryColor string
}

// 1. 定义一个代表你配置的唯一 OptionName (替代已废弃的 NewOptionCode)
const optConvertOptions renderer.OptionName = "ConvertOptions"

// 2. 编写可以将 Options 注入到 goldmark 渲染器的配置函数
func WithConvertOptions(opts Options) goldmark.Option {
	return goldmark.WithRendererOptions(
		renderer.WithOption(optConvertOptions, opts),
	)
}

// 修改结构体，让其持有 Options 字段
type customLiRenderer struct {
	opts Options
}

// 3. 实现 renderer.SetOptioner 接口，goldmark 初始化时会自动调用此方法分发配置
func (r *customLiRenderer) SetOption(name renderer.OptionName, value interface{}) {
	if name == optConvertOptions {
		if opts, ok := value.(Options); ok {
			r.opts = opts
		}
	}
}

func (r *customLiRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// 注册节点处理器
	reg.Register(ast.KindListItem, r.renderListItem)
}

func (r *customLiRenderer) renderListItem(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Printf("renderListItem: %v", n)
	log.Printf("parent: %v", n.Parent())

	parent := n.Parent()
	isUl := false
	if parent != nil {
		if list, ok := parent.(*ast.List); ok && !list.IsOrdered() {
			isUl = true
		}
	}

	if entering {
		// 动态使用从 SetOption 注入进来的配置
		if isUl && r.opts.PrimaryColor != "" {
			_, _ = w.WriteString(`<li style="color:` + r.opts.PrimaryColor + `;">`)
			_, _ = w.WriteString(`<span style="display:inline;">`)
		} else {
			_, _ = w.WriteString("<li>")
			if isUl {
				_, _ = w.WriteString(`<span style="display:inline;">`)
			}
		}
	} else {
		if isUl {
			_, _ = w.WriteString(`</span>`)
		}
		_, _ = w.WriteString("</li>")
	}
	return ast.WalkContinue, nil
}

// 4. 定义 Extender 并将渲染器实例作为指针传递
type CustomLiExtension struct{}

func NewCustomLiExtension() goldmark.Extender {
	return &CustomLiExtension{}
}

func (e *CustomLiExtension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			// 必须传指针 &customLiRenderer{}，否则 SetOption 方法无法修改内部的 opts 属性
			util.Prioritized(&customLiRenderer{}, 0),
		),
	)
}
