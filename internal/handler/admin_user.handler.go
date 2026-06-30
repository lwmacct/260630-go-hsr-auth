package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
)

type adminUserHandler struct {
	auth authHandler
}

func RegisterAdminUser(api huma.API, config Config, services Services) {
	handler := adminUserHandler{auth: authHandler{config: config, services: services}}
	admin := huma.NewGroup(api, "/admin")
	huma.Register(admin, huma.Operation{OperationID: "admin-list-users", Method: http.MethodGet, Path: "/users", Tags: []string{"Admin"}}, handler.listUsers)
	huma.Register(admin, huma.Operation{OperationID: "admin-create-user", Method: http.MethodPost, Path: "/users", DefaultStatus: http.StatusCreated, Tags: []string{"Admin"}}, handler.createUser)
	huma.Register(admin, huma.Operation{OperationID: "admin-update-user", Method: http.MethodPatch, Path: "/users/{id}", Tags: []string{"Admin"}}, handler.updateUser)
	huma.Register(admin, huma.Operation{OperationID: "admin-batch-user-role", Method: http.MethodPost, Path: "/users/batch-role", Tags: []string{"Admin"}}, handler.batchUserRole)
	huma.Register(admin, huma.Operation{OperationID: "admin-batch-user-status", Method: http.MethodPost, Path: "/users/batch-status", Tags: []string{"Admin"}}, handler.batchUserStatus)
	huma.Register(admin, huma.Operation{OperationID: "admin-batch-user-password", Method: http.MethodPost, Path: "/users/batch-password", Tags: []string{"Admin"}}, handler.batchUserPassword)
	huma.Register(admin, huma.Operation{OperationID: "admin-delete-users", Method: http.MethodDelete, Path: "/users", Tags: []string{"Admin"}}, handler.deleteUsers)
}

func (h adminUserHandler) listUsers(ctx context.Context, input *AdminUserListInputDTO) (*BodyDTO[AdminUserListDTO], error) {
	if _, err := h.requireAdmin(ctx, input.Session); err != nil {
		return nil, err
	}
	output, err := h.auth.services.AdminUsers.List(ctx, service.AdminUserListInput{
		Keyword:  input.Keyword,
		Role:     input.Role,
		Status:   input.Status,
		Page:     input.Page,
		PageSize: input.PageSize,
	})
	if err != nil {
		return nil, utilAdminAPIError(err)
	}
	return &BodyDTO[AdminUserListDTO]{Body: ToAdminUserListDTO(output)}, nil
}

func (h adminUserHandler) createUser(ctx context.Context, input *AdminCreateUserInputDTO) (*BodyDTO[AdminUserDTO], error) {
	if _, err := h.requireAdmin(ctx, input.Session); err != nil {
		return nil, err
	}
	user, err := h.auth.services.AdminUsers.Create(ctx, service.AdminUserCreateInput{
		Username:    input.Body.Username,
		DisplayName: utilStringValue(input.Body.DisplayName),
		Email:       utilStringValue(input.Body.Email),
		AvatarURL:   utilStringValue(input.Body.AvatarURL),
		Role:        utilStringValue(input.Body.Role),
		Password:    utilStringValue(input.Body.Password),
	})
	if err != nil {
		return nil, utilAdminAPIError(err)
	}
	return &BodyDTO[AdminUserDTO]{Body: ToAdminUserDTO(*user)}, nil
}

func (h adminUserHandler) updateUser(ctx context.Context, input *AdminUpdateUserInputDTO) (*BodyDTO[AdminUserDTO], error) {
	if _, err := h.requireAdmin(ctx, input.Session); err != nil {
		return nil, err
	}
	user, err := h.auth.services.AdminUsers.UpdateProfile(ctx, input.ID, service.AdminUserUpdateProfileInput{
		DisplayName: input.Body.DisplayName,
		Email:       utilStringValue(input.Body.Email),
		AvatarURL:   utilStringValue(input.Body.AvatarURL),
	})
	if err != nil {
		return nil, utilAdminAPIError(err)
	}
	return &BodyDTO[AdminUserDTO]{Body: ToAdminUserDTO(*user)}, nil
}

func (h adminUserHandler) batchUserRole(ctx context.Context, input *AdminBatchRoleInputDTO) (*BodyDTO[ActionDTO], error) {
	admin, err := h.requireAdmin(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.auth.services.AdminUsers.SetRoleBatch(ctx, admin.ID, input.Body.IDs, input.Body.Role); err != nil {
		return nil, utilAdminAPIError(err)
	}
	return &BodyDTO[ActionDTO]{Body: ActionDTO{OK: true}}, nil
}

func (h adminUserHandler) batchUserStatus(ctx context.Context, input *AdminBatchStatusInputDTO) (*BodyDTO[ActionDTO], error) {
	admin, err := h.requireAdmin(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.auth.services.AdminUsers.SetStatusBatch(ctx, admin.ID, input.Body.IDs, input.Body.Status); err != nil {
		return nil, utilAdminAPIError(err)
	}
	return &BodyDTO[ActionDTO]{Body: ActionDTO{OK: true}}, nil
}

func (h adminUserHandler) batchUserPassword(ctx context.Context, input *AdminBatchPasswordInputDTO) (*BodyDTO[ActionDTO], error) {
	admin, err := h.requireAdmin(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.auth.services.AdminUsers.ResetPasswordBatch(ctx, admin.ID, input.Body.IDs, input.Body.Password); err != nil {
		return nil, utilAdminAPIError(err)
	}
	return &BodyDTO[ActionDTO]{Body: ActionDTO{OK: true}}, nil
}

func (h adminUserHandler) deleteUsers(ctx context.Context, input *AdminBatchDeleteInputDTO) (*BodyDTO[ActionDTO], error) {
	admin, err := h.requireAdmin(ctx, input.Session)
	if err != nil {
		return nil, err
	}
	if err := h.auth.services.AdminUsers.DeleteBatch(ctx, admin.ID, input.Body.IDs); err != nil {
		return nil, utilAdminAPIError(err)
	}
	return &BodyDTO[ActionDTO]{Body: ActionDTO{OK: true}}, nil
}

func (h adminUserHandler) requireAdmin(ctx context.Context, sessionID string) (*AuthUserDTO, error) {
	user, err := h.auth.currentUser(ctx, sessionID)
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	if !user.Admin {
		return nil, huma.Error403Forbidden("forbidden")
	}
	return user, nil
}
