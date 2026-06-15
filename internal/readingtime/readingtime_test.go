package readingtime

import "testing"

func TestCountWords(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"Hello world", 2},
		{"你好世界", 4},
		{"Hello 世界", 2}, // 注意：正则匹配连续的中文和字母数字，这里“Hello”和“世界”是两个
		{"", 0}, 
		{"123 456", 2},
	}
	for _, tt := range tests {
		got := CountWords(tt.input)
		if got != tt.want {
			t.Errorf("CountWords(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestEstimate(t *testing.T) {
	minutes, words := Estimate("这是一段测试文字，用于估算阅读时间。" + "你好")
	if words == 0 || minutes == 0 {
		t.Errorf("Estimate returned zero: minutes=%d, words=%d", minutes, words)
	}
}
