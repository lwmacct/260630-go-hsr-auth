package auth

import (
	"context"
	"fmt"

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
	for _, model := range Models() {
		if _, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
			return fmt.Errorf("create auth table: %w", err)
		}
	}
	return nil
}

func ResetSchema(ctx context.Context, db *bun.DB) error {
	for _, table := range []string{"auth_oauth_accounts", "auth_oauth_flows", "auth_passwords", "auth_sessions", "users"} {
		if _, err := db.NewDropTable().Table(table).IfExists().Cascade().Exec(ctx); err != nil {
			return fmt.Errorf("drop auth table %s: %w", table, err)
		}
	}
	return ApplySchema(ctx, db)
}
