package auth

import (
	"context"

	"github.com/lwmacct/260630-go-hsr-auth/internal/handler"
	"github.com/lwmacct/260630-go-hsr-auth/internal/infra/oauthclient"
	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/challenge"
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
	return challenge.NewImageProvider(maxItems)
}

func NewRemoteTokenChallengeProvider(provider string, siteKey string, secret string, verifyURL string) (ChallengeProvider, error) {
	return challenge.NewRemoteTokenProvider(provider, siteKey, secret, verifyURL)
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
