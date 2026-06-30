package auth

import (
	"context"

	"github.com/lwmacct/260630-go-hsr-auth/internal/handler"
	"github.com/lwmacct/260630-go-hsr-auth/internal/infra/captcha"
	"github.com/lwmacct/260630-go-hsr-auth/internal/infra/oauthclient"
	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
)

type OAuthClientConfig struct {
	Enabled      bool
	ClientID     string
	ClientSecret string
	Scopes       []string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
}

func NewOAuthProvider(name string, config OAuthClientConfig) (OAuthProvider, error) {
	provider, err := oauthclient.New(name, oauthclient.ProviderConfig{
		Enabled:      config.Enabled,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       config.Scopes,
		AuthURL:      config.AuthURL,
		TokenURL:     config.TokenURL,
		UserInfoURL:  config.UserInfoURL,
	})
	if err != nil {
		return nil, err
	}
	return oauthClientProvider{provider: provider}, nil
}

func NewImageChallengeProvider(maxItems int) ChallengeProvider {
	return serviceChallengeProviderProxy{provider: captcha.NewImageProvider(maxItems)}
}

func NewRemoteTokenChallengeProvider(provider string, siteKey string, secret string, verifyURL string) (ChallengeProvider, error) {
	value, err := captcha.NewRemoteTokenProvider(provider, siteKey, secret, verifyURL)
	if err != nil {
		return nil, err
	}
	return serviceChallengeProviderProxy{provider: value}, nil
}

type oauthClientProvider struct {
	provider *oauthclient.Provider
}

func (p oauthClientProvider) Name() string {
	return p.provider.Name()
}

func (p oauthClientProvider) AuthorizationURL(state string, redirectURI string, codeVerifier string, nonce string) string {
	return p.provider.AuthorizationURL(state, redirectURI, codeVerifier, nonce)
}

func (p oauthClientProvider) ExchangeProfile(ctx context.Context, code string, redirectURI string, codeVerifier string) (OAuthProfile, error) {
	profile, err := p.provider.ExchangeProfile(ctx, oauthclient.TokenRequest{
		Code:         code,
		RedirectURI:  redirectURI,
		CodeVerifier: codeVerifier,
	})
	if err != nil {
		return OAuthProfile{}, err
	}
	return fromServiceOAuthProfile(profile), nil
}

type handlerOAuthProvider struct {
	provider OAuthProvider
}

func (p handlerOAuthProvider) Name() string {
	return p.provider.Name()
}

func (p handlerOAuthProvider) AuthorizationURL(state string, redirectURI string, codeVerifier string, nonce string) string {
	return p.provider.AuthorizationURL(state, redirectURI, codeVerifier, nonce)
}

func (p handlerOAuthProvider) ExchangeProfile(ctx context.Context, code string, redirectURI string, codeVerifier string) (service.AuthOauthAccountProfile, error) {
	profile, err := p.provider.ExchangeProfile(ctx, code, redirectURI, codeVerifier)
	if err != nil {
		return service.AuthOauthAccountProfile{}, err
	}
	return toServiceOAuthProfile(profile), nil
}

var _ handler.OAuthProviderAuth = handlerOAuthProvider{}

type serviceChallengeProvider struct {
	provider ChallengeProvider
}

func (p serviceChallengeProvider) Name() string {
	return p.provider.Name()
}

func (p serviceChallengeProvider) PublicConfig() service.AuthChallengePublicConfig {
	return toServiceChallengePublicConfig(p.provider.PublicConfig())
}

func (p serviceChallengeProvider) Create(ctx context.Context, input service.AuthChallengeInput) (*service.AuthChallenge, error) {
	challenge, err := p.provider.Create(ctx, fromServiceChallengeInput(input))
	if err != nil {
		return nil, err
	}
	return toServiceChallenge(challenge), nil
}

func (p serviceChallengeProvider) Verify(ctx context.Context, answer service.AuthChallengeAnswer, input service.AuthChallengeInput) error {
	return p.provider.Verify(ctx, fromServiceChallengeAnswer(answer), fromServiceChallengeInput(input))
}

type serviceChallengeProviderProxy struct {
	provider service.AuthChallengeProvider
}

func (p serviceChallengeProviderProxy) Name() string {
	return p.provider.Name()
}

func (p serviceChallengeProviderProxy) PublicConfig() ChallengePublicConfig {
	return fromServiceChallengePublicConfig(p.provider.PublicConfig())
}

func (p serviceChallengeProviderProxy) Create(ctx context.Context, input ChallengeInput) (*Challenge, error) {
	challenge, err := p.provider.Create(ctx, toServiceChallengeInput(input))
	if err != nil {
		return nil, err
	}
	return fromServiceChallenge(challenge), nil
}

func (p serviceChallengeProviderProxy) Verify(ctx context.Context, answer ChallengeAnswer, input ChallengeInput) error {
	return p.provider.Verify(ctx, toServiceChallengeAnswer(answer), toServiceChallengeInput(input))
}
