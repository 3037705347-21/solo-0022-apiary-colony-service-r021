package domain

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var identifierPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// ValidateIdentifier 校验蜂场与蜂箱编号：小写字母、数字与单连字符。
func ValidateIdentifier(field, value string, minLen, maxLen int) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Invalid(field, "%s 不能为空", field)
	}
	if len(trimmed) < minLen || len(trimmed) > maxLen {
		return Invalid(field, "%s 长度必须在 %d 到 %d 之间", field, minLen, maxLen)
	}
	if !identifierPattern.MatchString(trimmed) {
		return Invalid(field, "%s 只能包含小写字母、数字和单连字符", field)
	}
	return nil
}

// ValidateText 校验文本字段的去除空白后的字符数。
func ValidateText(field, value string, minRunes, maxRunes int) error {
	count := utf8.RuneCountInString(strings.TrimSpace(value))
	if count < minRunes {
		if minRunes == 1 {
			return Invalid(field, "%s 不能为空", field)
		}
		return Invalid(field, "%s 至少需要 %d 个字符", field, minRunes)
	}
	if count > maxRunes {
		return Invalid(field, "%s 长度不能超过 %d 个字符", field, maxRunes)
	}
	return nil
}

func ValidateIntRange(field string, value, minValue, maxValue int) error {
	if value < minValue || value > maxValue {
		return Invalid(field, "%s 必须落在 %d 到 %d 之间", field, minValue, maxValue)
	}
	return nil
}

func TrimText(value string) string {
	return strings.TrimSpace(value)
}
