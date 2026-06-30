package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lwmacct/260630-go-hsr-auth/internal/config"
)

func run(ctx context.Context, cfg *config.Config) error {
	if err := cfg.Server.HTTP.TLS.Validate(); err != nil {
		return err
	}

	deps, err := newDependencies(ctx, cfg)
	if err != nil {
		return err
	}
	defer deps.Close()

	srv := newHTTPServer(cfg, deps)
	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", srv.Addr)
	if err != nil {
		return err
	}
	defer func() { _ = ln.Close() }()

	errCh := make(chan error, 1)

	go func() {
		httpCfg := cfg.Server.HTTP
		slog.Info("web service starting", "listen", srv.Addr, "https", httpCfg.TLS.Enabled, "web_root", httpCfg.WebRoot)
		var serveErr error
		if httpCfg.TLS.Enabled {
			serveErr = srv.ServeTLS(ln, "", "")
		} else {
			serveErr = srv.Serve(ln)
		}
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errCh <- serveErr
		}
		close(errCh)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(sigCh)

	select {
	case <-ctx.Done():
		return shutdown(ctx, srv)
	case sig := <-sigCh:
		slog.Info("received shutdown signal", "signal", sig.String())
		return shutdown(ctx, srv)
	case err := <-errCh:
		return err
	}
}

func shutdown(ctx context.Context, srv *http.Server) error {
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	slog.Info("web service stopped")
	return nil
}
