package frontmatter

import (
	"regexp"

	"gopkg.in/yaml.v3"
)

var frontMatterRegex = regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`)

// FrontMatter 文档元数据
type FrontMatter struct {
	Title string `yaml:"title"`
}

// Parse 从内容中提取 Front Matter 并返回元数据和正文
func Parse(content []byte) (FrontMatter, []byte) {
	matches := frontMatterRegex.FindSubmatch(content)
	if len(matches) != 3 {
		return FrontMatter{}, content
	}

	var fm FrontMatter
	_ = yaml.Unmarshal(matches[1], &fm)
	return fm, matches[2]
}
