// Package handler_http provides HTTP handlers for the User Management System (UMS).
// It implements RESTful endpoints for creating, reading, updating, and deleting users.
// The package uses the Gin framework for routing and request handling.
package handler_http

import (
	"context"
	"log/slog"
	"net/http"
	"ums/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// errMessage represents error messages returned by the HTTP handlers.
// These messages provide meaningful feedback to clients in case of failures.
type errMessage string

var (
	// invalidJSON indicates that the request body contains invalid JSON format.
	invalidJSON errMessage = "invalid JSON format"

	// creationFailed indicates that the user creation process failed.
	creationFailed errMessage = "failed to create user"

	// updateFailed indicates that the user update process failed.
	updateFailed errMessage = "failed to update user"

	// readUsersFailed indicates that the process of reading users failed.
	readUsersFailed errMessage = "failed to read users"

	// idParam is the URL parameter name used to identify users in requests.
	idParam string = "id"
)

// UMS defines the interface for user management operations.
// It abstracts the business logic for creating, reading, updating, and deleting users.
type UMS interface {
	// CreateUser creates a new user in the system.
	// It returns the unique ID of the created user or an error if the operation fails.
	CreateUser(ctx context.Context, user models.User) (uuid.UUID, error)

	// ReadUsers retrieves a list of all users in the system.
	// It returns an empty list if no users are found.
	ReadUsers(ctx context.Context) ([]models.User, error)

	// UpdateUser updates an existing user by their ID.
	// It returns an error if the user does not exist or the update fails.
	UpdateUser(ctx context.Context, id uuid.UUID, user models.User) error

	// DeleteUser removes a user from the system by their ID.
	// It returns an error if the user does not exist or the deletion fails.
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// Handler is an HTTP handler responsible for managing user-related requests.
// It interacts with the UMS service to perform CRUD operations on users.
type Handler struct {
	// service is the interface that provides business logic for user management.
	service UMS
}

// New creates a new instance of the Handler.
// It takes an implementation of the UMS interface as a dependency.
// This allows for easy testing and decoupling of business logic from HTTP handling.
func New(service UMS) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateUser handles the HTTP POST request to create a new user.
// It expects a JSON payload in the request body containing user details.
// If successful, it returns a 201 Created status with the ID of the newly created user.
// Possible errors:
//   - 400 Bad Request: Invalid JSON or validation failure.
//   - 500 Internal Server Error: Failed to create the user in the service layer.
func (h *Handler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		slog.Error("Cannot bind body request to User model", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidJSON})
		return
	}
	err = validateUser(user)
	if err != nil {
		errorMessage := getValidationErrorMessage(err)
		slog.Error("Failed to validate request body", "err", errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}
	id, err := h.service.CreateUser(ctx, user)
	if err != nil {
		slog.Error("Failed to create user", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": creationFailed})
		return
	}
	slog.Info("User created", "id", id.String())
	c.JSON(http.StatusCreated, gin.H{"id": id.String()})
}

// ReadUsers handles the HTTP GET request to retrieve a list of all users.
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

// UpdateUser handles the HTTP PUT request to update an existing user by their ID.
// The user ID is provided as a URL parameter, and the updated user details are in the request body.
// If successful, it returns a 204 No Content status.
// Possible errors:
//   - 400 Bad Request: Invalid ID parameter or invalid JSON.
//   - 500 Internal Server Error: Failed to update the user in the service layer.
func (h *Handler) UpdateUser(c *gin.Context) {
	idParam := c.Param(idParam)
	err := validateId(idParam)
	if err != nil {
		errorMessage := getValidationErrorMessage(err)
		slog.Error("Failed to validate id param", "err", errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}
	id, err := uuid.Parse(idParam)
	if err != nil {
		slog.Error("Failed to build UUID from id param", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	user := models.User{}
	err = c.ShouldBindJSON(&user)
	if err != nil {
		slog.Error("Cannot bind body request to User model", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidJSON})
		return
	}
	err = validateUser(user)
	if err != nil {
		errorMessage := getValidationErrorMessage(err)
		slog.Error("Failed to validate request body", "err", errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}
	err = h.service.UpdateUser(ctx, id, user)
	if err != nil {
		slog.Error("Failed to update user", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateFailed})
		return
	}
	slog.Info("User updated", "id", id.String())
	c.Status(http.StatusNoContent)
}

// DeleteUser handles the HTTP DELETE request to remove a user by their ID.
// The user ID is provided as a URL parameter.
// If successful, it returns a 204 No Content status.
// Possible errors:
//   - 400 Bad Request: Invalid ID parameter.
//   - 500 Internal Server Error: Failed to delete the user in the service layer.
func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param(idParam)
	err := validateId(idParam)
	if err != nil {
		errorMessage := getValidationErrorMessage(err)
		slog.Error("Failed to validate id param", "err", errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}
	id, err := uuid.Parse(idParam)
	if err != nil {
		slog.Error("Failed to build UUID from id param", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	err = h.service.DeleteUser(ctx, id)
	if err != nil {
		slog.Error("Failed to delete user", "err", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}
