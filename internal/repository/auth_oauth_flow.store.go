package repository

import (
	"context"
	"time"
)

func (s *Store) CreateAuthOauthFlow(ctx context.Context, item AuthOauthFlowCreate) (*AuthOauthFlowRow, error) {
	row := AuthOauthFlowModel{
		StateHash:        item.StateHash,
		Provider:         item.Provider,
		PKCECodeVerifier: item.PKCECodeVerifier,
		Nonce:            item.Nonce,
		ReturnTo:         item.ReturnTo,
		ExpiresAt:        item.ExpiresAt,
		CreatedAt:        item.CreatedAt,
	}
	if _, err := s.db.NewInsert().Model(&row).Exec(ctx); err != nil {
		return nil, err
	}
	result := utilAuthOauthFlowRow(row)
	return &result, nil
}

func (s *Store) FetchAuthOauthFlow(ctx context.Context, stateHash []byte) (*AuthOauthFlowRow, error) {
	row := new(AuthOauthFlowModel)
	if err := s.db.NewSelect().Model(row).Where("state_hash = ?", stateHash).Scan(ctx); err != nil {
		return nil, WrapNotFound(err)
	}
	result := utilAuthOauthFlowRow(*row)
	return &result, nil
}

func (s *Store) DeleteAuthOauthFlow(ctx context.Context, stateHash []byte) error {
	_, err := s.db.NewDelete().
		Model((*AuthOauthFlowModel)(nil)).
		Where("state_hash = ?", stateHash).
		Exec(ctx)
	return err
}

func (s *Store) DeleteAuthOauthFlowExpired(ctx context.Context, expiresBefore time.Time) error {
	_, err := s.db.NewDelete().
		Model((*AuthOauthFlowModel)(nil)).
		Where("expires_at < ?", expiresBefore).
		Exec(ctx)
	return err
}
