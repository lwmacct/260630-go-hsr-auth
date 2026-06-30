package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/challenge"
	"github.com/uptrace/bun"

	"github.com/lwmacct/260630-go-hsr-auth/internal/handler"
	"github.com/lwmacct/260630-go-hsr-auth/internal/repository"
	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
)

type DB = bun.IDB

type Module struct {
	endpoint *handler.Endpoint
}

func New(options Options) (*Module, error) {
	if options.DB == nil {
		panic("auth.New: DB is nil")
	}
	if options.SessionTTL <= 0 {
		options.SessionTTL = service.AuthSessionDefaultTTL
	}
	if options.OAuthTTL <= 0 {
		options.OAuthTTL = 10 * time.Minute
	}
	if options.Config.ChallengeProvider == nil {
		options.Config.ChallengeProvider = NewImageChallengeProvider(1024)
	}
	if options.Config.Request == nil {
		options.Config.Request = RequestFromContext
	}
	options.Config.Session.Cookie = normalizeSessionCookieConfig(options.Config.Session.Cookie)

	store := repository.NewStore(options.DB)
	users := service.NewUserService(store)
	passwords := service.NewAuthPasswordService(store)
	oauthAccounts := service.NewAuthOauthAccountService(store)
	sessions := service.NewAuthSessionService(store, options.SessionTTL)
	oauthFlows := service.NewAuthOauthFlowService(store, options.OAuthTTL)
	challenges := challenge.NewService(options.Config.ChallengeProvider)
	adminUsers := service.NewAdminUserService(
		store,
		users,
		passwords,
		oauthAccounts,
		sessions,
		runtimeAdminChecker(options.Config.RuntimeAdmins),
	)

	endpoint := handler.NewEndpoint(toHandlerConfig(options.Config), handler.Services{
		Users:         users,
		Passwords:     passwords,
		OAuthAccounts: oauthAccounts,
		Sessions:      sessions,
		OAuthFlows:    oauthFlows,
		Challenges:    challenges,
		AdminUsers:    adminUsers,
	})
	return &Module{endpoint: endpoint}, nil
}

func MustNew(options Options) *Module {
	module, err := New(options)
	if err != nil {
		panic(err)
	}
	return module
}

func (m *Module) Handler() http.Handler {
	return m.endpoint.Handler()
}

func (m *Module) Register(api huma.API) {
	m.endpoint.Register(api)
}

func toHandlerConfig(config Config) handler.Config {
	return handler.Config{
		LocalLoginEnabled:        config.Local.LoginEnabled,
		LocalRegistrationEnabled: config.Local.RegistrationEnabled,
		OAuthEnabled:             config.OAuth.Enabled,
		OAuthAutoRegister:        config.OAuth.AutoRegister,
		OAuthCallbackBaseURL:     config.OAuth.CallbackBaseURL,
		OAuthProviders:           toHandlerOAuthProviderConfigs(config.OAuth.Providers),
		SessionCookie:            toHandlerSessionCookieConfig(config.Session.Cookie),
		RuntimeAdmins:            config.RuntimeAdmins,
		Request:                  toHandlerRequest(config.Request),
		OAuthProvider:            toHandlerOAuthProvider(config.OAuthProvider),
	}
}

func normalizeSessionCookieConfig(config SessionCookieConfig) SessionCookieConfig {
	if config.Name == "" {
		config.Name = "web_session"
	}
	if config.Path == "" {
		config.Path = "/api"
	}
	if config.SameSite == 0 {
		config.SameSite = http.SameSiteStrictMode
	}
	return config
}

func toHandlerSessionCookieConfig(config SessionCookieConfig) handler.SessionCookieConfig {
	return handler.SessionCookieConfig{
		Name:     config.Name,
		Path:     config.Path,
		Secure:   config.Secure,
		SameSite: config.SameSite,
	}
}

func toHandlerOAuthProviderConfigs(values []OAuthProviderConfig) []handler.OAuthProviderConfig {
	items := make([]handler.OAuthProviderConfig, 0, len(values))
	for _, value := range values {
		items = append(items, handler.OAuthProviderConfig{
			Provider: value.Provider,
			Label:    value.Label,
		})
	}
	return items
}

func toHandlerRequest(fn RequestFunc) handler.RequestAuth {
	return func(ctx context.Context) (service.AuthSessionInput, bool) {
		request, ok := fn(ctx)
		return toServiceSessionRequest(request), ok
	}
}

func toHandlerOAuthProvider(fn func(string) (OAuthProvider, error)) func(string) (handler.OAuthProviderAuth, error) {
	if fn == nil {
		return nil
	}
	return func(provider string) (handler.OAuthProviderAuth, error) {
		value, err := fn(provider)
		if err != nil {
			return nil, err
		}
		return handlerOAuthProvider{provider: value}, nil
	}
}

func runtimeAdminChecker(admins []string) func(string) bool {
	return func(username string) bool {
		for _, admin := range admins {
			if strings.EqualFold(strings.TrimSpace(admin), username) {
				return true
			}
		}
		return false
	}
}
