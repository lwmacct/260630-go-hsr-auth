package auth

import (
	"context"

	"github.com/lwmacct/260630-go-hsr-auth/internal/handler"
	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
)

func ContextWithRequest(ctx context.Context, request SessionRequest) context.Context {
	return handler.ContextWithRequest(ctx, toServiceSessionRequest(request))
}

func RequestFromContext(ctx context.Context) (SessionRequest, bool) {
	request, ok := handler.RequestFromContext(ctx)
	return fromServiceSessionRequest(request), ok
}

func toServiceSessionRequest(value SessionRequest) service.AuthSessionInput {
	return service.AuthSessionInput{
		IP:         value.IP,
		Host:       value.Host,
		UserAgent:  value.UserAgent,
		Method:     value.Method,
		Path:       value.Path,
		RemoteAddr: value.RemoteAddr,
	}
}

func fromServiceSessionRequest(value service.AuthSessionInput) SessionRequest {
	return SessionRequest{
		IP:         value.IP,
		Host:       value.Host,
		UserAgent:  value.UserAgent,
		Method:     value.Method,
		Path:       value.Path,
		RemoteAddr: value.RemoteAddr,
	}
}

func toServiceChallengeInput(value ChallengeInput) service.AuthChallengeInput {
	return service.AuthChallengeInput{
		IP:         value.IP,
		UserAgent:  value.UserAgent,
		Method:     value.Method,
		Path:       value.Path,
		RemoteAddr: value.RemoteAddr,
	}
}

func fromServiceChallengeInput(value service.AuthChallengeInput) ChallengeInput {
	return ChallengeInput{
		IP:         value.IP,
		UserAgent:  value.UserAgent,
		Method:     value.Method,
		Path:       value.Path,
		RemoteAddr: value.RemoteAddr,
	}
}

func toServiceChallengeAnswer(value ChallengeAnswer) service.AuthChallengeAnswer {
	return service.AuthChallengeAnswer{
		Provider:    value.Provider,
		ChallengeID: value.ChallengeID,
		Answer:      value.Answer,
		Token:       value.Token,
	}
}

func fromServiceChallengeAnswer(value service.AuthChallengeAnswer) ChallengeAnswer {
	return ChallengeAnswer{
		Provider:    value.Provider,
		ChallengeID: value.ChallengeID,
		Answer:      value.Answer,
		Token:       value.Token,
	}
}

func toServiceChallengePublicConfig(value ChallengePublicConfig) service.AuthChallengePublicConfig {
	return service.AuthChallengePublicConfig{
		Provider: value.Provider,
		SiteKey:  value.SiteKey,
	}
}

func fromServiceChallengePublicConfig(value service.AuthChallengePublicConfig) ChallengePublicConfig {
	return ChallengePublicConfig{
		Provider: value.Provider,
		SiteKey:  value.SiteKey,
	}
}

func toServiceChallenge(value *Challenge) *service.AuthChallenge {
	if value == nil {
		return nil
	}
	return &service.AuthChallenge{
		Provider:    value.Provider,
		ChallengeID: value.ChallengeID,
		Image:       value.Image,
		ExpiresAt:   value.ExpiresAt,
	}
}

func fromServiceChallenge(value *service.AuthChallenge) *Challenge {
	if value == nil {
		return nil
	}
	return &Challenge{
		Provider:    value.Provider,
		ChallengeID: value.ChallengeID,
		Image:       value.Image,
		ExpiresAt:   value.ExpiresAt,
	}
}

func toServiceOAuthProfile(value OAuthProfile) service.AuthOauthAccountProfile {
	return service.AuthOauthAccountProfile{
		Provider:              value.Provider,
		Subject:               value.Subject,
		ProviderEmail:         value.ProviderEmail,
		ProviderEmailVerified: value.ProviderEmailVerified,
		ProviderDisplayName:   value.ProviderDisplayName,
		ProviderAvatarURL:     value.ProviderAvatarURL,
		ProviderProfile:       value.ProviderProfile,
	}
}

func fromServiceOAuthProfile(value service.AuthOauthAccountProfile) OAuthProfile {
	return OAuthProfile{
		Provider:              value.Provider,
		Subject:               value.Subject,
		ProviderEmail:         value.ProviderEmail,
		ProviderEmailVerified: value.ProviderEmailVerified,
		ProviderDisplayName:   value.ProviderDisplayName,
		ProviderAvatarURL:     value.ProviderAvatarURL,
		ProviderProfile:       value.ProviderProfile,
	}
}
