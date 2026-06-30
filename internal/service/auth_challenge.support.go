package service

import (
	"context"
	"errors"
	"time"
)

const (
	AuthChallengeProviderImage     = "image"
	AuthChallengeProviderHCaptcha  = "hcaptcha"
	AuthChallengeProviderTurnstile = "turnstile"
)

type AuthChallengeInput struct {
	IP         string
	UserAgent  string
	Method     string
	Path       string
	RemoteAddr string
}

type AuthChallengePublicConfig struct {
	Provider string
	SiteKey  string
}

type AuthChallenge struct {
	Provider    string
	ChallengeID string
	Image       string
	ExpiresAt   time.Time
}

type AuthChallengeAnswer struct {
	Provider    string
	ChallengeID string
	Answer      string
	Token       string
}

var (
	ErrAuthChallengeInvalid       = errors.New("invalid challenge")
	ErrAuthChallengeUnsupported   = errors.New("challenge provider unsupported")
	ErrAuthChallengeLimitExceeded = errors.New("challenge limit exceeded")
)

type AuthChallengeProvider interface {
	Name() string
	PublicConfig() AuthChallengePublicConfig
	Create(context.Context, AuthChallengeInput) (*AuthChallenge, error)
	Verify(context.Context, AuthChallengeAnswer, AuthChallengeInput) error
}
