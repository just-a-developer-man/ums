// Package handler_http provides HTTP handlers for the User Management System (UMS).
// It implements RESTful endpoints for creating, reading, updating, and deleting users.
// The package uses the Gin framework for routing and request handling.
// Each handler method is designed to handle specific HTTP requests and interact with the UMS service layer.
// Error handling is consistent across all handlers, providing meaningful feedback to clients in case of failures.
package handler_http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"ums/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// errMessage represents error messages returned by the HTTP handlers.
// These messages provide meaningful feedback to clients in case of failures.
type errMessage string

var (
	// invalidJSON indicates that the request body contains invalid JSON format.
	invalidJSON errMessage = "invalid JSON format"

	// createUserFailed indicates that the user creation process failed.
	createUserFailed errMessage = "failed to create user"

	// readUsersFailed indicates that the process of reading users failed.
	readUsersFailed errMessage = "failed to read users"

	// readUserFailed indicates that the process of reading a single user failed.
	readUserFailed errMessage = "failed to read user"

	// updateUserPasswordFailed indicates that the process of updating a user's password failed.
	updateUserPasswordFailed errMessage = "failed to update the user's password"

	// updateUserNameFailed indicates that the process of updating a user's name failed.
	updateUserNameFailed errMessage = "failed to update the user's name"

	// updateUserEmailFailed indicates that the process of updating a user's email failed.
	updateUserEmailFailed errMessage = "failed to update the user's email"

	// updateUserDataFailed indicates that the process of updating a user's data failed.
	updateUserDataFailed errMessage = "failed to update the user's data"

	// deleteUserFailed indicates that the process of deleting a user failed.Context
	deleteUserFailed errMessage = "failed to delete user"

	// uidParam is the URL parameter name used to identify users in requests.
	uidParam string = "id"
)

// UMS defines the interface for user management operations.
// It abstracts the business logic for creating, reading, updating, and deleting users.
// This interface ensures loose coupling between the HTTP handlers and the underlying service implementation.
type UMS interface {
	// CreateUser creates a new user in the system.
	// It returns the unique ID of the created user or an error if the operation fails.
	CreateUser(ctx context.Context, createUser dto.CreateUserRequest) (uuid.UUID, error)

	// CreateUserWithRole creates a new user in the system, allowing the setting of the user's role.
	// It returns the unique ID of the created user or an error if the operation fails.
	CreateUserWithRole(ctx context.Context, createUser dto.CreateUserWithRoleRequest) (uuid.UUID, error)

	// UpdateUserPassword updates a user's password in the system.
	// It validates the new password and ensures the user exists before performing the update.
	// Returns an error if the user does not exist, the new password is invalid, or the update fails.
	UpdateUserPassword(ctx context.Context, uid uuid.UUID, updatePassword dto.UpdateUserPasswordRequest) error

	// UpdateUserName updates a user's name in the system.
	// It validates the new name and ensures the user exists before performing the update.
	// Returns an error if the user does not exist, the new name is invalid, or the update fails.
	UpdateUserName(ctx context.Context, uid uuid.UUID, updateUserName dto.UpdateUserNameRequest) error

	// UpdateUserEmail updates a user's email in the system.
	// It validates the new email and ensures the user exists before performing the update.
	// Returns an error if the user does not exist, the new email is invalid, or the update fails.
	UpdateUserEmail(ctx context.Context, uid uuid.UUID, updateUserEmail dto.UpdateUserEmailRequest) error

	// ReadUsers retrieves a list of all users from the system, excluding sensitive information like passwords.
	// Returns an error if the read operation fails.
	ReadUsers(ctx context.Context) (dto.ReadUsersResponse, error)

	// ReadUser retrieves detailed information about a single user, excluding sensitive information like passwords.
	// Returns an error if the read operation fails or the user does not exist.
	ReadUser(ctx context.Context, uid uuid.UUID) (dto.ReadUserResponse, error)

	// UpdateUserData updates all user details except the UserID.
	// It validates the provided data and ensures the user exists before performing the update.
	// Returns an error if the user does not exist, the data is invalid, or the update fails.
	UpdateUserData(ctx context.Context, uid uuid.UUID, updateUserData dto.UpdateUserDataRequest) error

	// UpdateUserRole updates a user's role in the system.
	// It validates the new role and ensures the user exists before performing the update.
	// Returns an error if the user does not exist, the role is invalid, or the update fails.
	UpdateUserRole(ctx context.Context, uid uuid.UUID, updateUserRole dto.UpdateUserRoleRequest) error

	// DeleteUser removes a user from the system.
	// Returns an error if the deletion fails or the user does not exist.
	DeleteUser(ctx context.Context, uid uuid.UUID) error
}

