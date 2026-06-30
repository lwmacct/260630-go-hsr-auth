package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/lwmacct/260630-go-hsr-auth/internal/config"
)

const httpAPIPrefix = "/api"

func newHTTPServer(cfg *config.Config, deps *dependencies) *http.Server {
	httpCfg := cfg.Server.HTTP
	srv := &http.Server{
		Addr:              httpCfg.Listen,
		Handler:           newHTTPHandler(cfg, deps),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       httpCfg.ReadTimeout,
		WriteTimeout:      httpCfg.WriteTimeout,
		IdleTimeout:       httpCfg.IdleTimeout,
	}

	if deps.tls == nil || deps.tls.TLSConfig() == nil {
		return srv
	}

	srv.TLSConfig = deps.tls.TLSConfig()
	return srv
}

func newHTTPHandler(cfg *config.Config, deps *dependencies) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(httpAPIPrefix+"/", http.StripPrefix(httpAPIPrefix, newHTTPAPIHandler(cfg, deps)))

	if cfg.Server.HTTP.WebRoot != "" {
		mux.Handle("/", http.FileServer(http.Dir(cfg.Server.HTTP.WebRoot)))
	}

	return deps.requests.Wrap(mux)
}

func newHTTPAPIHandler(cfg *config.Config, deps *dependencies) http.Handler {
	maxBodyBytes := cfg.Server.HTTP.MaxAPIBodyBytes
	if maxBodyBytes < 0 {
		maxBodyBytes = 0
	}
	return limitRequestBody(deps.auth.Handler(), maxBodyBytes)
}

func limitRequestBody(next http.Handler, maxBytes int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maxBytes > 0 && shouldLimitRequestBody(r) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		}
		next.ServeHTTP(w, r)
	})
}

func shouldLimitRequestBody(r *http.Request) bool {
	if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Body == nil || r.Body == http.NoBody {
		return false
	}
	for _, value := range r.Header.Values("Upgrade") {
		if strings.EqualFold(strings.TrimSpace(value), "websocket") {
			return false
		}
	}
	return true
}
