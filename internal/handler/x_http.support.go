package handler

import (
	"context"

	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
)

type Config struct {
	LocalLoginEnabled        bool
	LocalRegistrationEnabled bool
	OAuthEnabled             bool
	OAuthAutoRegister        bool
	OAuthCallbackBaseURL     string
	OAuthProviders           []OAuthProviderConfig
	TLSEnabled               bool
	SecureCookies            bool
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
	Challenges    *service.AuthChallengeService
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
