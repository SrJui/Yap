package auth

import (
	"net/mail"
	"regexp"
)

// 3-18 chars
// _ allowed but not as first char
// . allowed but not as first char
var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9-][a-zA-Z0-9._-]{2,17}$`)

// min. 12 chars
// ascii allowed without space
var passwordPattern = regexp.MustCompile(`^[\x21-\x7E]{12,}$`)

func isValidPassword(password string) bool {
	return passwordPattern.MatchString(password)
}

func isValidEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return false
	}
	return true
}

func isValidUsername(username string) bool {
	return usernamePattern.MatchString(username)
}
