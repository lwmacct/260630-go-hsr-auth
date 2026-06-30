package handler

import (
	"context"
	"net/http"

	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/challenge"
)

type Config struct {
	LocalLoginEnabled        bool
	LocalRegistrationEnabled bool
	SessionCookie            SessionCookieConfig
	RuntimeAdmins            []string
	Request                  RequestAuth
}

type Services struct {
	Users      *service.UserService
	Passwords  *service.AuthPasswordService
	Sessions   *service.AuthSessionService
	Challenges *challenge.Service
	AdminUsers *service.AdminUserService
}

type RequestAuth func(context.Context) (service.AuthSessionInput, bool)

type SessionCookieConfig struct {
	Name     string
	Path     string
	Secure   bool
	SameSite http.SameSite
}
