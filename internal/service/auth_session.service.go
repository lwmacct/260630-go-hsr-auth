package service

import (
	"context"
	"time"

	"github.com/lwmacct/260630-go-hsr-auth/internal/repository"
	"github.com/lwmacct/260630-go-hsr-shared/pkg/token"
)

type AuthSessionService struct {
	store *repository.Store
	ttl   time.Duration
}

func NewAuthSessionService(store *repository.Store, ttl time.Duration) *AuthSessionService {
	if store == nil {
		panic("NewAuthSessionService: store is nil")
	}
	if ttl <= 0 {
		ttl = AuthSessionDefaultTTL
	}
	return &AuthSessionService{store: store, ttl: ttl}
}

func (s *AuthSessionService) Create(ctx context.Context, userID int64, request AuthSessionInput) (string, time.Time, error) {
	sessionID, err := token.NewWithPrefix("sess")
	if err != nil {
		return "", time.Time{}, err
	}
	now := time.Now().UTC()
	expiresAt := now.Add(s.ttl)
	_, err = s.store.CreateAuthSession(ctx, repository.AuthSessionCreate{
		TokenHash:     utilAuthSessionTokenHash(sessionID),
		UserID:        userID,
		LoginIP:       request.IP,
		LastIP:        request.IP,
		UserAgentHash: utilAuthSessionTokenHash(request.UserAgent),
		ExpiresAt:     expiresAt,
		CreatedAt:     now,
		LastSeenAt:    now,
	})
	return sessionID, expiresAt, err
}

func (s *AuthSessionService) Delete(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}
	return s.store.DeleteAuthSession(ctx, utilAuthSessionTokenHash(sessionID))
}

func (s *AuthSessionService) RevokeForUsers(ctx context.Context, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}
	_, err := s.store.UpdateAuthSessionRevokedForUsers(ctx, userIDs, time.Now().UTC())
	return err
}

func (s *AuthSessionService) DeleteForUsers(ctx context.Context, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}
	return s.store.DeleteAuthSessionForUsers(ctx, userIDs)
}

func (s *AuthSessionService) User(ctx context.Context, sessionID string, request AuthSessionInput, users *UserService) (*AuthSessionUser, error) {
	tokenHash := utilAuthSessionTokenHash(sessionID)
	row, err := s.store.FetchAuthSession(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	if !row.ExpiresAt.After(now) {
		_ = s.Delete(ctx, sessionID)
		return nil, repository.ErrNotFound
	}
	user, err := users.ByID(ctx, row.UserID)
	if err != nil {
		return nil, err
	}
	if user.DisabledAt != nil {
		return nil, repository.ErrNotFound
	}
	_, err = s.store.UpdateAuthSessionTouch(ctx, tokenHash, request.IP, now)
	if err != nil {
		return nil, err
	}
	return &AuthSessionUser{ID: user.ID, Username: user.Username, ExpiresAt: row.ExpiresAt}, nil
}
