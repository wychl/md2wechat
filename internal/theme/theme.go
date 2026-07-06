package theme

import (
	"embed"
	"fmt"
)

//go:embed themes/*.css
var themeFS embed.FS

// Theme 主题结构
type Theme struct {
	Name    string
	MainCSS string
}

// Load 加载主题 CSS（主题特定 CSS）
func Load(name string) (string, error) {
	themeCSS, err := themeFS.ReadFile("themes/" + name + ".css")
	if err != nil {
		// 如果主题不存在，回退到 default
		if name != "default" {
			themeCSS, err = themeFS.ReadFile("themes/default.css")
			if err != nil {
				return "", fmt.Errorf("默认主题加载失败: %w", err)
			}
		} else {
			return "", fmt.Errorf("主题文件 %s 不存在: %w", name, err)
		}
	}

	return string(themeCSS), nil
}

// GenerateStyle 生成最终的样式块
func GenerateStyle(themeName, primaryColor string) string {
	css, err := Load(themeName)
	if err != nil {
		// 若加载失败，返回最小样式
		css = "/* theme load failed */"
	}
	vars := fmt.Sprintf("<style>:root{--md-primary-color:%s}</style>", primaryColor)
	return vars + "<style>" + css + "</style>"
}
