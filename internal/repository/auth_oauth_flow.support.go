package repository

import "time"

type AuthOauthFlowCreate struct {
	StateHash        []byte
	Provider         string
	PKCECodeVerifier string
	Nonce            string
	ReturnTo         string
	ExpiresAt        time.Time
	CreatedAt        time.Time
}

type AuthOauthFlowRow struct {
	StateHash        []byte
	Provider         string
	PKCECodeVerifier string
	Nonce            string
	ReturnTo         string
	ExpiresAt        time.Time
	CreatedAt        time.Time
}

func utilAuthOauthFlowRow(row AuthOauthFlowModel) AuthOauthFlowRow {
	return AuthOauthFlowRow{
		StateHash:        row.StateHash,
		Provider:         row.Provider,
		PKCECodeVerifier: row.PKCECodeVerifier,
		Nonce:            row.Nonce,
		ReturnTo:         row.ReturnTo,
		ExpiresAt:        row.ExpiresAt,
		CreatedAt:        row.CreatedAt,
	}
}
