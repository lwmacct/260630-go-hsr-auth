package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

func (s *Store) CreateAuthPassword(ctx context.Context, item AuthPasswordCreate) (*AuthPasswordRow, error) {
	row := AuthPasswordModel{
		UserID:            item.UserID,
		PasswordHash:      item.PasswordHash,
		PasswordChangedAt: item.PasswordChangedAt,
		CreatedAt:         item.CreatedAt,
		UpdatedAt:         item.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(&row).Exec(ctx); err != nil {
		return nil, err
	}
	result := utilAuthPasswordRow(row)
	return &result, nil
}

func (s *Store) UpsertAuthPassword(ctx context.Context, item AuthPasswordCreate) (*AuthPasswordRow, error) {
	row := AuthPasswordModel{
		UserID:            item.UserID,
		PasswordHash:      item.PasswordHash,
		PasswordChangedAt: item.PasswordChangedAt,
		CreatedAt:         item.CreatedAt,
		UpdatedAt:         item.UpdatedAt,
	}
	if _, err := s.db.NewInsert().
		Model(&row).
		On("CONFLICT (user_id) DO UPDATE").
		Set("password_hash = EXCLUDED.password_hash").
		Set("password_changed_at = EXCLUDED.password_changed_at").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx); err != nil {
		return nil, err
	}
	result := utilAuthPasswordRow(row)
	return &result, nil
}

func (s *Store) FetchAuthPassword(ctx context.Context, userID string) (*AuthPasswordRow, error) {
	row := new(AuthPasswordModel)
	if err := s.db.NewSelect().Model(row).Where("user_id = ?", userID).Scan(ctx); err != nil {
		return nil, WrapNotFound(err)
	}
	result := utilAuthPasswordRow(*row)
	return &result, nil
}

func (s *Store) UpdateAuthPassword(ctx context.Context, userID string, passwordHash string, updatedAt time.Time) (*AuthPasswordChange, error) {
	result, err := s.db.NewUpdate().
		Model((*AuthPasswordModel)(nil)).
		Set("password_hash = ?", passwordHash).
		Set("password_changed_at = ?", updatedAt).
		Set("updated_at = ?", updatedAt).
		Where("user_id = ?", userID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err == nil && affected == 0 {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &AuthPasswordChange{Affected: affected}, nil
}

func (s *Store) DeleteAuthPasswordForUsers(ctx context.Context, userIDs []string) error {
	_, err := s.db.NewDelete().
		Model((*AuthPasswordModel)(nil)).
		Where("user_id IN (?)", bun.List(userIDs)).
		Exec(ctx)
	return err
}
