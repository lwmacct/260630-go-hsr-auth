package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/lwmacct/260630-go-hsr-shared/pkg/challenge"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/requestctx"
)

const (
	UserRoleUser  = "user"
	UserRoleAdmin = "admin"
)

const (
	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"
)

const (
	ChallengeProviderImage     = challenge.ProviderImage
	ChallengeProviderHCaptcha  = challenge.ProviderHCaptcha
	ChallengeProviderTurnstile = challenge.ProviderTurnstile
)

const (
	OAuthProviderGitHub = "github"
	OAuthProviderGoogle = "google"
)

type Config struct {
	Local             LocalConfig
	OAuth             OAuthConfig
	Session           SessionConfig
	RuntimeAdmins     []string
	OAuthProvider     func(string) (OAuthProvider, error)
	ChallengeProvider ChallengeProvider
	Request           RequestFunc
}

type LocalConfig struct {
	LoginEnabled        bool
	RegistrationEnabled bool
}

type OAuthConfig struct {
	Enabled         bool
	AutoRegister    bool
	CallbackBaseURL string
	Providers       []OAuthProviderConfig
}

type OAuthProviderConfig struct {
	Provider string
	Label    string
}

type SessionConfig struct {
	Cookie SessionCookieConfig
}

type SessionCookieConfig struct {
	Name     string
	Path     string
	Secure   bool
	SameSite http.SameSite
}

type Options struct {
	DB         DB
	Config     Config
	SessionTTL time.Duration
	OAuthTTL   time.Duration
}

type RequestFunc func(context.Context) (SessionRequest, bool)

type SessionRequest = requestctx.Request

type ChallengeProvider = challenge.Provider

type ChallengeInput = challenge.Input

type ChallengePublicConfig = challenge.PublicConfig

type Challenge = challenge.Challenge

type ChallengeAnswer = challenge.Answer

type OAuthProvider interface {
	Name() string
	AuthorizationURL(state string, redirectURI string, codeVerifier string, nonce string) string
	ExchangeProfile(ctx context.Context, code string, redirectURI string, codeVerifier string) (OAuthProfile, error)
}

type OAuthProfile struct {
	Provider              string
	Subject               string
	ProviderEmail         string
	ProviderEmailVerified bool
	ProviderDisplayName   string
	ProviderAvatarURL     string
	ProviderProfile       string
}
