package service

import (
	"crypto/sha256"
	"time"

	"github.com/lwmacct/260630-go-hsr-auth/internal/repository"
)

type AuthOauthFlow struct {
	Provider         string
	PKCECodeVerifier string
	Nonce            string
	ReturnTo         string
	ExpiresAt        time.Time
	CreatedAt        time.Time
}

func utilAuthOauthFlow(row repository.AuthOauthFlowRow) AuthOauthFlow {
	return AuthOauthFlow{
		Provider:         row.Provider,
		PKCECodeVerifier: row.PKCECodeVerifier,
		Nonce:            row.Nonce,
		ReturnTo:         row.ReturnTo,
		ExpiresAt:        row.ExpiresAt,
		CreatedAt:        row.CreatedAt,
	}
}

func utilAuthOauthFlowStateHash(value string) []byte {
	hash := sha256.Sum256([]byte(value))
	return hash[:]
}
