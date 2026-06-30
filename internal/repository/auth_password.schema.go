package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/lwmacct/260630-go-hsr-auth/internal/infra/dbschema"
)

type AuthPasswordModel struct {
	bun.BaseModel `bun:"table:auth_passwords,alias:ap"`

	UserID            int64     `bun:"user_id,pk"`
	PasswordHash      string    `bun:"password_hash,notnull"`
	PasswordChangedAt time.Time `bun:"password_changed_at,notnull"`
	CreatedAt         time.Time `bun:"created_at,notnull"`
	UpdatedAt         time.Time `bun:"updated_at,notnull"`
}

func (*AuthPasswordModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(user_id) REFERENCES users (id) ON DELETE CASCADE")
	return nil
}

func AuthPasswordSchema() dbschema.Schema {
	return dbschema.Schema{Models: []any{(*AuthPasswordModel)(nil)}}
}
