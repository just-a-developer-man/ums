package v1

import (
	"context"
	"fmt"
	"ums/internal/dto"

	"github.com/google/uuid"
)

// UMS defines the interface for user management operations.
type UMS interface {
	CreateUser(ctx context.Context, createUser dto.CreateUserRequest) (uuid.UUID, error)
	CreateUserWithRole(ctx context.Context, createUser dto.CreateUserWithRoleRequest) (uuid.UUID, error)
	ReadUsers(ctx context.Context) (dto.ReadUsersResponse, error)
	ReadUser(ctx context.Context, uid uuid.UUID) (dto.ReadUserResponse, error)
	UpdateUserData(ctx context.Context, uid uuid.UUID, updateUserData dto.UpdateUserDataRequest) error
	UpdateUserRole(ctx context.Context, uid uuid.UUID, updateUserRole dto.UpdateUserRoleRequest) error
	UpdateUserPassword(ctx context.Context, uid uuid.UUID, updatePassword dto.UpdateUserPasswordRequest) error
	UpdateUserName(ctx context.Context, uid uuid.UUID, updateUserName dto.UpdateUserNameRequest) error
	UpdateUserEmail(ctx context.Context, uid uuid.UUID, updateUserEmail dto.UpdateUserEmailRequest) error
	DeleteUser(ctx context.Context, uid uuid.UUID) error
}

// Validator defines the interface for validating tagged structs and UID strings.
type Validator interface {
	ValidateUID(string) error
	ValidateStruct(interface{}) error
}

// Handler is an HTTP handler responsible for managing user-related requests.
type Handler struct {
	service   UMS
	validator Validator
}

// New creates a new instance of the Handler.
func New(service UMS, validator Validator) (*Handler, error) {
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}
	if validator == nil {
		return nil, fmt.Errorf("validator cannot be nil")
	}
	return &Handler{
		service:   service,
		validator: validator,
	}, nil
}
