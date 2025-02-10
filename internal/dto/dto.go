package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateUserRequest represents the data required to create a new user.
type CreateUserRequest struct {
	// Name is the user's name. It must consist of alphanumeric characters, hyphens, or underscores.
	Name string `json:"name" validate:"required,regexp=^[a-zA-Z0-9\\-_]+$"`
	// Email is the user's email address. It must be a valid email format.
	Email string `json:"email" validate:"required,email"`
	// Password is the user's password. It must meet complexity requirements.
	Password string `json:"password" validate:"required,min=8,max=64,passwordComplexity"`
}

// CreateUserWithRoleRequest represents the data required for an admin to create a new user.
// Unlike CreateUserRequest, it allows setting the Role of the user.
type CreateUserWithRoleRequest struct {
	// Name is the user's name. It must consist of alphanumeric characters, hyphens, or underscores.
	Name string `json:"name" validate:"required,regexp=^[a-zA-Z0-9\\-_]+$"`
	// Email is the user's email address. It must be a valid email format.
	Email string `json:"email" validate:"required,email"`
	// Password is the user's password. It must meet complexity requirements.
	Password string `json:"password" validate:"required,min=8,max=64,passwordComplexity"`
	// Role defines the user's role in the system (e.g., "admin" or "user").
	Role string `json:"role" validate:"required,oneof=user admin"`
}

type CreateUserResponse struct {
	ID uuid.UUID `json:"id"`
}

// UpdateUserPasswordRequest represents the data required to update a user's password.
type UpdateUserPasswordRequest struct {
	// NewPassword is the user's new password. It must meet complexity requirements.
	NewPassword string `json:"new_password" validate:"required,min=8,max=64,passwordComplexity"`
}

// UpdateUserNameRequest represents the data required to update a user's name.
type UpdateUserNameRequest struct {
	// Name is the user's new name. It must consist of alphanumeric characters, hyphens, or underscores.
	Name string `json:"name" validate:"required,regexp=^[a-zA-Z0-9\\-_]+$"`
}

// UpdateUserEmailRequest represents the data required to update a user's email.
type UpdateUserEmailRequest struct {
	// Email is the user's new email address. It must be a valid email format.
	Email string `json:"email" validate:"required,email"`
}

// ReadUsersResponse represents the response containing a list of users.
// This DTO is used by admins to retrieve user data without sensitive information like passwords.
type ReadUsersResponse struct {
	// Users is a list of user details.
	Users []ReadUserResponse `json:"users"`
}

// ReadUserResponse represents the details of a single user.
// This DTO excludes sensitive information like passwords.
type ReadUserResponse struct {
	// ID is the unique identifier of the user.
	ID uuid.UUID `json:"id"`
	// Name is the user's name.
	Name string `json:"name"`
	// Email is the user's email address.
	Email string `json:"email"`
	// Role defines the user's role in the system (e.g., "admin" or "user").
	Role string `json:"role"`
	// CreatedAt is the timestamp when the user was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the timestamp when the user was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateUserDataRequest represents the data required for an admin to update all user details except the UserID.
type UpdateUserDataRequest struct {
	// Name is the user's new name. It must consist of alphanumeric characters, hyphens, or underscores.
	Name string `json:"name" validate:"required,regexp=^[a-zA-Z0-9\\-_]+$"`
	// Email is the user's new email address. It must be a valid email format.
	Email string `json:"email" validate:"required,email"`
	// Role defines the user's new role in the system (e.g., "admin" or "user").
	Role string `json:"role" validate:"required,oneof=user admin"`
}

// UpdateUserRoleRequest represents the data required for an admin to update a user's role.
type UpdateUserRoleRequest struct {
	// Role defines the user's new role in the system (e.g., "admin" or "user").
	Role string `json:"role" validate:"required,oneof=user admin"`
}
