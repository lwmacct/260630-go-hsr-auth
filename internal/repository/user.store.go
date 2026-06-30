package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

func (s *Store) CreateUser(ctx context.Context, item UserCreate) (*UserRow, error) {
	row := UserModel{
		Username:    item.Username,
		DisplayName: item.DisplayName,
		Email:       item.Email,
		AvatarURL:   item.AvatarURL,
		Role:        item.Role,
		Status:      item.Status,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(&row).Exec(ctx); err != nil {
		return nil, err
	}
	result := utilUserRow(row)
	return &result, nil
}

func (s *Store) FetchUser(ctx context.Context, id int64) (*UserRow, error) {
	row := new(UserModel)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, WrapNotFound(err)
	}
	result := utilUserRow(*row)
	return &result, nil
}

func (s *Store) FetchUserByUsername(ctx context.Context, username string) (*UserRow, error) {
	row := new(UserModel)
	if err := s.db.NewSelect().Model(row).Where("username = ?", username).Scan(ctx); err != nil {
		return nil, WrapNotFound(err)
	}
	result := utilUserRow(*row)
	return &result, nil
}

func (s *Store) ListUsers(ctx context.Context) ([]UserRow, error) {
	var rows []UserModel
	if err := s.db.NewSelect().Model(&rows).Order("username ASC").Scan(ctx); err != nil {
		return nil, err
	}
	return utilUserRows(rows), nil
}

func (s *Store) ListUsersByFilter(ctx context.Context, filter UserFilter) ([]UserRow, error) {
	var rows []UserModel
	query := utilApplyUserFilter(s.db.NewSelect().Model(&rows), filter)
	if err := query.
		Order("id DESC").
		Limit(filter.PageSize).
		Offset((filter.Page - 1) * filter.PageSize).
		Scan(ctx); err != nil {
		return nil, err
	}
	return utilUserRows(rows), nil
}

func (s *Store) FetchUserTotalByFilter(ctx context.Context, filter UserFilter) (*UserTotal, error) {
	query := utilApplyUserFilter(s.db.NewSelect().Model((*UserModel)(nil)), filter)
	count, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}
	return &UserTotal{Count: count}, nil
}

func (s *Store) UpdateUserProfile(ctx context.Context, id int64, patch UserProfilePatch) (*UserRow, error) {
	result, err := s.db.NewUpdate().
		Model((*UserModel)(nil)).
		Set("display_name = ?", patch.DisplayName).
		Set("email = ?", utilNullableString(patch.Email)).
		Set("avatar_url = ?", utilNullableString(patch.AvatarURL)).
		Set("updated_at = ?", patch.UpdatedAt).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return nil, ErrNotFound
	}
	return s.FetchUser(ctx, id)
}

func (s *Store) UpdateUserLastLogin(ctx context.Context, id int64, lastLoginAt time.Time) (*UserChange, error) {
	result, err := s.db.NewUpdate().
		Model((*UserModel)(nil)).
		Set("last_login_at = ?", lastLoginAt).
		Set("updated_at = ?", lastLoginAt).
		Where("id = ?", id).
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
	return &UserChange{Affected: affected}, nil
}

func (s *Store) UpdateUserRoleBatch(ctx context.Context, ids []int64, role string, updatedAt time.Time) (*UserChange, error) {
	result, err := s.db.NewUpdate().
		Model((*UserModel)(nil)).
		Set("role = ?", role).
		Set("updated_at = ?", updatedAt).
		Where("id IN (?)", bun.List(ids)).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &UserChange{Affected: affected}, nil
}

func (s *Store) UpdateUserStatusBatch(ctx context.Context, ids []int64, status string, updatedAt time.Time, disabledAt *time.Time) (*UserChange, error) {
	query := s.db.NewUpdate().
		Model((*UserModel)(nil)).
		Set("status = ?", status).
		Set("updated_at = ?", updatedAt).
		Where("id IN (?)", bun.List(ids))
	if disabledAt == nil {
		query = query.Set("disabled_at = NULL")
	} else {
		query = query.Set("disabled_at = ?", *disabledAt)
	}
	result, err := query.Exec(ctx)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &UserChange{Affected: affected}, nil
}

func (s *Store) DeleteUserBatch(ctx context.Context, ids []int64) error {
	_, err := s.db.NewDelete().
		Model((*UserModel)(nil)).
		Where("id IN (?)", bun.List(ids)).
		Exec(ctx)
	return err
}
