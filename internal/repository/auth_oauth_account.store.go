package repository

import (
	"context"

	"github.com/uptrace/bun"
)

func (s *Store) FetchAuthOauthAccountByProviderSubject(ctx context.Context, provider string, subject string) (*AuthOauthAccountRow, error) {
	row := new(AuthOauthAccountModel)
	if err := s.db.NewSelect().Model(row).
		Where("provider = ?", provider).
		Where("subject = ?", subject).
		Scan(ctx); err != nil {
		return nil, WrapNotFound(err)
	}
	result := utilAuthOauthAccountRow(*row)
	return &result, nil
}

func (s *Store) CreateAuthOauthAccount(ctx context.Context, item AuthOauthAccountCreate) (*AuthOauthAccountRow, error) {
	row := AuthOauthAccountModel{
		UserID:                item.UserID,
		Provider:              item.Provider,
		Subject:               item.Subject,
		ProviderEmail:         item.ProviderEmail,
		ProviderEmailVerified: item.ProviderEmailVerified,
		ProviderDisplayName:   item.ProviderDisplayName,
		ProviderAvatarURL:     item.ProviderAvatarURL,
		ProviderProfile:       item.ProviderProfile,
		CreatedAt:             item.CreatedAt,
		UpdatedAt:             item.UpdatedAt,
	}
	if _, err := s.db.NewInsert().Model(&row).Exec(ctx); err != nil {
		return nil, err
	}
	result := utilAuthOauthAccountRow(row)
	return &result, nil
}

func (s *Store) DeleteAuthOauthAccountForUsers(ctx context.Context, userIDs []int64) error {
	_, err := s.db.NewDelete().
		Model((*AuthOauthAccountModel)(nil)).
		Where("user_id IN (?)", bun.List(userIDs)).
		Exec(ctx)
	return err
}
