package service

import (
	"crypto/sha256"
	"time"
)

const AuthSessionDefaultTTL = 7 * 24 * time.Hour

type AuthSessionUser struct {
	ID        string
	Username  string
	Admin     bool
	ExpiresAt time.Time
}

type AuthSessionInput struct {
	IP         string
	Scheme     string
	Host       string
	UserAgent  string
	Method     string
	Path       string
	RemoteAddr string
}

func utilAuthSessionTokenHash(value string) []byte {
	hash := sha256.Sum256([]byte(value))
	return hash[:]
}
