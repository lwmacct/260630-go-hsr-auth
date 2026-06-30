package auth

import "github.com/lwmacct/260630-go-hsr-shared/pkg/challenge"

func NewImageChallengeProvider(maxItems int) ChallengeProvider {
	return challenge.NewImageProvider(maxItems)
}

func NewRemoteTokenChallengeProvider(provider string, siteKey string, secret string, verifyURL string) (ChallengeProvider, error) {
	return challenge.NewRemoteTokenProvider(provider, siteKey, secret, verifyURL)
}
