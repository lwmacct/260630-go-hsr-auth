package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lwmacct/260630-go-hsr-auth/internal/repository"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/challenge"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/identity"
	"github.com/uptrace/bun"

	"github.com/lwmacct/260630-go-hsr-auth/internal/handler"
	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
)

type DB = bun.IDB

type Module struct {
	endpoint       *handler.Endpoint
	users          *service.UserService
	sessions       *service.AuthSessionService
	runtimeAdminFn func(string) bool
	sessionCookie  SessionCookieConfig
}

func New(options Options) (*Module, error) {
	if options.DB == nil {
		panic("auth.New: DB is nil")
	}
	if options.SessionTTL <= 0 {
		options.SessionTTL = service.AuthSessionDefaultTTL
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
	sessions := service.NewAuthSessionService(store, options.SessionTTL)
	challenges := challenge.NewService(options.Config.ChallengeProvider)
	runtimeAdminFn := runtimeAdminChecker(options.Config.RuntimeAdmins)
	adminUsers := service.NewAdminUserService(
		store,
		users,
		passwords,
		sessions,
		runtimeAdminFn,
	)

	endpoint := handler.NewEndpoint(toHandlerConfig(options.Config), handler.Services{
		Users:      users,
		Passwords:  passwords,
		Sessions:   sessions,
		Challenges: challenges,
		AdminUsers: adminUsers,
	})
	return &Module{
		endpoint:       endpoint,
		users:          users,
		sessions:       sessions,
		runtimeAdminFn: runtimeAdminFn,
		sessionCookie:  options.Config.Session.Cookie,
	}, nil
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

func (m *Module) UserByID(ctx context.Context, id string) (*User, error) {
	user, err := m.users.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return m.toUser(user), nil
}

func (m *Module) Principal(ctx context.Context, subject string) (*identity.Principal, error) {
	user, err := m.users.ByID(ctx, subject)
	if err != nil {
		return nil, err
	}
	return toPrincipal(m.toUser(user)), nil
}

func (m *Module) Principals(ctx context.Context, subjects []string) (map[string]*identity.Principal, error) {
	out := make(map[string]*identity.Principal, len(subjects))
	for _, subject := range subjects {
		if subject == "" {
			continue
		}
		if _, ok := out[subject]; ok {
			continue
		}
		principal, err := m.Principal(ctx, subject)
		if err != nil {
			continue
		}
		out[subject] = principal
	}
	return out, nil
}

func (m *Module) CreateExternalUser(ctx context.Context, input ExternalUserInput) (*User, error) {
	user, err := m.users.CreateWithUniqueUsername(ctx, service.CreateUserInput{
		Username:    input.Username,
		DisplayName: input.DisplayName,
		Email:       input.Email,
		AvatarURL:   input.AvatarURL,
	})
	if err != nil {
		return nil, err
	}
	return m.toUser(user), nil
}

func (m *Module) CreateSession(ctx context.Context, userID string, request SessionRequest) (*Session, error) {
	sessionID, expiresAt, err := m.sessions.Create(ctx, userID, toServiceSessionRequest(request))
	if err != nil {
		return nil, err
	}
	if err := m.users.MarkLogin(ctx, userID); err != nil {
		return nil, err
	}
	user, err := m.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &Session{
		ID:        sessionID,
		ExpiresAt: expiresAt,
		SetCookie: sessionCookieValue(sessionID, expiresAt, m.sessionCookie),
		User:      user,
	}, nil
}

func (m *Module) CurrentUser(ctx context.Context, sessionID string, request SessionRequest) (*User, error) {
	sessionUser, err := m.sessions.User(ctx, sessionID, toServiceSessionRequest(request), m.users)
	if err != nil {
		return nil, err
	}
	return m.UserByID(ctx, sessionUser.ID)
}

func (m *Module) CurrentPrincipal(ctx context.Context, sessionID string, request SessionRequest) (*identity.Principal, error) {
	user, err := m.CurrentUser(ctx, sessionID, request)
	if err != nil {
		return nil, err
	}
	return toPrincipal(user), nil
}

func (m *Module) RequireAdmin(ctx context.Context, sessionID string, request SessionRequest) (*User, error) {
	user, err := m.CurrentUser(ctx, sessionID, request)
	if err != nil {
		return nil, err
	}
	if !user.Admin {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

func toPrincipal(user *User) *identity.Principal {
	if user == nil {
		return nil
	}
	return &identity.Principal{
		ID:          user.ID,
		Subject:     user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
		Status:      user.Status,
		Admin:       user.Admin,
		DisabledAt:  user.DisabledAt,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func toHandlerConfig(config Config) handler.Config {
	return handler.Config{
		LocalLoginEnabled:        config.Local.LoginEnabled,
		LocalRegistrationEnabled: config.Local.RegistrationEnabled,
		SessionCookie:            toHandlerSessionCookieConfig(config.Session.Cookie),
		RuntimeAdmins:            config.RuntimeAdmins,
		Request:                  toHandlerRequest(config.Request),
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

func toHandlerRequest(fn RequestFunc) handler.RequestAuth {
	return func(ctx context.Context) (service.AuthSessionInput, bool) {
		request, ok := fn(ctx)
		return toServiceSessionRequest(request), ok
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

func (m *Module) toUser(user *service.User) *User {
	if user == nil {
		return nil
	}
	return &User{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
		Status:      user.Status,
		Admin:       user.Role == service.UserRoleAdmin || m.runtimeAdminFn(user.Username),
		DisabledAt:  user.DisabledAt,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
