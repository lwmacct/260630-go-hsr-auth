package service

import (
	"errors"
	"strings"
	"time"

	"github.com/lwmacct/260630-go-hsr-auth/internal/repository"
)

const (
	UserRoleUser  = "user"
	UserRoleAdmin = "admin"
)

const (
	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"
)

type User struct {
	ID          int64
	Username    string
	DisplayName string
	Email       string
	AvatarURL   string
	Role        string
	Status      string
	DisabledAt  *time.Time
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateUserInput struct {
	Username    string
	DisplayName string
	Email       string
	AvatarURL   string
	Role        string
}

type UpdateUserProfileInput struct {
	DisplayName string
	Email       string
	AvatarURL   string
}

type UserListFilter struct {
	Keyword  string
	Role     string
	Status   string
	Page     int
	PageSize int
}

var (
	ErrUserInvalidCredentials = errors.New("invalid credentials")
	ErrUserDisabled           = errors.New("user disabled")
	ErrUserEmptySelection     = errors.New("empty user selection")
	ErrUserInvalidRole        = errors.New("invalid user role")
	ErrUserInvalidStatus      = errors.New("invalid user status")
	ErrUserUsernameTaken      = errors.New("username taken")
)

func utilUser(row repository.UserRow) User {
	return User{
		ID:          row.ID,
		Username:    row.Username,
		DisplayName: row.DisplayName,
		Email:       row.Email,
		AvatarURL:   row.AvatarURL,
		Role:        row.Role,
		Status:      row.Status,
		DisabledAt:  row.DisabledAt,
		LastLoginAt: row.LastLoginAt,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func utilIsUserUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "unique constraint failed") ||
		strings.Contains(message, "duplicate key value violates unique constraint") ||
		strings.Contains(message, "constraint failed: unique")
}

func normalizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	var builder strings.Builder
	lastDash := false
	for _, r := range username {
		valid := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if valid {
			builder.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			builder.WriteByte('-')
			lastDash = true
		}
	}
	username = strings.Trim(builder.String(), "-")
	if len(username) > 64 {
		username = strings.Trim(username[:64], "-")
	}
	return username
}

func validateUserRole(role string) bool {
	return role == UserRoleUser || role == UserRoleAdmin
}

func validateUserStatus(status string) bool {
	return status == UserStatusActive || status == UserStatusDisabled
}

func normalizeUserPage(page int, pageSize int) (int, int) {
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
