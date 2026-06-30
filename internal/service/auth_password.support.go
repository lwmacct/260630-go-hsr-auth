package service

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	ErrAuthPasswordInvalidCredentials = errors.New("invalid credentials")
	ErrAuthPasswordWeakPassword       = errors.New("weak password")
)

type AuthPasswordRegisterInput struct {
	Username string
	Password string
}

func validateAuthPassword(username string, password string) error {
	length := utf8.RuneCountInString(password)
	if length < 8 || length > 128 || strings.TrimSpace(password) == "" {
		return ErrAuthPasswordWeakPassword
	}
	normalized := strings.ToLower(password)
	username = strings.ToLower(strings.TrimSpace(username))
	if username != "" && utf8.RuneCountInString(username) >= 3 && strings.Contains(normalized, username) {
		return ErrAuthPasswordWeakPassword
	}
	return nil
}
