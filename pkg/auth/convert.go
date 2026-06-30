package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/requestctx"
)

func ContextWithRequest(ctx context.Context, request SessionRequest) context.Context {
	return requestctx.ContextWithRequest(ctx, request)
}

func RequestFromContext(ctx context.Context) (SessionRequest, bool) {
	return requestctx.RequestFromContext(ctx)
}

func toServiceSessionRequest(value SessionRequest) service.AuthSessionInput {
	return service.AuthSessionInput{
		IP:         value.IP,
		Scheme:     value.Scheme,
		Host:       value.Host,
		UserAgent:  value.UserAgent,
		Method:     value.Method,
		Path:       value.Path,
		RemoteAddr: value.RemoteAddr,
	}
}

func sessionCookieValue(value string, expiresAt time.Time, config SessionCookieConfig) string {
	//nolint:gosec // Secure is an application cookie policy; local HTTP development intentionally uses insecure cookies.
	return (&http.Cookie{
		Name:     config.Name,
		Value:    value,
		Path:     config.Path,
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		HttpOnly: true,
		Secure:   config.Secure,
		SameSite: config.SameSite,
	}).String()
}
