package repository

import "time"

type AuthOauthAccountCreate struct {
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

type AuthOauthAccountRow struct {
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

func utilAuthOauthAccountRow(row AuthOauthAccountModel) AuthOauthAccountRow {
	return AuthOauthAccountRow{
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
