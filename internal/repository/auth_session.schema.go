package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/lwmacct/260630-go-hsr-auth/internal/infra/dbschema"
)

type AuthSessionModel struct {
	bun.BaseModel `bun:"table:auth_sessions,alias:as"`

	TokenHash     []byte     `bun:"token_hash,pk"`
	UserID        int64      `bun:"user_id,notnull"`
	LoginIP       string     `bun:"login_ip,notnull"`
	LastIP        string     `bun:"last_ip,notnull"`
	UserAgentHash []byte     `bun:"user_agent_hash,notnull"`
	ExpiresAt     time.Time  `bun:"expires_at,notnull"`
	CreatedAt     time.Time  `bun:"created_at,notnull"`
	LastSeenAt    time.Time  `bun:"last_seen_at,notnull"`
	RevokedAt     *time.Time `bun:"revoked_at,nullzero"`
}

func (*AuthSessionModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(user_id) REFERENCES users (id) ON DELETE CASCADE")
	return nil
}

func AuthSessionSchema() dbschema.Schema {
	return dbschema.Schema{Models: []any{(*AuthSessionModel)(nil)}}
}