// Validator defines the interface for validating incoming HTTP requests.
// It ensures that the data provided in requests is properly formatted and valid.
type Validator interface {
	ValidateCreateUserReq(dto.CreateUserRequest) error
	ValidateCreateUserWithRoleReq(dto.CreateUserWithRoleRequest) error
	ValidateUID(string) error
	ValidateUpdateUserPasswordReq(dto.UpdateUserPasswordRequest) error
	ValidateUpdateUserNameReq(dto.UpdateUserNameRequest) error
	ValidateUpdateUserEmailReq(dto.UpdateUserEmailRequest) error
	ValidateUpdateUserDataReq(dto.UpdateUserDataRequest) error
	ValidateUpdateUserRoleReq(dto.UpdateUserRoleRequest) error
}

// Handler is an HTTP handler responsible for managing user-related requests.
// It interacts with the UMS service to perform CRUD operations on users.
// The handler also uses a Validator to ensure incoming requests are properly formatted and valid.
type Handler struct {
	// service is the interface that provides business logic for user management.
	service UMS

	// validator is the interface that provides methods for validating request data.
	validator Validator
}

// New creates a new instance of the Handler.
// It takes an implementation of the UMS interface and a Validator as dependencies.
// This allows for easy testing and decoupling of business logic and validation from HTTP handling.
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

