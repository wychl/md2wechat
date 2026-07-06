package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config 持有所有配置项
type Config struct {
	Theme        string
	PrimaryColor string
	OutputFile   string
	InputFile    string
	Help         bool
	SkillMode    bool // 技能模式
	JSONOutput   bool // 输出 JSON
	Quiet        bool // 静默模式
}

// Load 从命令行、环境变量、配置文件加载配置
func Load() (*Config, error) {
	// 1. 定义命令行标志（使用 pflag）
	fs := pflag.NewFlagSet("md2wechat", pflag.ContinueOnError)
	fs.String("theme", "default", "主题名称 (default, simple, grace, fresh, warm)")
	fs.String("color", "#0F4C81", "主色调，如 #0F4C81")
	fs.Bool("reading-time", false, "显示阅读时间预估")
	fs.String("output", "", "输出 HTML 文件 (默认 stdout)")
	fs.Bool("help", false, "显示帮助")
	fs.Bool("skill", false, "技能模式（-json -quiet）")
	fs.Bool("json", false, "输出 JSON 格式")
	fs.Bool("quiet", false, "静默模式，不输出提示信息")
	fs.String("input", "", "输入 Markdown 文件 (也可作为位置参数提供)")

	// 解析命令行参数
	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("解析命令行参数失败: %w", err)
	}

	// 2. 配置 viper
	v := viper.New()

	// 设置环境变量前缀，自动匹配大写、下划线
	v.SetEnvPrefix("MD2WECHAT")
	v.AutomaticEnv()

	// 绑定命令行标志到 viper
	if err := v.BindPFlags(fs); err != nil {
		return nil, fmt.Errorf("绑定命令行标志失败: %w", err)
	}

	// 3. 配置文件支持 (可选)
	v.SetConfigName(".md2wechat") // 配置文件名 (无扩展名)
	v.SetConfigType("yaml")       // 默认 yaml，也支持 json/toml
	v.AddConfigPath(".")          // 当前目录
	v.AddConfigPath("$HOME")      // 用户目录

	if err := v.ReadInConfig(); err != nil {
		// 配置文件不存在或不可读是正常情况，不报错
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 4. 提取配置值
	cfg := &Config{
		Theme:        v.GetString("theme"),
		PrimaryColor: v.GetString("color"),
		OutputFile:   v.GetString("output"),
		Help:         v.GetBool("help"),
		SkillMode:    v.GetBool("skill"),
		JSONOutput:   v.GetBool("json"),
		Quiet:        v.GetBool("quiet"),
		InputFile:    v.GetString("input"), // 可能来自 -input 或位置参数
	}

	// 处理 skill 模式
	if cfg.SkillMode {
		cfg.JSONOutput = true
		cfg.Quiet = true
	}

	// 如果已经提供了帮助，直接返回，不需要检查输入文件
	if cfg.Help {
		return cfg, nil
	}

	// 5. 确定输入文件（优先级：-input > 位置参数 > 标准输入）
	if cfg.InputFile == "" {
		// 获取未被标志消耗的位置参数
		args := fs.Args()
		if len(args) > 0 {
			cfg.InputFile = args[0]
		}
	}

	if cfg.InputFile == "" {
		// 检查是否从标准输入读取（管道）
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// 有管道数据，不要求输入文件
			cfg.InputFile = ""
		} else {
			return nil, fmt.Errorf("请指定输入文件（使用 -input 或位置参数）或通过管道传入内容")
		}
	}

	return cfg, nil
}

// PrintHelp 打印帮助信息
func PrintHelp() {
	fmt.Println(`md2wechat - Markdown 转微信公众号 HTML

用法:
  md2wechat [选项] [输入文件]

选项:
  -theme <名称>       主题: default, simple, grace, fresh, warm
  -color <颜色>       主色调，如 #0F4C81
  -reading-time       显示阅读时间预估
  -output <文件>      输出 HTML 文件
  -json               输出 JSON 格式
  -quiet              静默模式
  -skill              技能模式（等价于 -json -quiet）
  -input <文件>       输入 Markdown 文件（也可作为位置参数提供）
  -help               显示帮助

环境变量:
  MD2WECHAT_THEME, MD2WECHAT_COLOR, MD2WECHAT_READING_TIME, MD2WECHAT_OUTPUT

配置文件:
  支持 .md2wechat.yaml / .md2wechat.json / .md2wechat.toml，放在当前目录或用户目录。
  示例 (.md2wechat.yaml):
    theme: simple
    color: "#FF6600"
    reading-time: true
    output: "article.html"

优先级: 命令行 > 环境变量 > 配置文件 > 默认值

示例:
  md2wechat README.md -theme simple -output wechat.html
  cat doc.md | md2wechat -skill`)
}
