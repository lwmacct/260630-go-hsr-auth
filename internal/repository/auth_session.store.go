package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

func (s *Store) CreateAuthSession(ctx context.Context, item AuthSessionCreate) (*AuthSessionRow, error) {
	row := AuthSessionModel{
		TokenHash:     item.TokenHash,
		UserID:        item.UserID,
		LoginIP:       item.LoginIP,
		LastIP:        item.LastIP,
		UserAgentHash: item.UserAgentHash,
		ExpiresAt:     item.ExpiresAt,
		CreatedAt:     item.CreatedAt,
		LastSeenAt:    item.LastSeenAt,
	}
	if _, err := s.db.NewInsert().Model(&row).Exec(ctx); err != nil {
		return nil, err
	}
	result := utilAuthSessionRow(row)
	return &result, nil
}

func (s *Store) FetchAuthSession(ctx context.Context, tokenHash []byte) (*AuthSessionRow, error) {
	row := new(AuthSessionModel)
	if err := s.db.NewSelect().
		Model(row).
		Where("token_hash = ?", tokenHash).
		Where("revoked_at IS NULL").
		Scan(ctx); err != nil {
		return nil, WrapNotFound(err)
	}
	result := utilAuthSessionRow(*row)
	return &result, nil
}

func (s *Store) UpdateAuthSessionTouch(ctx context.Context, tokenHash []byte, lastIP string, lastSeenAt time.Time) (*AuthSessionChange, error) {
	result, err := s.db.NewUpdate().
		Model((*AuthSessionModel)(nil)).
		Set("last_seen_at = ?", lastSeenAt).
		Set("last_ip = ?", lastIP).
		Where("token_hash = ?", tokenHash).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &AuthSessionChange{Affected: affected}, nil
}

func (s *Store) UpdateAuthSessionRevokedForUsers(ctx context.Context, userIDs []string, revokedAt time.Time) (*AuthSessionChange, error) {
	result, err := s.db.NewUpdate().
		Model((*AuthSessionModel)(nil)).
		Set("revoked_at = ?", revokedAt).
		Where("user_id IN (?)", bun.List(userIDs)).
		Where("revoked_at IS NULL").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &AuthSessionChange{Affected: affected}, nil
}

func (s *Store) DeleteAuthSession(ctx context.Context, tokenHash []byte) error {
	_, err := s.db.NewDelete().
		Model((*AuthSessionModel)(nil)).
		Where("token_hash = ?", tokenHash).
		Exec(ctx)
	return err
}

func (s *Store) DeleteAuthSessionForUsers(ctx context.Context, userIDs []string) error {
	_, err := s.db.NewDelete().
		Model((*AuthSessionModel)(nil)).
		Where("user_id IN (?)", bun.List(userIDs)).
		Exec(ctx)
	return err
}
