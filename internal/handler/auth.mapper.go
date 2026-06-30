package handler

import (
	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/challenge"
)

func ToAuthUserDTO(user *service.User, runtimeAdmin bool) *AuthUserDTO {
	if user == nil {
		return nil
	}
	admin := user.Role == service.UserRoleAdmin || runtimeAdmin
	return &AuthUserDTO{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
		Admin:       admin,
	}
}

func ToAuthChallengeCreateDTO(value *challenge.Challenge) AuthChallengeCreateDTO {
	if value == nil {
		return AuthChallengeCreateDTO{}
	}
	return AuthChallengeCreateDTO{
		Provider:    value.Provider,
		ChallengeID: value.ChallengeID,
		Image:       value.Image,
		ExpiresAt:   value.ExpiresAt,
	}
}

func ToAuthChallengeAnswer(dto AuthChallengeDTO) challenge.Answer {
	return challenge.Answer{
		Provider:    dto.Provider,
		ChallengeID: dto.ChallengeID,
		Answer:      dto.Answer,
		Token:       dto.Token,
	}
}

func ToAuthChallengeInput(request service.AuthSessionInput) challenge.Input {
	return challenge.Input{
		IP:         request.IP,
		UserAgent:  request.UserAgent,
		Method:     request.Method,
		Path:       request.Path,
		RemoteAddr: request.RemoteAddr,
	}
}
