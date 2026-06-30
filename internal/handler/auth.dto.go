package handler

import "time"

type AuthConfigDTO struct {
	Local struct {
		LoginEnabled        bool `json:"loginEnabled"`
		RegistrationEnabled bool `json:"registrationEnabled"`
	} `json:"local"`
	OAuth struct {
		Enabled   bool              `json:"enabled"`
		Providers []AuthProviderDTO `json:"providers"`
	} `json:"oauth"`
	Challenge struct {
		Provider string `json:"provider"`
		SiteKey  string `json:"sitekey,omitempty"`
	} `json:"challenge"`
}

type AuthProviderDTO struct {
	Provider string `json:"provider"`
	Label    string `json:"label"`
}

type AuthChallengeCreateDTO struct {
	Provider    string    `json:"provider"`
	ChallengeID string    `json:"challengeId,omitempty"`
	Image       string    `json:"image,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
}

type AuthChallengeDTO struct {
	Provider    string `json:"provider"`
	ChallengeID string `json:"challengeId,omitempty"`
	Answer      string `json:"answer,omitempty"`
	Token       string `json:"token,omitempty"`
}

type AuthCredentialsDTO struct {
	Username  string           `json:"username"`
	Password  string           `json:"password"`
	Challenge AuthChallengeDTO `json:"challenge"`
}

type AuthPasswordChangeDTO struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type AuthUserDTO struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email,omitempty"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	Role        string `json:"role"`
	Admin       bool   `json:"admin"`
}

type AuthSessionDTO struct {
	Authenticated bool         `json:"authenticated"`
	ExpiresAt     time.Time    `json:"expiresAt,omitempty"`
	User          *AuthUserDTO `json:"user,omitempty"`
}

type AuthSessionResponseDTO struct {
	SetCookie string `header:"Set-Cookie"`
	Body      AuthSessionDTO
}

type AuthLogoutInputDTO struct {
	Session string `cookie:"web_session"`
}

type AuthPasswordChangeInputDTO struct {
	Session string `cookie:"web_session"`
	Body    AuthPasswordChangeDTO
}

type AuthOAuthStartInputDTO struct {
	Provider string `query:"provider"`
	ReturnTo string `query:"returnTo"`
}

type AuthOAuthCallbackInputDTO struct {
	Provider string `query:"provider"`
	Code     string `query:"code"`
	State    string `query:"state"`
}

type AuthOAuthRedirectResponseDTO struct {
	Location  string `header:"Location"`
	SetCookie string `header:"Set-Cookie,omitempty"`
}
