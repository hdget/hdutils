package text

import (
	"regexp"
	"strings"
	"unicode"
)

type Utils interface {
	OnlyNumeric(s string) string
	OnlyAlphaNumeric(s string) string
	OnlyChinese(s string) string
	CleanString(origStr string, args ...bool) string
}

var (
	nonNumericRegex      = regexp.MustCompile(`[^0-9 ]+`)       // 非数字
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`) // 非英文字符和数字
	nonChineseRegex      = regexp.MustCompile(`[^\p{Han}]+`)    // 非汉字
)

// OnlyNumeric 去除数字以外的所有字符
func OnlyNumeric(s string) string {
	return nonNumericRegex.ReplaceAllString(s, "")
}

// OnlyAlphaNumeric 去除字符数字以外的所有字符
func OnlyAlphaNumeric(s string) string {
	return nonAlphanumericRegex.ReplaceAllString(s, "")
}

// OnlyChinese 去除中文以外的所有字符
func OnlyChinese(s string) string {
	return nonChineseRegex.ReplaceAllString(s, "")
}

// CleanString 处理字符串, args[0]为是否转换为小写
func CleanString(origStr string, args ...bool) string {
	// 1. 去除前后空格
	s := strings.TrimSpace(origStr)

	// 2. 是否转换小写
	toLower := false
	if len(args) > 0 {
		toLower = args[0]
	}

	if toLower {
		s = strings.ToLower(s)
	}

	// 去除不可见字符
	s = removeInvisibleCharacter(s)
	return s
}

// removeInvisibleCharacter 去除掉不能显示的字符
func removeInvisibleCharacter(origStr string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, origStr)
}
