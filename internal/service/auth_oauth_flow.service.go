package service

import (
	"context"
	"errors"
	"time"

	"github.com/lwmacct/260630-go-hsr-auth/internal/infra/token"
	"github.com/lwmacct/260630-go-hsr-auth/internal/repository"
)

type AuthOauthFlowService struct {
	store *repository.Store
	ttl   time.Duration
}

func NewAuthOauthFlowService(store *repository.Store, ttl time.Duration) *AuthOauthFlowService {
	if store == nil {
		panic("NewAuthOauthFlowService: store is nil")
	}
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &AuthOauthFlowService{store: store, ttl: ttl}
}

func (s *AuthOauthFlowService) Create(ctx context.Context, provider string, returnTo string) (state string, flow AuthOauthFlow, err error) {
	state, err = token.NewWithPrefix("oauth_state")
	if err != nil {
		return "", AuthOauthFlow{}, err
	}
	codeVerifier, err := token.NewBase62(64)
	if err != nil {
		return "", AuthOauthFlow{}, err
	}
	nonce, err := token.NewWithPrefix("nonce")
	if err != nil {
		return "", AuthOauthFlow{}, err
	}
	now := time.Now().UTC()
	row, err := s.store.CreateAuthOauthFlow(ctx, repository.AuthOauthFlowCreate{
		StateHash:        utilAuthOauthFlowStateHash(state),
		Provider:         provider,
		PKCECodeVerifier: codeVerifier,
		Nonce:            nonce,
		ReturnTo:         returnTo,
		ExpiresAt:        now.Add(s.ttl),
		CreatedAt:        now,
	})
	if err != nil {
		return "", AuthOauthFlow{}, err
	}
	return state, utilAuthOauthFlow(*row), nil
}

func (s *AuthOauthFlowService) Consume(ctx context.Context, state string) (*AuthOauthFlow, error) {
	stateHash := utilAuthOauthFlowStateHash(state)
	row, err := s.store.FetchAuthOauthFlow(ctx, stateHash)
	if err != nil {
		return nil, err
	}
	_ = s.store.DeleteAuthOauthFlow(ctx, stateHash)
	if time.Now().UTC().After(row.ExpiresAt) {
		return nil, repository.ErrNotFound
	}
	flow := utilAuthOauthFlow(*row)
	return &flow, nil
}

func (s *AuthOauthFlowService) Cleanup(ctx context.Context) error {
	err := s.store.DeleteAuthOauthFlowExpired(ctx, time.Now().UTC())
	if errors.Is(err, repository.ErrNotFound) {
		return nil
	}
	return err
}
