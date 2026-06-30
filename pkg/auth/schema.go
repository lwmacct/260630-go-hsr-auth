package auth

import (
	"context"

	"github.com/lwmacct/260630-go-hsr-shared/pkg/schema"
	"github.com/uptrace/bun"

	"github.com/lwmacct/260630-go-hsr-auth/internal/repository"
)

func Models() []any {
	return []any{
		(*repository.UserModel)(nil),
		(*repository.AuthPasswordModel)(nil),
		(*repository.AuthOauthAccountModel)(nil),
		(*repository.AuthSessionModel)(nil),
		(*repository.AuthOauthFlowModel)(nil),
	}
}

func ApplySchema(ctx context.Context, db *bun.DB) error {
	return schema.Apply(ctx, db, Models()...)
}

func ResetSchema(ctx context.Context, db *bun.DB) error {
	return schema.Reset(ctx, db, []string{"auth_oauth_accounts", "auth_oauth_flows", "auth_passwords", "auth_sessions", "users"}, Models()...)
}
