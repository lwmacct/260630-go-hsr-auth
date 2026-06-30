package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/lwmacct/260630-go-hsr-auth/internal/repository"
)

type AuthPasswordService struct {
	store *repository.Store
}

func NewAuthPasswordService(store *repository.Store) *AuthPasswordService {
	if store == nil {
		panic("NewAuthPasswordService: store is nil")
	}
	return &AuthPasswordService{store: store}
}

func (s *AuthPasswordService) CheckStrength(username string, password string) error {
	return validateAuthPassword(username, password)
}

func (s *AuthPasswordService) Register(ctx context.Context, input AuthPasswordRegisterInput) (*User, error) {
	if err := validateAuthPassword(input.Username, input.Password); err != nil {
		return nil, err
	}
	var user *User
	err := s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		users := NewUserService(txStore)
		passwords := NewAuthPasswordService(txStore)
		created, err := users.Create(ctx, CreateUserInput{Username: input.Username})
		if err != nil {
			return err
		}
		if err := passwords.Set(ctx, created.Username, created.ID, input.Password); err != nil {
			return err
		}
		user = created
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthPasswordService) Set(ctx context.Context, username string, userID int64, password string) error {
	if err := validateAuthPassword(username, password); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	_, err = s.store.CreateAuthPassword(ctx, repository.AuthPasswordCreate{
		UserID:            userID,
		PasswordHash:      string(hash),
		PasswordChangedAt: now,
		CreatedAt:         now,
		UpdatedAt:         now,
	})
	return err
}

func (s *AuthPasswordService) Reset(ctx context.Context, username string, userID int64, password string) error {
	if err := validateAuthPassword(username, password); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	_, err = s.store.UpsertAuthPassword(ctx, repository.AuthPasswordCreate{
		UserID:            userID,
		PasswordHash:      string(hash),
		PasswordChangedAt: now,
		CreatedAt:         now,
		UpdatedAt:         now,
	})
	return err
}

func (s *AuthPasswordService) DeleteForUsers(ctx context.Context, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}
	return s.store.DeleteAuthPasswordForUsers(ctx, userIDs)
}

func (s *AuthPasswordService) Change(ctx context.Context, username string, userID int64, currentPassword string, newPassword string) error {
	if currentPassword == "" {
		return ErrAuthPasswordInvalidCredentials
	}
	row, err := s.store.FetchAuthPassword(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrAuthPasswordInvalidCredentials
		}
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(row.PasswordHash), []byte(currentPassword))
	if err != nil {
		return ErrAuthPasswordInvalidCredentials
	}
	err = validateAuthPassword(username, newPassword)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.store.UpdateAuthPassword(ctx, userID, string(hash), time.Now().UTC())
	return err
}

func (s *AuthPasswordService) Authenticate(ctx context.Context, username string, password string, users *UserService) (*User, error) {
	if strings.TrimSpace(username) == "" || password == "" {
		return nil, ErrAuthPasswordInvalidCredentials
	}
	user, err := users.ByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAuthPasswordInvalidCredentials
		}
		return nil, err
	}
	err = users.EnsureActive(user)
	if err != nil {
		return nil, err
	}
	row, err := s.store.FetchAuthPassword(ctx, user.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAuthPasswordInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(row.PasswordHash), []byte(password)); err != nil {
		return nil, ErrAuthPasswordInvalidCredentials
	}
	return user, nil
}
