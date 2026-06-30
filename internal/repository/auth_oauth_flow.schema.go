package repository

import (
	"time"

	"github.com/uptrace/bun"

	"github.com/lwmacct/260630-go-hsr-auth/internal/infra/dbschema"
)

type AuthOauthFlowModel struct {
	bun.BaseModel `bun:"table:auth_oauth_flows,alias:aof"`

	StateHash        []byte    `bun:"state_hash,pk"`
	Provider         string    `bun:"provider,notnull"`
	PKCECodeVerifier string    `bun:"pkce_code_verifier,notnull"`
	Nonce            string    `bun:"nonce,nullzero"`
	ReturnTo         string    `bun:"return_to,nullzero"`
	ExpiresAt        time.Time `bun:"expires_at,notnull"`
	CreatedAt        time.Time `bun:"created_at,notnull"`
}

func AuthOauthFlowSchema() dbschema.Schema {
	return dbschema.Schema{Models: []any{(*AuthOauthFlowModel)(nil)}}
}
