package handler

import "time"

type AdminUserListInputDTO struct {
	Session  string `cookie:"web_session"`
	Keyword  string `query:"keyword"`
	Role     string `query:"role"`
	Status   string `query:"status"`
	Page     int    `query:"page" default:"1"`
	PageSize int    `query:"pageSize" default:"20"`
}

type AdminCreateUserInputDTO struct {
	Session string `cookie:"web_session"`
	Body    AdminCreateUserDTO
}

type AdminUpdateUserInputDTO struct {
	Session string `cookie:"web_session"`
	ID      string `path:"id"`
	Body    AdminUpdateUserDTO
}

type AdminBatchRoleInputDTO struct {
	Session string `cookie:"web_session"`
	Body    AdminBatchRoleDTO
}

type AdminBatchStatusInputDTO struct {
	Session string `cookie:"web_session"`
	Body    AdminBatchStatusDTO
}

type AdminBatchPasswordInputDTO struct {
	Session string `cookie:"web_session"`
	Body    AdminBatchPasswordDTO
}

type AdminBatchDeleteInputDTO struct {
	Session string `cookie:"web_session"`
	Body    AdminBatchIDsDTO
}

type AdminUserDTO struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	DisplayName string     `json:"displayName"`
	Email       string     `json:"email,omitempty"`
	AvatarURL   string     `json:"avatarUrl,omitempty"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	Admin       bool       `json:"admin"`
	DisabledAt  *time.Time `json:"disabledAt,omitempty"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type AdminUserListDTO struct {
	Items    []AdminUserDTO `json:"items"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

type AdminCreateUserDTO struct {
	Username    string  `json:"username"`
	DisplayName *string `json:"displayName,omitempty"`
	Email       *string `json:"email,omitempty"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
	Role        *string `json:"role,omitempty"`
	Password    *string `json:"password,omitempty"`
}

type AdminUpdateUserDTO struct {
	DisplayName string  `json:"displayName"`
	Email       *string `json:"email,omitempty"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
}

type AdminBatchRoleDTO struct {
	IDs  []string `json:"ids"`
	Role string   `json:"role"`
}

type AdminBatchStatusDTO struct {
	IDs    []string `json:"ids"`
	Status string   `json:"status"`
}

type AdminBatchPasswordDTO struct {
	IDs      []string `json:"ids"`
	Password string   `json:"password"`
}

type AdminBatchIDsDTO struct {
	IDs []string `json:"ids"`
}
