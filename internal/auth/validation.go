package auth

import (
	"unicode/utf8"
)

func ValidatePassword(password string) bool {
	if utf8.RuneCountInString(password) < 12 {
		return false
	}
	return true
}
