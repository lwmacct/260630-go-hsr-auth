package handler

import "github.com/lwmacct/260630-go-hsr-auth/internal/service"

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

func ToAuthChallengeCreateDTO(challenge *service.AuthChallenge) AuthChallengeCreateDTO {
	if challenge == nil {
		return AuthChallengeCreateDTO{}
	}
	return AuthChallengeCreateDTO{
		Provider:    challenge.Provider,
		ChallengeID: challenge.ChallengeID,
		Image:       challenge.Image,
		ExpiresAt:   challenge.ExpiresAt,
	}
}

func ToAuthChallengeAnswer(dto AuthChallengeDTO) service.AuthChallengeAnswer {
	return service.AuthChallengeAnswer{
		Provider:    dto.Provider,
		ChallengeID: dto.ChallengeID,
		Answer:      dto.Answer,
		Token:       dto.Token,
	}
}

func ToAuthChallengeInput(request service.AuthSessionInput) service.AuthChallengeInput {
	return service.AuthChallengeInput{
		IP:         request.IP,
		UserAgent:  request.UserAgent,
		Method:     request.Method,
		Path:       request.Path,
		RemoteAddr: request.RemoteAddr,
	}
}
