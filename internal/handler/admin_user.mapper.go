package handler

import "github.com/lwmacct/260630-go-hsr-auth/internal/service"

func ToAdminUserListDTO(output *service.AdminUserListOutput) AdminUserListDTO {
	if output == nil {
		return AdminUserListDTO{Items: []AdminUserDTO{}}
	}
	body := AdminUserListDTO{
		Items:    []AdminUserDTO{},
		Total:    output.Total,
		Page:     output.Page,
		PageSize: output.PageSize,
	}
	for _, item := range output.Items {
		body.Items = append(body.Items, ToAdminUserDTO(item))
	}
	return body
}

func ToAdminUserDTO(user service.AdminUser) AdminUserDTO {
	return AdminUserDTO{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
		Status:      user.Status,
		Admin:       user.Admin,
		DisabledAt:  user.DisabledAt,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
