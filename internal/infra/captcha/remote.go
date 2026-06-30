package captcha

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
)

type RemoteTokenProvider struct {
	provider  string
	siteKey   string
	secret    string
	verifyURL string
	client    *http.Client
}

func NewRemoteTokenProvider(provider string, siteKey string, secret string, verifyURL string) (*RemoteTokenProvider, error) {
	if strings.TrimSpace(provider) == "" || strings.TrimSpace(siteKey) == "" || strings.TrimSpace(secret) == "" || strings.TrimSpace(verifyURL) == "" {
		return nil, service.ErrAuthChallengeUnsupported
	}
	return &RemoteTokenProvider{
		provider:  strings.TrimSpace(provider),
		siteKey:   strings.TrimSpace(siteKey),
		secret:    strings.TrimSpace(secret),
		verifyURL: strings.TrimSpace(verifyURL),
		client:    &http.Client{Timeout: 5 * time.Second},
	}, nil
}

func (p *RemoteTokenProvider) Name() string {
	return p.provider
}

func (p *RemoteTokenProvider) PublicConfig() service.AuthChallengePublicConfig {
	return service.AuthChallengePublicConfig{Provider: p.provider, SiteKey: p.siteKey}
}

func (p *RemoteTokenProvider) Create(context.Context, service.AuthChallengeInput) (*service.AuthChallenge, error) {
	return nil, service.ErrAuthChallengeUnsupported
}

func (p *RemoteTokenProvider) Verify(ctx context.Context, response service.AuthChallengeAnswer, request service.AuthChallengeInput) error {
	if strings.TrimSpace(response.Token) == "" {
		return service.ErrAuthChallengeInvalid
	}
	form := url.Values{}
	form.Set("secret", p.secret)
	form.Set("response", response.Token)
	if request.IP != "" {
		form.Set("remoteip", request.IP)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.verifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return service.ErrAuthChallengeInvalid
	}
	var body struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}
	if !body.Success {
		return service.ErrAuthChallengeInvalid
	}
	return nil
}
