package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type AuthPasswordModel struct {
	bun.BaseModel `bun:"table:auth_passwords,alias:ap"`

	UserID            string    `bun:"user_id,pk,type:uuid"`
	PasswordHash      string    `bun:"password_hash,notnull"`
	PasswordChangedAt time.Time `bun:"password_changed_at,notnull"`
	CreatedAt         time.Time `bun:"created_at,notnull"`
	UpdatedAt         time.Time `bun:"updated_at,notnull"`
}

func (*AuthPasswordModel) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(user_id) REFERENCES users (id) ON DELETE CASCADE")
	return nil
}
