package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type UserModel struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID          int64      `bun:"id,pk,autoincrement"`
	Username    string     `bun:"username,notnull,unique"`
	DisplayName string     `bun:"display_name,notnull"`
	Email       string     `bun:"email,nullzero,unique"`
	AvatarURL   string     `bun:"avatar_url,nullzero"`
	Role        string     `bun:"role,notnull"`
	Status      string     `bun:"status,notnull"`
	DisabledAt  *time.Time `bun:"disabled_at,nullzero"`
	LastLoginAt *time.Time `bun:"last_login_at,nullzero"`
	CreatedAt   time.Time  `bun:"created_at,notnull"`
	UpdatedAt   time.Time  `bun:"updated_at,notnull"`
}
