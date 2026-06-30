package service

import (
	"errors"
	"time"
)

type AdminUser struct {
	ID          string
	Username    string
	DisplayName string
	Email       string
	AvatarURL   string
	Role        string
	Status      string
	Admin       bool
	DisabledAt  *time.Time
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AdminUserListInput struct {
	Keyword  string
	Role     string
	Status   string
	Page     int
	PageSize int
}

type AdminUserListOutput struct {
	Items    []AdminUser
	Total    int
	Page     int
	PageSize int
}

type AdminUserCreateInput struct {
	Username    string
	DisplayName string
	Email       string
	AvatarURL   string
	Role        string
	Password    string
}

type AdminUserUpdateProfileInput struct {
	DisplayName string
	Email       string
	AvatarURL   string
}

var (
	ErrAdminUserCannotOperateSelf = errors.New("cannot operate current user")
	ErrAdminUserEmptySelection    = errors.New("empty user selection")
	ErrAdminUserNotFound          = errors.New("user not found")
)

func utilAdminUser(user User, runtimeAdmin bool) AdminUser {
	return AdminUser{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
		Status:      user.Status,
		Admin:       user.Role == UserRoleAdmin || runtimeAdmin,
		DisabledAt:  user.DisabledAt,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func validateAdminUserSelection(actorID string, ids []string) error {
	if len(ids) == 0 {
		return ErrAdminUserEmptySelection
	}
	for _, id := range ids {
		if id == actorID {
			return ErrAdminUserCannotOperateSelf
		}
	}
	return nil
}

func normalizeAdminUserPage(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
