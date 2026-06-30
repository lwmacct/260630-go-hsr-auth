package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lwmacct/260630-go-hsr-auth/internal/repository"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/idgen"
)

type UserService struct {
	store *repository.Store
}

func NewUserService(store *repository.Store) *UserService {
	if store == nil {
		panic("NewUserService: store is nil")
	}
	return &UserService{store: store}
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	username := normalizeUsername(input.Username)
	if username == "" {
		return nil, ErrUserInvalidCredentials
	}
	displayName := strings.TrimSpace(input.DisplayName)
	if displayName == "" {
		displayName = username
	}
	role := strings.TrimSpace(input.Role)
	if role == "" {
		role = UserRoleUser
	}
	if !validateUserRole(role) {
		return nil, ErrUserInvalidRole
	}
	now := time.Now().UTC()
	row, err := s.store.CreateUser(ctx, repository.UserCreate{
		ID:          idgen.NewUUID7(),
		Username:    username,
		DisplayName: displayName,
		Email:       strings.TrimSpace(input.Email),
		AvatarURL:   strings.TrimSpace(input.AvatarURL),
		Role:        role,
		Status:      UserStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		if utilIsUserUniqueViolation(err) {
			return nil, fmt.Errorf("%w: %w", ErrUserUsernameTaken, err)
		}
		return nil, err
	}
	user := utilUser(*row)
	return &user, nil
}

func (s *UserService) CreateWithUniqueUsername(ctx context.Context, input CreateUserInput) (*User, error) {
	base := normalizeUsername(input.Username)
	if base == "" {
		base = "user"
	}
	const maxAttempts = 100
	for attempt := range maxAttempts {
		candidate := base
		if attempt > 0 {
			candidate = fmt.Sprintf("%s-%d", base, attempt+1)
		}
		input.Username = candidate
		user, err := s.Create(ctx, input)
		if err == nil {
			return user, nil
		}
		if !errors.Is(err, ErrUserUsernameTaken) {
			return nil, err
		}
	}
	return nil, ErrUserUsernameTaken
}

func (s *UserService) ByID(ctx context.Context, id string) (*User, error) {
	row, err := s.store.FetchUser(ctx, id)
	if err != nil {
		return nil, err
	}
	user := utilUser(*row)
	return &user, nil
}

func (s *UserService) ByUsername(ctx context.Context, username string) (*User, error) {
	row, err := s.store.FetchUserByUsername(ctx, normalizeUsername(username))
	if err != nil {
		return nil, err
	}
	user := utilUser(*row)
	return &user, nil
}

func (s *UserService) EnsureActive(user *User) error {
	if user == nil {
		return repository.ErrNotFound
	}
	if user.Status == UserStatusDisabled || user.DisabledAt != nil {
		return ErrUserDisabled
	}
	return nil
}

func (s *UserService) List(ctx context.Context) ([]User, error) {
	rows, err := s.store.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(rows))
	for _, row := range rows {
		users = append(users, utilUser(row))
	}
	return users, nil
}

func (s *UserService) Search(ctx context.Context, filter UserListFilter) ([]User, int, error) {
	page, pageSize := normalizeUserPage(filter.Page, filter.PageSize)
	repoFilter := repository.UserFilter{
		Keyword:  filter.Keyword,
		Role:     filter.Role,
		Status:   filter.Status,
		Page:     page,
		PageSize: pageSize,
	}
	rows, err := s.store.ListUsersByFilter(ctx, repoFilter)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.store.FetchUserTotalByFilter(ctx, repoFilter)
	if err != nil {
		return nil, 0, err
	}
	users := make([]User, 0, len(rows))
	for _, row := range rows {
		users = append(users, utilUser(row))
	}
	return users, total.Count, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, id string, input UpdateUserProfileInput) (*User, error) {
	displayName := strings.TrimSpace(input.DisplayName)
	if displayName == "" {
		return nil, ErrUserInvalidCredentials
	}
	row, err := s.store.UpdateUserProfile(ctx, id, repository.UserProfilePatch{
		DisplayName: displayName,
		Email:       strings.TrimSpace(input.Email),
		AvatarURL:   strings.TrimSpace(input.AvatarURL),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}
	user := utilUser(*row)
	return &user, nil
}

func (s *UserService) SetRoleBatch(ctx context.Context, ids []string, role string) error {
	role = strings.TrimSpace(role)
	if len(ids) == 0 {
		return ErrUserEmptySelection
	}
	if !validateUserRole(role) {
		return ErrUserInvalidRole
	}
	_, err := s.store.UpdateUserRoleBatch(ctx, ids, role, time.Now().UTC())
	return err
}

func (s *UserService) SetStatusBatch(ctx context.Context, ids []string, status string) error {
	status = strings.TrimSpace(status)
	if len(ids) == 0 {
		return ErrUserEmptySelection
	}
	if !validateUserStatus(status) {
		return ErrUserInvalidStatus
	}
	now := time.Now().UTC()
	var disabledAt *time.Time
	if status == UserStatusDisabled {
		disabledAt = &now
	}
	_, err := s.store.UpdateUserStatusBatch(ctx, ids, status, now, disabledAt)
	return err
}

func (s *UserService) MarkLogin(ctx context.Context, id string) error {
	_, err := s.store.UpdateUserLastLogin(ctx, id, time.Now().UTC())
	return err
}

func (s *UserService) DeleteBatch(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return ErrUserEmptySelection
	}
	return s.store.DeleteUserBatch(ctx, ids)
}
