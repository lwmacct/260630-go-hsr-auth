package handler

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
)

func utilAdminAPIError(err error) error {
	switch {
	case errors.Is(err, service.ErrAdminUserNotFound):
		return huma.Error404NotFound("user not found")
	case errors.Is(err, service.ErrAdminUserCannotOperateSelf):
		return huma.Error403Forbidden("cannot operate current user")
	case errors.Is(err, service.ErrAdminUserEmptySelection):
		return huma.Error400BadRequest("empty user selection")
	case errors.Is(err, service.ErrUserInvalidCredentials):
		return huma.Error400BadRequest("invalid user")
	case errors.Is(err, service.ErrUserInvalidRole):
		return huma.Error400BadRequest("invalid role")
	case errors.Is(err, service.ErrUserInvalidStatus):
		return huma.Error400BadRequest("invalid status")
	case errors.Is(err, service.ErrUserUsernameTaken):
		return huma.Error400BadRequest("username taken")
	case errors.Is(err, service.ErrAuthPasswordWeakPassword):
		return huma.Error400BadRequest("weak password")
	default:
		return huma.Error500InternalServerError("internal server error")
	}
}

func utilStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
