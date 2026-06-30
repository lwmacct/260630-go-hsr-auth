package auth

import (
	"context"
	"net/http"
	"time"

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
	ChallengeProviderImage     = "image"
	ChallengeProviderHCaptcha  = "hcaptcha"
	ChallengeProviderTurnstile = "turnstile"
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

type ChallengeProvider interface {
	Name() string
	PublicConfig() ChallengePublicConfig
	Create(context.Context, ChallengeInput) (*Challenge, error)
	Verify(context.Context, ChallengeAnswer, ChallengeInput) error
}

type ChallengeInput struct {
	IP         string
	UserAgent  string
	Method     string
	Path       string
	RemoteAddr string
}

type ChallengePublicConfig struct {
	Provider string
	SiteKey  string
}

type Challenge struct {
	Provider    string
	ChallengeID string
	Image       string
	ExpiresAt   time.Time
}

type ChallengeAnswer struct {
	Provider    string
	ChallengeID string
	Answer      string
	Token       string
}

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
