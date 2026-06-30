package repository

import "time"

type AuthSessionCreate struct {
	TokenHash     []byte
	UserID        int64
	LoginIP       string
	LastIP        string
	UserAgentHash []byte
	ExpiresAt     time.Time
	CreatedAt     time.Time
	LastSeenAt    time.Time
}

type AuthSessionRow struct {
	TokenHash     []byte
	UserID        int64
	LoginIP       string
	LastIP        string
	UserAgentHash []byte
	ExpiresAt     time.Time
	CreatedAt     time.Time
	LastSeenAt    time.Time
	RevokedAt     *time.Time
}

type AuthSessionChange struct {
	Affected int64
}

func utilAuthSessionRow(row AuthSessionModel) AuthSessionRow {
	return AuthSessionRow{
		TokenHash:     row.TokenHash,
		UserID:        row.UserID,
		LoginIP:       row.LoginIP,
		LastIP:        row.LastIP,
		UserAgentHash: row.UserAgentHash,
		ExpiresAt:     row.ExpiresAt,
		CreatedAt:     row.CreatedAt,
		LastSeenAt:    row.LastSeenAt,
		RevokedAt:     row.RevokedAt,
	}
}
