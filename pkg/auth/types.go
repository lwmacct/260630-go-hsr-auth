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

type Config struct {
	Local             LocalConfig
	Session           SessionConfig
	RuntimeAdmins     []string
	ChallengeProvider ChallengeProvider
	Request           RequestFunc
}

type LocalConfig struct {
	LoginEnabled        bool
	RegistrationEnabled bool
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
}

type RequestFunc func(context.Context) (SessionRequest, bool)

type SessionRequest = requestctx.Request

type ChallengeProvider = challenge.Provider

type ChallengeInput = challenge.Input

type ChallengePublicConfig = challenge.PublicConfig

type Challenge = challenge.Challenge

type ChallengeAnswer = challenge.Answer

type User struct {
	ID          int64
	Username    string
	DisplayName string
	Email       string
	AvatarURL   string
	Role        string
	Status      string
	Admin       bool
	DisabledAt  *time.Time
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Session struct {
	ID        string
	ExpiresAt time.Time
	SetCookie string
	User      *User
}

type ExternalUserInput struct {
	Username    string
	DisplayName string
	Email       string
	AvatarURL   string
}
