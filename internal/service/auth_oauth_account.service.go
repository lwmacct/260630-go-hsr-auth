package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lwmacct/260630-go-hsr-auth/internal/repository"
)

type AuthOauthAccountService struct {
	store *repository.Store
}

func NewAuthOauthAccountService(store *repository.Store) *AuthOauthAccountService {
	if store == nil {
		panic("NewAuthOauthAccountService: store is nil")
	}
	return &AuthOauthAccountService{store: store}
}

func (s *AuthOauthAccountService) ByProviderSubject(ctx context.Context, provider string, subject string) (*AuthOauthAccount, error) {
	row, err := s.store.FetchAuthOauthAccountByProviderSubject(ctx, strings.TrimSpace(provider), strings.TrimSpace(subject))
	if err != nil {
		return nil, err
	}
	identity := utilAuthOauthAccount(*row)
	return &identity, nil
}

func (s *AuthOauthAccountService) ResolveUser(ctx context.Context, input AuthOauthAccountResolveInput) (*User, error) {
	profile := input.Profile
	identity, err := s.ByProviderSubject(ctx, profile.Provider, profile.Subject)
	if err == nil {
		users := NewUserService(s.store)
		user, userErr := users.ByID(ctx, identity.UserID)
		if userErr != nil {
			return nil, userErr
		}
		if activeErr := users.EnsureActive(user); activeErr != nil {
			return nil, activeErr
		}
		return user, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	if !input.AutoRegister {
		return nil, ErrAuthOauthAccountRegistrationDisabled
	}
	var user *User
	err = s.store.RunInTx(ctx, func(ctx context.Context, txStore *repository.Store) error {
		users := NewUserService(txStore)
		oauthAccounts := NewAuthOauthAccountService(txStore)
		created, createErr := users.CreateWithUniqueUsername(ctx, CreateUserInput{
			Username:    input.Username,
			DisplayName: profile.ProviderDisplayName,
			Email:       profile.ProviderEmail,
			AvatarURL:   profile.ProviderAvatarURL,
		})
		if createErr != nil {
			return createErr
		}
		_, identityErr := oauthAccounts.Create(ctx, created.ID, profile)
		if identityErr != nil {
			return identityErr
		}
		user = created
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthOauthAccountService) Create(ctx context.Context, userID int64, profile AuthOauthAccountProfile) (*AuthOauthAccount, error) {
	now := time.Now().UTC()
	item := repository.AuthOauthAccountCreate{
		UserID:                userID,
		Provider:              strings.TrimSpace(profile.Provider),
		Subject:               strings.TrimSpace(profile.Subject),
		ProviderEmail:         strings.TrimSpace(profile.ProviderEmail),
		ProviderEmailVerified: profile.ProviderEmailVerified,
		ProviderDisplayName:   strings.TrimSpace(profile.ProviderDisplayName),
		ProviderAvatarURL:     strings.TrimSpace(profile.ProviderAvatarURL),
		ProviderProfile:       profile.ProviderProfile,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	if item.Provider == "" || item.Subject == "" || item.UserID == 0 {
		return nil, ErrAuthOauthAccountTaken
	}
	row, err := s.store.CreateAuthOauthAccount(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAuthOauthAccountTaken, err)
	}
	identity := utilAuthOauthAccount(*row)
	return &identity, nil
}

func (s *AuthOauthAccountService) DeleteForUsers(ctx context.Context, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}
	return s.store.DeleteAuthOauthAccountForUsers(ctx, userIDs)
}
