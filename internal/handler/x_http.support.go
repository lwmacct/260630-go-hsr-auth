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
	OAuthEnabled             bool
	OAuthAutoRegister        bool
	OAuthCallbackBaseURL     string
	OAuthProviders           []OAuthProviderConfig
	SessionCookie            SessionCookieConfig
	RuntimeAdmins            []string
	Request                  RequestAuth
	OAuthProvider            func(string) (OAuthProviderAuth, error)
}

type Services struct {
	Users         *service.UserService
	Passwords     *service.AuthPasswordService
	OAuthAccounts *service.AuthOauthAccountService
	Sessions      *service.AuthSessionService
	OAuthFlows    *service.AuthOauthFlowService
	Challenges    *challenge.Service
	AdminUsers    *service.AdminUserService
}

type OAuthProviderAuth interface {
	Name() string
	AuthorizationURL(state string, redirectURI string, codeVerifier string, nonce string) string
	ExchangeProfile(ctx context.Context, code string, redirectURI string, codeVerifier string) (service.AuthOauthAccountProfile, error)
}

type RequestAuth func(context.Context) (service.AuthSessionInput, bool)

type OAuthProviderConfig struct {
	Provider string
	Label    string
}

type SessionCookieConfig struct {
	Name     string
	Path     string
	Secure   bool
	SameSite http.SameSite
}
