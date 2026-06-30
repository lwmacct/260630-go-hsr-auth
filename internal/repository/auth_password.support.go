package repository

import "time"

type AuthPasswordCreate struct {
	UserID            string
	PasswordHash      string
	PasswordChangedAt time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type AuthPasswordRow struct {
	UserID            string
	PasswordHash      string
	PasswordChangedAt time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type AuthPasswordChange struct {
	Affected int64
}

func utilAuthPasswordRow(row AuthPasswordModel) AuthPasswordRow {
	return AuthPasswordRow{
		UserID:            row.UserID,
		PasswordHash:      row.PasswordHash,
		PasswordChangedAt: row.PasswordChangedAt,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}
