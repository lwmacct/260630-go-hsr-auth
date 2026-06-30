package dbschema

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type Schema struct {
	Models []any
}

func Apply(ctx context.Context, db *bun.DB, schemas ...Schema) error {
	for _, schema := range schemas {
		for _, model := range schema.Models {
			if _, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
				return fmt.Errorf("create table: %w", err)
			}
		}
	}
	return nil
}

func Reset(ctx context.Context, db *bun.DB, tables []string, schemas ...Schema) error {
	for _, table := range tables {
		if _, err := db.NewDropTable().Table(table).IfExists().Cascade().Exec(ctx); err != nil {
			return fmt.Errorf("drop table %s: %w", table, err)
		}
	}
	return Apply(ctx, db, schemas...)
}
