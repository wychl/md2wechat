package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wychl/md2wechat/internal/config"
	"github.com/wychl/md2wechat/internal/input"
	"github.com/wychl/md2wechat/internal/output"
	"github.com/wychl/md2wechat/pkg/convert"
)

func main() {
	run()
}

func run() {
	cfg, err := config.Load()
	if err != nil {
		outputJSONError(cfg, err)
		os.Exit(1)
	}
	if cfg.Help {
		config.PrintHelp()
		return
	}

	// 读取 markdown 内容
	var content []byte
	inputReader := input.NewReader(cfg)
	content, err = inputReader.Read()
	if err != nil {
		outputJSONError(cfg, fmt.Errorf("读取输入失败: %w", err))
		os.Exit(1)
	}

	// 转换 Markdown
	result, err := convert.Convert(content, convert.Options{
		Theme:        cfg.Theme,
		PrimaryColor: cfg.PrimaryColor,
		ReadingTime:  cfg.ReadingTime,
	})
	if err != nil {
		outputJSONError(cfg, fmt.Errorf("转换 Markdown 失败: %w", err))
		os.Exit(1)
	}

	// 非静默模式输出日志
	if !cfg.Quiet && !cfg.JSONOutput && cfg.OutputFile != "" {
		fmt.Fprintf(os.Stderr, "转换完成，输出文件: %s\n", cfg.OutputFile)
	}

	// 输出结果
	writer := output.NewWriter(cfg)
	if err := writer.Write(result); err != nil {
		outputJSONError(cfg, fmt.Errorf("写入输出失败: %w", err))
		os.Exit(1)
	}

	if !cfg.JSONOutput && cfg.OutputFile != "" && !cfg.Quiet {
		fmt.Printf("✅ 已写入: %s\n", cfg.OutputFile)
	}
}

func outputJSONError(cfg *config.Config, err error) {
	if cfg != nil && cfg.JSONOutput {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
	} else {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
	}
}
