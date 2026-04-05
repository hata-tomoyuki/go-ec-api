package middleware

import "strings"

// TrimAndLimit は文字列の前後の空白を除去し、maxLen を超える場合は切り詰める。
// 各ハンドラーで入力を正規化する共通ヘルパーとして使用する。
//
// 例:
//
//	TrimAndLimit("  hello world  ", 5) → "hello"
//	TrimAndLimit("  short  ", 100)     → "short"
func TrimAndLimit(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	s = strings.TrimSpace(s)
	runes := []rune(s)
	if len(runes) > maxLen {
		s = string(runes[:maxLen])
	}
	return s
}
