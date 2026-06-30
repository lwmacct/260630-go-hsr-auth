package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lwmacct/260614-go-pkg-tlsreload/pkg/adapters/op"
	"github.com/lwmacct/260614-go-pkg-tlsreload/pkg/tlsreload"

	"github.com/lwmacct/260630-go-hsr-auth/internal/config"
	"github.com/lwmacct/260630-go-hsr-auth/internal/infra/database"
	"github.com/lwmacct/260630-go-hsr-auth/pkg/auth"

	"github.com/uptrace/bun"
)

type dependencies struct {
	db       *bun.DB
	auth     *auth.Module
	requests requestContextMiddleware
	tls      *tlsreload.Manager
}

func newDependencies(ctx context.Context, cfg *config.Config) (*dependencies, error) {
	deps, err := newDependenciesWithoutTLS(ctx, cfg)
	if err != nil {
		return nil, err
	}

	tlsManager, err := tlsreload.New(ctx, cfg.Server.HTTP.TLS, tlsreload.Options{
		Logger: slog.Default(),
		Adapters: []tlsreload.Adapter{
			op.New(op.Options{}),
		},
	})
	if err != nil {
		deps.Close()
		return nil, fmt.Errorf("configure tls: %w", err)
	}
	deps.tls = tlsManager
	return deps, nil
}

func newDependenciesWithoutTLS(ctx context.Context, cfg *config.Config) (*dependencies, error) {
	db, err := database.Open(ctx, cfg.Server.Database)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := auth.ApplySchema(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply database schema: %w", err)
	}

	deps := &dependencies{
		db:       db,
		requests: newRequestContextMiddleware(cfg.Server.HTTP.TrustedProxies),
	}
	module, err := auth.New(auth.Options{
		DB:         db,
		Config:     newAuthConfig(cfg),
		SessionTTL: cfg.Server.HTTP.SessionTTL,
	})
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("configure auth module: %w", err)
	}
	deps.auth = module
	return deps, nil
}

func (d *dependencies) Close() {
	if d == nil {
		return
	}
	if d.tls != nil {
		d.tls.Close()
		d.tls = nil
	}
	if d.db != nil {
		_ = d.db.Close()
		d.db = nil
	}
	d.auth = nil
}
