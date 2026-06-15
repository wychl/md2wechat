package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wychl/md2wechat/internal/config"
	"github.com/wychl/md2wechat/pkg/convert"
)

// Writer 输出接口
type Writer interface {
	Write(result *convert.Result) error
}

// FileWriter 输出到文件
type FileWriter struct {
	Filename string
}

func (w FileWriter) Write(result *convert.Result) error {
	return os.WriteFile(w.Filename, []byte(result.HTML), 0644)
}

// StdoutWriter 输出到标准输出
type StdoutWriter struct{}

func (w StdoutWriter) Write(result *convert.Result) error {
	_, err := fmt.Print(result.HTML)
	return err
}

// JsonWriter 输出到标准输出
type JsonWriter struct{}

func (w JsonWriter) Write(result *convert.Result) error {
	// JSON 模式：输出结构化结果到 stdout
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		return fmt.Errorf("JSON 编码失败: %w", err)
	}
	return nil
}

// NewWriter 根据配置创建输出器
func NewWriter(cfg *config.Config) Writer {
	if cfg == nil {
		return StdoutWriter{}
	}

	if cfg.JSONOutput {
		return JsonWriter{}
	}

	if cfg.OutputFile != "" {
		return FileWriter{Filename: cfg.OutputFile}
	}
	return StdoutWriter{}
}
