package readingtime

import (
	"regexp"
)

var wordRegex = regexp.MustCompile(`[\p{Han}\p{Latin}\p{N}]+`)

// CountWords 统计中英文及数字字符数量
func CountWords(text string) int {
	return len(wordRegex.FindAllString(text, -1))
}

// Estimate 根据字数估算阅读分钟数（假设 200 字/分钟）
func Estimate(text string) (minutes int, words int) {
	words = CountWords(text)
	minutes = (words + 199) / 200
	return
}
