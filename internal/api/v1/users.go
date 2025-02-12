package v1

import (
	"net/http"
	"ums/internal/dto"
	"ums/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	uidParam                 = "id"
	logHandleBodyFailed      = "Failed to handle request body"
	logHandleUIDFailed       = "Failed to handle UID"
	logUpdatePasswordFailed  = "Failed to update user password"
	logUpdateUserNameFailed  = "Failed to update user name"
	logUpdateUserEmailFailed = "Failed to update user email"
	logUpdateUserRoleFailed  = "Failed to update user role"
	logUpdateUserDatatFailed = "Failed to update user data"
	logReadUsersFailed       = "Failed to read users"
	logReadUserFailed        = "Failed to read user"
	logDeleteUserFailed      = "Failed to delete user"
	logUserCreated           = "User created"
	logPasswordUpdated       = "User password updated"
	logUserNameUpdated       = "User name updated"
	logUserEmailUpdated      = "User email updated"
	logUserRoleUpdated       = "User role updated"
	logUserDataUpdated       = "User data updated"
	logUsersRead             = "Users were read"
	logUserRead              = "User was read"
	logUserDeleted           = "User was deleted"
)

// CreateUser handless the HTTP POST request to create a new user.
func (h *Handler) CreateUser(c *gin.Context) {
	createUser := dto.CreateUserRequest{}
	err := h.handleBody(c, &createUser)
	if err != nil {
		h.handleError(c, logHandleBodyFailed, err)
		return
	}

	ctx := c.Request.Context()
	uid, err := h.service.CreateUser(ctx, createUser)
	if err != nil {
		h.handleError(c, "Failed to create user", err)
		return
	}

	h.handleOK(
		c,
		http.StatusCreated,
		dto.CreateUserResponse{ID: uid},
		logUserCreated,
	)
}

// CreateUser handles the HTTP POST request to create a new user and define his
// role.
func (h *Handler) CreateUserWithRole(c *gin.Context) {
	createUserWithRole := dto.CreateUserWithRoleRequest{}
	err := h.handleBody(c, &createUserWithRole)
	if err != nil {
		h.handleError(c, logHandleBodyFailed, err)
		return
	}

	ctx := c.Request.Context()
	uid, err := h.service.CreateUserWithRole(ctx, createUserWithRole)
	if err != nil {
		h.handleError(c, "Failed to create user", err)
		return
	}

	h.handleOK(
		c,
		http.StatusCreated,
		dto.CreateUserResponse{ID: uid},
		"User created",
	)
}

// UpdateUserRole handles the HTTP PUT request to update user's password in the system.
func (h *Handler) UpdateUserPassword(c *gin.Context) {
	updatePassword := dto.UpdateUserPasswordRequest{}
	h.updateUserData(c, updatePassword, logUpdatePasswordFailed, logPasswordUpdated)
}

// UpdateUserRole handles the HTTP PUT request to update user's name in the system.
func (h *Handler) UpdateUserName(c *gin.Context) {
	updateUsername := dto.UpdateUserNameRequest{}
	h.updateUserData(c, updateUsername, logUpdateUserNameFailed, logUserNameUpdated)
}

// UpdateUserRole handles the HTTP PUT request to update user's email in the system.
func (h *Handler) UpdateUserEmail(c *gin.Context) {
	updateEmail := dto.UpdateUserEmailRequest{}
	h.updateUserData(c, updateEmail, logUpdateUserEmailFailed, logUserEmailUpdated)
}

// UpdateUserRole handles the HTTP PUT request to update user's role in the system.
func (h *Handler) UpdateUserRole(c *gin.Context) {
	updateRole := dto.UpdateUserRoleRequest{}
	h.updateUserData(c, updateRole, logUpdateUserRoleFailed, logUserRoleUpdated)
}

// UpdateUserData handles the HTTP PUT request to update user's data in the system.
func (h *Handler) UpdateUserData(c *gin.Context) {
	updateData := dto.UpdateUserDataRequest{}
	h.updateUserData(c, updateData, logUpdateUserDatatFailed, logUserDataUpdated)
}

// ReadUsers handles the HTTP GET request to read information about all users in the system and returns it as
// the answere.
func (h *Handler) ReadUsers(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := h.service.ReadUsers(ctx)
	if err != nil {
		h.handleError(c, logReadUsersFailed, err)
		return
	}

	h.handleOK(c, http.StatusOK, users, logUsersRead)
}

// ReadUser  handles the HTTP GET request to read single user data from the system and returns it as the answere.
func (h *Handler) ReadUser(c *gin.Context) {
	uid, err := h.handleUIDParam(c)
	if err != nil {
		h.handleError(c, logHandleUIDFailed, err)
		return
	}

	ctx := c.Request.Context()
	users, err := h.service.ReadUser(ctx, uid)
	if err != nil {
		h.handleError(c, logReadUserFailed, err)
		return
	}

	h.handleOK(c, http.StatusOK, users, logUserRead)
}

// DeleteUser  handles the HTTP DELETE request to delete user from the system.
func (h *Handler) DeleteUser(c *gin.Context) {
	uid, err := h.handleUIDParam(c)
	if err != nil {
		h.handleError(c, logHandleUIDFailed, err)
		return
	}

	ctx := c.Request.Context()
	err = h.service.DeleteUser(ctx, uid)
	if err != nil {
		h.handleError(c, logDeleteUserFailed, err)
		return
	}

	h.handleOK(c, http.StatusNoContent, nil, logUserDeleted)
}

// updateUserData is unified handler for all update operations in the system.
func (h *Handler) updateUserData(c *gin.Context, dto any, failLogMessage string, okLogMessage string) {
	uid, err := h.handleUIDParam(c)
	if err != nil {
		h.handleError(c, logHandleUIDFailed, err)
		return
	}

	err = h.handleBody(c, dto)
	if err != nil {
		h.handleError(c, logHandleBodyFailed, err)
		return
	}

	ctx := c.Request.Context()
	err = h.service.UpdateUserData(ctx, uid, dto)
	if err != nil {
		h.handleError(c, failLogMessage, err)
		return
	}

	h.handleOK(c, http.StatusNoContent, nil, okLogMessage)
}

// handleUIDParam tryes to read UID param from request URI and validates it.
func (h *Handler) handleUIDParam(c *gin.Context) (uuid.UUID, error) {
	uidFromParam := c.Param(uidParam)
	err := h.validator.ValidateUID(uidFromParam)
	if err != nil {
		return uuid.Nil, &errors.UIDError{Err: err}
	}

	uid, err := uuid.Parse(uidFromParam)
	if err != nil {
		return uuid.Nil, &errors.UIDError{Err: err}
	}

	return uid, nil
}

// handleBody tryes to unmarshal request body to given dto, validates it and
// returns appropriate error if fails.
func (h *Handler) handleBody(c *gin.Context, dto any) error {
	err := c.ShouldBindJSON(dto)
	if err != nil {
		return &errors.JSONError{Err: err}
	}

	err = h.validator.ValidateStruct(dto)
	if err != nil {
		return &errors.ValidateError{Err: err}
	}

	return nil
}
