package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	challengepkg "github.com/lwmacct/260630-go-hsr-shared/pkg/challenge"

	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
)

type authHandler struct {
	config   Config
	services Services
}

func RegisterAuth(api huma.API, config Config, services Services) {
	handler := authHandler{config: config, services: services}
	auth := huma.NewGroup(api, "/auth")
	huma.Register(auth, huma.Operation{OperationID: "get-auth-config", Method: http.MethodGet, Path: "/config", Tags: []string{"Auth"}}, handler.configOutput)
	huma.Register(auth, huma.Operation{OperationID: "create-auth-challenge", Method: http.MethodPost, Path: "/challenges", Tags: []string{"Auth"}}, handler.createChallenge)
	huma.Register(auth, huma.Operation{OperationID: "register-password-user", Method: http.MethodPost, Path: "/password/register", DefaultStatus: http.StatusCreated, Tags: []string{"Auth"}}, handler.passwordRegister)
	huma.Register(auth, huma.Operation{OperationID: "login-password", Method: http.MethodPost, Path: "/password/login", Tags: []string{"Auth"}}, handler.passwordLogin)
	huma.Register(auth, huma.Operation{OperationID: "change-password", Method: http.MethodPost, Path: "/password/change", Tags: []string{"Auth"}}, handler.passwordChange)
	huma.Register(auth, huma.Operation{OperationID: "logout", Method: http.MethodPost, Path: "/logout", Tags: []string{"Auth"}}, handler.logout)
	huma.Register(auth, huma.Operation{OperationID: "get-current-user", Method: http.MethodGet, Path: "/me", Tags: []string{"Auth"}}, handler.me)
}

func (h authHandler) configOutput(_ context.Context, _ *struct{}) (*BodyDTO[AuthConfigDTO], error) {
	challenge := h.services.Challenges.PublicConfig()
	body := AuthConfigDTO{}
	body.Local.LoginEnabled = h.config.LocalLoginEnabled
	body.Local.RegistrationEnabled = h.config.LocalRegistrationEnabled
	body.Challenge.Provider = challenge.Provider
	body.Challenge.SiteKey = challenge.SiteKey
	return &BodyDTO[AuthConfigDTO]{Body: body}, nil
}

func (h authHandler) createChallenge(ctx context.Context, _ *struct{}) (*BodyDTO[AuthChallengeCreateDTO], error) {
	request, err := utilRequest(ctx, h.config)
	if err != nil {
		return nil, err
	}
	challenge, err := h.services.Challenges.Create(ctx, ToAuthChallengeInput(request))
	if err != nil {
		if errors.Is(err, challengepkg.ErrLimitExceeded) {
			return nil, huma.Error429TooManyRequests("too many challenges")
		}
		return nil, huma.Error400BadRequest("challenge creation unsupported")
	}
	return &BodyDTO[AuthChallengeCreateDTO]{Body: ToAuthChallengeCreateDTO(challenge)}, nil
}

func (h authHandler) passwordRegister(ctx context.Context, input *BodyInputDTO[AuthCredentialsDTO]) (*AuthSessionResponseDTO, error) {
	if !h.config.LocalLoginEnabled || !h.config.LocalRegistrationEnabled {
		return nil, huma.Error403Forbidden("password registration disabled")
	}
	request, err := utilRequest(ctx, h.config)
	if err != nil {
		return nil, err
	}
	challengeErr := h.verifyChallenge(ctx, input.Body.Challenge, request)
	if challengeErr != nil {
		return nil, huma.Error401Unauthorized("invalid challenge")
	}
	user, err := h.services.Passwords.Register(ctx, service.AuthPasswordRegisterInput{
		Username: input.Body.Username,
		Password: input.Body.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserUsernameTaken) || errors.Is(err, service.ErrAuthPasswordWeakPassword) || errors.Is(err, service.ErrUserInvalidCredentials) {
			return nil, huma.Error400BadRequest(utilRegisterErrorMessage(err))
		}
		return nil, huma.Error500InternalServerError("register failed")
	}
	return h.createSessionResponse(ctx, user.ID, request)
}

func (h authHandler) passwordLogin(ctx context.Context, input *BodyInputDTO[AuthCredentialsDTO]) (*AuthSessionResponseDTO, error) {
	if !h.config.LocalLoginEnabled {
		return nil, huma.Error403Forbidden("password login disabled")
	}
	request, err := utilRequest(ctx, h.config)
	if err != nil {
		return nil, err
	}
	challengeErr := h.verifyChallenge(ctx, input.Body.Challenge, request)
	if challengeErr != nil {
		return nil, huma.Error401Unauthorized("invalid challenge")
	}
	user, err := h.services.Passwords.Authenticate(ctx, input.Body.Username, input.Body.Password, h.services.Users)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	return h.createSessionResponse(ctx, user.ID, request)
}

