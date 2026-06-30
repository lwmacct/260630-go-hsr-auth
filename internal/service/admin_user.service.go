package service

import (
	"context"
	"errors"
	"strings"

	"github.com/lwmacct/260630-go-hsr-auth/internal/repository"
)

type AdminUserService struct {
	store          *repository.Store
	users          *UserService
	passwords      *AuthPasswordService
	sessions       *AuthSessionService
	runtimeAdminFn func(string) bool
}

func NewAdminUserService(
	store *repository.Store,
	users *UserService,
	passwords *AuthPasswordService,
	sessions *AuthSessionService,
	runtimeAdminFn func(string) bool,
) *AdminUserService {
	if store == nil {
		panic("NewAdminUserService: store is nil")
	}
	if users == nil || passwords == nil || sessions == nil {
		panic("NewAdminUserService: dependency is nil")
	}
	if runtimeAdminFn == nil {
		runtimeAdminFn = func(string) bool { return false }
	}
	return &AdminUserService{
		store:          store,
		users:          users,
		passwords:      passwords,
		sessions:       sessions,
		runtimeAdminFn: runtimeAdminFn,
	}
}

func (s *AdminUserService) List(ctx context.Context, input AdminUserListInput) (*AdminUserListOutput, error) {
	page, pageSize := normalizeAdminUserPage(input.Page, input.PageSize)
	users, total, err := s.users.Search(ctx, UserListFilter{
		Keyword:  input.Keyword,
		Role:     input.Role,
		Status:   input.Status,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}
	items := make([]AdminUser, 0, len(users))
	for _, user := range users {
		items = append(items, utilAdminUser(user, s.runtimeAdminFn(user.Username)))
	}
	return &AdminUserListOutput{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *AdminUserService) Create(ctx context.Context, input AdminUserCreateInput) (*AdminUser, error) {
	var user *User
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		users := NewUserService(txStore)
		passwords := NewAuthPasswordService(txStore)
		created, err := users.Create(ctx, CreateUserInput{
			Username:    input.Username,
			DisplayName: input.DisplayName,
			Email:       input.Email,
			AvatarURL:   input.AvatarURL,
			Role:        input.Role,
		})
		if err != nil {
			return err
		}
		if strings.TrimSpace(input.Password) != "" {
			if err := passwords.Set(ctx, created.Username, created.ID, input.Password); err != nil {
				return err
			}
		}
		user = created
		return nil
	})
	if err != nil {
		return nil, err
	}
	dto := utilAdminUser(*user, s.runtimeAdminFn(user.Username))
	return &dto, nil
}

func (s *AdminUserService) UpdateProfile(ctx context.Context, id int64, input AdminUserUpdateProfileInput) (*AdminUser, error) {
	user, err := s.users.UpdateProfile(ctx, id, UpdateUserProfileInput(input))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAdminUserNotFound
		}
		return nil, err
	}
	dto := utilAdminUser(*user, s.runtimeAdminFn(user.Username))
	return &dto, nil
}

func (s *AdminUserService) SetRoleBatch(ctx context.Context, actorID int64, ids []int64, role string) error {
	if err := validateAdminUserSelection(actorID, ids); err != nil {
		return err
	}
	return s.users.SetRoleBatch(ctx, ids, role)
}

func (s *AdminUserService) SetStatusBatch(ctx context.Context, actorID int64, ids []int64, status string) error {
	if err := validateAdminUserSelection(actorID, ids); err != nil {
		return err
	}
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		users := NewUserService(txStore)
		sessions := NewAuthSessionService(txStore, AuthSessionDefaultTTL)
		if err := users.SetStatusBatch(ctx, ids, status); err != nil {
			return err
		}
		if status == UserStatusDisabled {
			return sessions.RevokeForUsers(ctx, ids)
		}
		return nil
	})
	return err
}

func (s *AdminUserService) ResetPasswordBatch(ctx context.Context, actorID int64, ids []int64, password string) error {
	if err := validateAdminUserSelection(actorID, ids); err != nil {
		return err
	}
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		users := NewUserService(txStore)
		passwords := NewAuthPasswordService(txStore)
		sessions := NewAuthSessionService(txStore, AuthSessionDefaultTTL)
		for _, id := range ids {
			user, err := users.ByID(ctx, id)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return ErrAdminUserNotFound
				}
				return err
			}
			if err := passwords.Reset(ctx, user.Username, user.ID, password); err != nil {
				return err
			}
		}
		return sessions.RevokeForUsers(ctx, ids)
	})
}

func (s *AdminUserService) DeleteBatch(ctx context.Context, actorID int64, ids []int64) error {
	if err := validateAdminUserSelection(actorID, ids); err != nil {
		return err
	}
	return s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		passwords := NewAuthPasswordService(txStore)
		sessions := NewAuthSessionService(txStore, AuthSessionDefaultTTL)
		users := NewUserService(txStore)
		if err := sessions.DeleteForUsers(ctx, ids); err != nil {
			return err
		}
		if err := passwords.DeleteForUsers(ctx, ids); err != nil {
			return err
		}
		return users.DeleteBatch(ctx, ids)
	})
}
