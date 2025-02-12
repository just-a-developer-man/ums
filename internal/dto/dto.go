package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateUserRequest represents the data required to create a new user.
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,regexp=^[a-zA-Z0-9\\-_]+$"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64,passwordComplexity"`
}

// CreateUserWithRoleRequest represents the data required for an admin to create a new user.
type CreateUserWithRoleRequest struct {
	Name     string `json:"name" validate:"required,regexp=^[a-zA-Z0-9\\-_]+$"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64,passwordComplexity"`
	Role     string `json:"role" validate:"required,oneof=user admin"`
}

// CreateUserResponse represents the answere on CreateUser request.
type CreateUserResponse struct {
	ID uuid.UUID `json:"id"`
}

// UpdateUserPasswordRequest represents the data required to update a user's password.
type UpdateUserPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8,max=64,passwordComplexity"`
}

// UpdateUserNameRequest represents the data required to update a user's name.
type UpdateUserNameRequest struct {
	Name string `json:"name" validate:"required,regexp=^[a-zA-Z0-9\\-_]+$"`
}

// UpdateUserEmailRequest represents the data required to update a user's email.
type UpdateUserEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ReadUsersResponse represents the response containing a list of users.
type ReadUsersResponse struct {
	Users []ReadUserResponse `json:"users"`
}

// ReadUserResponse represents the details of a single user.
type ReadUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateUserDataRequest represents the data required for an admin to update all user details.
type UpdateUserDataRequest struct {
	Name  string `json:"name" validate:"required,regexp=^[a-zA-Z0-9\\-_]+$"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=user admin"`
}

// UpdateUserRoleRequest represents the data required for an admin to update a user's role.
type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=user admin"`
}