func (h authHandler) passwordChange(ctx context.Context, input *AuthPasswordChangeInputDTO) (*BodyDTO[AuthSessionDTO], error) {
	if !h.config.LocalLoginEnabled {
		return nil, huma.Error403Forbidden("password login disabled")
	}
	if input.Session == "" {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	request, err := utilRequest(ctx, h.config)
	if err != nil {
		return nil, err
	}
	sessionUser, err := h.services.Sessions.User(ctx, input.Session, request, h.services.Users)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	user, err := h.services.Users.ByID(ctx, sessionUser.ID)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	if err := h.services.Passwords.Change(ctx, user.Username, user.ID, input.Body.CurrentPassword, input.Body.NewPassword); err != nil {
		if errors.Is(err, service.ErrAuthPasswordInvalidCredentials) {
			return nil, huma.Error401Unauthorized("current password is incorrect")
		}
		if errors.Is(err, service.ErrAuthPasswordWeakPassword) {
			return nil, huma.Error400BadRequest("weak password")
		}
		return nil, huma.Error500InternalServerError("password change failed")
	}
	return &BodyDTO[AuthSessionDTO]{Body: AuthSessionDTO{
		Authenticated: true,
		ExpiresAt:     sessionUser.ExpiresAt,
		User:          h.toAuthUserDTO(user),
	}}, nil
}

func (h authHandler) logout(ctx context.Context, input *AuthLogoutInputDTO) (*AuthSessionResponseDTO, error) {
	if input.Session != "" {
		_ = h.services.Sessions.Delete(ctx, input.Session)
	}
	return &AuthSessionResponseDTO{
		SetCookie: utilClearSessionCookie(h.config.SessionCookie),
		Body:      AuthSessionDTO{Authenticated: false},
	}, nil
}

func (h authHandler) me(ctx context.Context, input *AuthLogoutInputDTO) (*BodyDTO[AuthSessionDTO], error) {
	if input.Session == "" {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	request, err := utilRequest(ctx, h.config)
	if err != nil {
		return nil, err
	}
	sessionUser, err := h.services.Sessions.User(ctx, input.Session, request, h.services.Users)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	user, err := h.services.Users.ByID(ctx, sessionUser.ID)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	return &BodyDTO[AuthSessionDTO]{Body: AuthSessionDTO{
		Authenticated: true,
		ExpiresAt:     sessionUser.ExpiresAt,
		User:          h.toAuthUserDTO(user),
	}}, nil
}

func (h authHandler) verifyChallenge(ctx context.Context, challenge AuthChallengeDTO, request service.AuthSessionInput) error {
	return h.services.Challenges.Verify(ctx, ToAuthChallengeAnswer(challenge), ToAuthChallengeInput(request))
}

func (h authHandler) createSessionResponse(ctx context.Context, userID int64, request service.AuthSessionInput) (*AuthSessionResponseDTO, error) {
	sessionID, expiresAt, err := h.services.Sessions.Create(ctx, userID, request)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}
	err = h.services.Users.MarkLogin(ctx, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}
	user, err := h.services.Users.ByID(ctx, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}
	return &AuthSessionResponseDTO{
		SetCookie: utilSessionCookieValue(sessionID, expiresAt, h.config.SessionCookie),
		Body: AuthSessionDTO{
			Authenticated: true,
			ExpiresAt:     expiresAt,
			User:          h.toAuthUserDTO(user),
		},
	}, nil
}

func (h authHandler) toAuthUserDTO(user *service.User) *AuthUserDTO {
	return ToAuthUserDTO(user, h.isRuntimeAdmin(user))
}

func (h authHandler) isRuntimeAdmin(user *service.User) bool {
	if user == nil {
		return false
	}
	return h.isRuntimeAdminUsername(user.Username)
}

func (h authHandler) isRuntimeAdminUsername(username string) bool {
	for _, admin := range h.config.RuntimeAdmins {
		if strings.EqualFold(strings.TrimSpace(admin), username) {
			return true
		}
	}
	return false
}

func (h authHandler) currentUser(ctx context.Context, sessionID string) (*AuthUserDTO, error) {
	if sessionID == "" {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	request, err := utilRequest(ctx, h.config)
	if err != nil {
		return nil, err
	}
	sessionUser, err := h.services.Sessions.User(ctx, sessionID, request, h.services.Users)
	if err != nil {
		return nil, err
	}
	user, err := h.services.Users.ByID(ctx, sessionUser.ID)
	if err != nil {
		return nil, err
	}
	return h.toAuthUserDTO(user), nil
}
