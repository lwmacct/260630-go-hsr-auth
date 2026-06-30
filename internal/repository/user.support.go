package repository

import (
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type UserCreate struct {
	Username    string
	DisplayName string
	Email       string
	AvatarURL   string
	Role        string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserRow struct {
	ID          int64
	Username    string
	DisplayName string
	Email       string
	AvatarURL   string
	Role        string
	Status      string
	DisabledAt  *time.Time
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserFilter struct {
	Keyword  string
	Role     string
	Status   string
	Page     int
	PageSize int
}

type UserProfilePatch struct {
	DisplayName string
	Email       string
	AvatarURL   string
	UpdatedAt   time.Time
}

type UserTotal struct {
	Count int
}

type UserChange struct {
	Affected int64
}

func utilUserRow(row UserModel) UserRow {
	return UserRow{
		ID:          row.ID,
		Username:    row.Username,
		DisplayName: row.DisplayName,
		Email:       row.Email,
		AvatarURL:   row.AvatarURL,
		Role:        row.Role,
		Status:      row.Status,
		DisabledAt:  row.DisabledAt,
		LastLoginAt: row.LastLoginAt,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func utilNullableString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func utilUserRows(rows []UserModel) []UserRow {
	items := make([]UserRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, utilUserRow(row))
	}
	return items
}

func utilApplyUserFilter(query *bun.SelectQuery, filter UserFilter) *bun.SelectQuery {
	keyword := strings.TrimSpace(filter.Keyword)
	if keyword != "" {
		like := "%" + strings.ToLower(keyword) + "%"
		query = query.Where("(LOWER(username) LIKE ? OR LOWER(display_name) LIKE ? OR LOWER(email) LIKE ?)", like, like, like)
	}
	if role := strings.TrimSpace(filter.Role); role != "" {
		query = query.Where("role = ?", role)
	}
	if status := strings.TrimSpace(filter.Status); status != "" {
		query = query.Where("status = ?", status)
	}
	return query
}