// CreateUser handles the HTTP POST request to create a new user.
// It expects a JSON payload in the request body containing user details.
// If successful, it returns a 201 Created status with the ID of the newly created user.
// Possible errors:
//   - 400 Bad Request: Invalid JSON or validation failure.
//   - 500 Internal Server Error: Failed to create the user in the service layer.
func (h *Handler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	createUser := dto.CreateUserRequest{}
	err := c.ShouldBindJSON(&createUser)
	if err != nil {
		slog.Error("Cannot bind body request dto", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidJSON})
		return
	}
	err = h.validator.ValidateCreateUserReq(createUser)
	if err != nil {
		slog.Error("Failed to validate request body", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := h.service.CreateUser(ctx, createUser)
	if err != nil {
		slog.Error("Failed to create user", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": createUserFailed})
		return
	}
	slog.Info("User created")
	c.JSON(http.StatusCreated, gin.H{"id": uid.String()})
}

// CreateUserWithRole handles the HTTP POST request to create a new user with administrative privileges.
// Unlike CreateUser, it allows setting the Role of the user.
// It expects a JSON payload in the request body containing user details.
// If successful, it returns a 201 Created status with the ID of the newly created user.
// Possible errors:
//   - 400 Bad Request: Invalid JSON or validation failure.
//   - 500 Internal Server Error: Failed to create the user in the service layer.
func (h *Handler) CreateUserWithRole(c *gin.Context) {
	ctx := c.Request.Context()
	createUser := dto.CreateUserWithRoleRequest{}
	err := c.ShouldBindJSON(&createUser)
	if err != nil {
		slog.Error("Cannot bind body request to dto", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidJSON})
		return
	}
	err = h.validator.ValidateCreateUserWithRoleReq(createUser)
	if err != nil {
		slog.Error("Failed to validate request body", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := h.service.CreateUserWithRole(ctx, createUser)
	if err != nil {
		slog.Error("Failed to create user", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": createUserFailed})
		return
	}
	slog.Info("User created")
	c.JSON(http.StatusCreated, gin.H{"id": uid.String()})
}

// UpdateUserPassword handles the HTTP PUT request to update a user's password by their ID.
// The user ID is provided as a URL parameter.
// If successful, it returns a 204 No Content status.
// Possible errors:
//   - 400 Bad Request: Invalid ID parameter or invalid password format.
//   - 500 Internal Server Error: Failed to update the user's password in the service layer.
func (h *Handler) UpdateUserPassword(c *gin.Context) {
	uidFromParam := c.Param(uidParam)
	err := h.validator.ValidateUID(uidFromParam)
	if err != nil {
		slog.Error("Failed to validate id param", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := uuid.Parse(uidFromParam)
	if err != nil {
		slog.Error("Failed to build UUID from id param", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateUserPassword := dto.UpdateUserPasswordRequest{}
	err = c.ShouldBindBodyWithJSON(&updateUserPassword)
	if err != nil {
		slog.Error("Failed to bind request body to dto", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidJSON})
		return
	}
	err = h.validator.ValidateUpdateUserPasswordReq(updateUserPassword)
	if err != nil {
		slog.Error("Failed to validate request body", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	err = h.service.UpdateUserPassword(ctx, uid, updateUserPassword)
	if err != nil {
		slog.Error("Failed to update user's password", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateUserPasswordFailed})
		return
	}
	slog.Info("User's password updated")
	c.Status(http.StatusNoContent)
}

// UpdateUserName handles the HTTP PUT request to update a user's name by their ID.
// The user ID is provided as a URL parameter.
// If successful, it returns a 204 No Content status.
// Possible errors:
//   - 400 Bad Request: Invalid ID parameter or invalid name format.
//   - 500 Internal Server Error: Failed to update the user's name in the service layer.
func (h *Handler) UpdateUserName(c *gin.Context) {
	uidFromParam := c.Param(uidParam)
	err := h.validator.ValidateUID(uidFromParam)
	if err != nil {
		slog.Error("Failed to validate id param", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := uuid.Parse(uidFromParam)
	if err != nil {
		slog.Error("Failed to build UUID from id param", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateUserName := dto.UpdateUserNameRequest{}
	err = c.ShouldBindBodyWithJSON(&updateUserName)
	if err != nil {
		slog.Error("Failed to bind request body to dto", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidJSON})
		return
	}
	err = h.validator.ValidateUpdateUserNameReq(updateUserName)
	if err != nil {
		slog.Error("Failed to validate request body", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	err = h.service.UpdateUserName(ctx, uid, updateUserName)
	if err != nil {
		slog.Error("Failed to update user's name", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateUserNameFailed})
		return
	}
	slog.Info("User's name updated")
	c.Status(http.StatusNoContent)
}

// UpdateUserEmail handles the HTTP PUT request to update a user's email by their ID.
// The user ID is provided as a URL parameter.
// If successful, it returns a 204 No Content status.
// Possible errors:
//   - 400 Bad Request: Invalid ID parameter or invalid email format.
//   - 500 Internal Server Error: Failed to update the user's email in the service layer.
func (h *Handler) UpdateUserEmail(c *gin.Context) {
	uidFromParam := c.Param(uidParam)
	err := h.validator.ValidateUID(uidFromParam)
	if err != nil {
		slog.Error("Failed to validate id param", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := uuid.Parse(uidFromParam)
	if err != nil {
		slog.Error("Failed to build UUID from id param", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateUserEmail := dto.UpdateUserEmailRequest{}
	err = c.ShouldBindBodyWithJSON(&updateUserEmail)
	if err != nil {
		slog.Error("Failed to bind request body to dto", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidJSON})
		return
	}
	err = h.validator.ValidateUpdateUserEmailReq(updateUserEmail)
	if err != nil {
		slog.Error("Failed to validate request body", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	err = h.service.UpdateUserEmail(ctx, uid, updateUserEmail)
	if err != nil {
		slog.Error("Failed to update user's email", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateUserEmailFailed})
		return
	}
	slog.Info("User's email updated")
	c.Status(http.StatusNoContent)
}

// ReadUsers handles the HTTP GET request to retrieve a list of all users.
// This endpoint is typically used by admins to retrieve user data without sensitive information like passwords.
// If successful, it returns a 200 OK status with the list of users in JSON format.
// Possible errors:
//   - 500 Internal Server Error: Failed to read users from the service layer.
func (h *Handler) ReadUsers(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := h.service.ReadUsers(ctx)
	if err != nil {
		slog.Error("Failed to read users", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": readUsersFailed})
		return
	}
	c.JSON(http.StatusOK, users)
}

// ReadUser handles the HTTP GET request to retrieve detailed information about a single user.
// This endpoint is typically used by admins to retrieve user data without sensitive information like passwords.
// If successful, it returns a 200 OK status with the user details in JSON format.
// Possible errors:
//   - 500 Internal Server Error: Failed to read user from the service layer.
func (h *Handler) ReadUser(c *gin.Context) {
	uidFromParam := c.Param(uidParam)
	err := h.validator.ValidateUID(uidFromParam)
	if err != nil {
		slog.Error("Failed to validate id param", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := uuid.Parse(uidFromParam)
	if err != nil {
		slog.Error("Failed to build UUID from id param", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	user, err := h.service.ReadUser(ctx, uid)
	if err != nil {
		slog.Error("Failed to read user", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": readUserFailed})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUserData handles the HTTP PUT request to update all user details except the UserID.
// The user ID is provided as a URL parameter.
// If successful, it returns a 204 No Content status.
// Possible errors:
//   - 400 Bad Request: Invalid ID parameter or invalid user data format.
//   - 500 Internal Server Error: Failed to update the user's data in the service layer.
func (h *Handler) UpdateUserData(c *gin.Context) {
	uidFromParam := c.Param(uidParam)
	err := h.validator.ValidateUID(uidFromParam)
	if err != nil {
		slog.Error("Failed to validate id param", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := uuid.Parse(uidFromParam)
	if err != nil {
		slog.Error("Failed to build UUID from id param", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateUserData := dto.UpdateUserDataRequest{}
	err = c.ShouldBindBodyWithJSON(&updateUserData)
	if err != nil {
		slog.Error("Failed to bind request body to dto", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidJSON})
		return
	}
	err = h.validator.ValidateUpdateUserDataReq(updateUserData)
	if err != nil {
		slog.Error("Failed to validate request body", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	err = h.service.UpdateUserData(ctx, uid, updateUserData)
	if err != nil {
		slog.Error("Failed to update user's data", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateUserDataFailed})
		return
	}
	slog.Info("User's data updated")
	c.Status(http.StatusNoContent)
}

// UpdateUserRole handles the HTTP PUT request to update a user's role by their ID.
// The user ID is provided as a URL parameter.
// If successful, it returns a 204 No Content status.
// Possible errors:
//   - 400 Bad Request: Invalid ID parameter or invalid role.
//   - 500 Internal Server Error: Failed to update the user's role in the service layer.
func (h *Handler) UpdateUserRole(c *gin.Context) {
	uidFromParam := c.Param(uidParam)
	err := h.validator.ValidateUID(uidFromParam)
	if err != nil {
		slog.Error("Failed to validate id param", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := uuid.Parse(uidFromParam)
	if err != nil {
		slog.Error("Failed to build UUID from id param", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateUserRole := dto.UpdateUserRoleRequest{}
	err = c.ShouldBindBodyWithJSON(&updateUserRole)
	if err != nil {
		slog.Error("Failed to bind request body to dto", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidJSON})
		return
	}
	err = h.validator.ValidateUpdateUserRoleReq(updateUserRole)
	if err != nil {
		slog.Error("Failed to validate request body", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	err = h.service.UpdateUserRole(ctx, uid, updateUserRole)
	if err != nil {
		slog.Error("Failed to update user's role", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateUserDataFailed})
		return
	}
	slog.Info("User's role updated")
	c.Status(http.StatusNoContent)
}

// DeleteUser handles the HTTP DELETE request to remove a user by their ID.
// The user ID is provided as a URL parameter.
// If successful, it returns a 204 No Content status.
// Possible errors:
//   - 400 Bad Request: Invalid ID parameter.
//   - 500 Internal Server Error: Failed to delete the user in the service layer.
func (h *Handler) DeleteUser(c *gin.Context) {
	uidFromParam := c.Param(uidParam)
	err := h.validator.ValidateUID(uidFromParam)
	if err != nil {
		slog.Error("Failed to validate id param", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := uuid.Parse(uidFromParam)
	if err != nil {
		slog.Error("Failed to build UUID from id param", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	err = h.service.DeleteUser(ctx, uid)
	if err != nil {
		slog.Error("Failed to delete user", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": deleteUserFailed})
		return
	}
	slog.Info("User deleted")
	c.Status(http.StatusNoContent)
}
