package service

import (
	"errors"
	"time"

	"github.com/lwmacct/260630-go-hsr-auth/internal/repository"
)

type AuthOauthAccount struct {
	ID                    int64
	UserID                int64
	Provider              string
	Subject               string
	ProviderEmail         string
	ProviderEmailVerified bool
	ProviderDisplayName   string
	ProviderAvatarURL     string
	ProviderProfile       string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type AuthOauthAccountProfile struct {
	Provider              string
	Subject               string
	ProviderEmail         string
	ProviderEmailVerified bool
	ProviderDisplayName   string
	ProviderAvatarURL     string
	ProviderProfile       string
}

var (
	ErrAuthOauthAccountRegistrationDisabled = errors.New("oauth account registration disabled")
	ErrAuthOauthAccountTaken                = errors.New("oauth account taken")
)

type AuthOauthAccountResolveInput struct {
	Profile      AuthOauthAccountProfile
	AutoRegister bool
	Username     string
}

func utilAuthOauthAccount(row repository.AuthOauthAccountRow) AuthOauthAccount {
	return AuthOauthAccount{
		ID:                    row.ID,
		UserID:                row.UserID,
		Provider:              row.Provider,
		Subject:               row.Subject,
		ProviderEmail:         row.ProviderEmail,
		ProviderEmailVerified: row.ProviderEmailVerified,
		ProviderDisplayName:   row.ProviderDisplayName,
		ProviderAvatarURL:     row.ProviderAvatarURL,
		ProviderProfile:       row.ProviderProfile,
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
	}
}
