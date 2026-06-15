package input

import (
	"io"
	"os"

	"github.com/wychl/md2wechat/internal/config"
)

// Reader 输入接口
type Reader interface {
	Read() ([]byte, error)
}

type FileReader struct {
	Filename string
}

func (r FileReader) Read() ([]byte, error) {
	return os.ReadFile(r.Filename)
}

type StdinReader struct{}

func (r StdinReader) Read() ([]byte, error) {
	return io.ReadAll(os.Stdin)
}

func NewReader(cfg *config.Config) Reader {
	if cfg == nil {
		return StdinReader{}
	}

	if cfg.InputFile != "" {
		return FileReader{Filename: cfg.InputFile}
	}
	return StdinReader{}
}
