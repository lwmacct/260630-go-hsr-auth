package auth

import (
	"context"

	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/requestctx"
)

func ContextWithRequest(ctx context.Context, request SessionRequest) context.Context {
	return requestctx.ContextWithRequest(ctx, request)
}

func RequestFromContext(ctx context.Context) (SessionRequest, bool) {
	return requestctx.RequestFromContext(ctx)
}

func toServiceSessionRequest(value SessionRequest) service.AuthSessionInput {
	return service.AuthSessionInput{
		IP:         value.IP,
		Scheme:     value.Scheme,
		Host:       value.Host,
		UserAgent:  value.UserAgent,
		Method:     value.Method,
		Path:       value.Path,
		RemoteAddr: value.RemoteAddr,
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
