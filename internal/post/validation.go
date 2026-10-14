package post

import (
	"strings"
	"unicode/utf8"
)

const maxContentLength = 1000

func isValidContent(content string) bool {
	if utf8.RuneCountInString(strings.TrimSpace(content)) == 0 {
		return false
	}
	if utf8.RuneCountInString(content) > maxContentLength {
		return false
	}
	return true
}
