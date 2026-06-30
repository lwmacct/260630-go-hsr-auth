package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
)

func utilRegisterErrorMessage(err error) string {
	switch {
	case errors.Is(err, service.ErrAuthPasswordWeakPassword):
		return "weak password"
	case errors.Is(err, service.ErrUserUsernameTaken):
		return "username taken"
	case errors.Is(err, service.ErrUserInvalidCredentials):
		return "invalid credentials"
	default:
		return "register failed"
	}
}

func utilRequest(ctx context.Context, config Config) (service.AuthSessionInput, error) {
	if config.Request == nil {
		return service.AuthSessionInput{}, huma.Error400BadRequest("invalid request source")
	}
	request, ok := config.Request(ctx)
	if !ok {
		return service.AuthSessionInput{}, huma.Error400BadRequest("invalid request source")
	}
	return request, nil
}

func utilSessionCookieValue(value string, expiresAt time.Time, config SessionCookieConfig) string {
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

func utilClearSessionCookie(config SessionCookieConfig) string {
	//nolint:gosec // Secure is an application cookie policy; local HTTP development intentionally uses insecure cookies.
	return (&http.Cookie{
		Name:     config.Name,
		Value:    "",
		Path:     config.Path,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   config.Secure,
		SameSite: config.SameSite,
	}).String()
}
