package v1

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	unknownError string = "unknown error occured"
	invalidJSON  string = "invalid JSON format"
)

type BindJSONError struct {
	UnderlyingErr error
}

func (e *BindJSONError) Error() string {
	return fmt.Sprintf("failed to unmarshal request body from JSON: %s", e.UnderlyingErr.Error())
}

type ValidationError struct {
	Message          string `json:"message"`
	ValidationErrors error  `json:"errors"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("failed to validate request: %s", e.ValidationErrors.Error())
}

type ErrorResponse struct {
	Code    int
	Message string
}

func (h *Handler) handleError(
	c *gin.Context,
	logMessage string,
	err error,
) {
	if err != nil {
		slog.Error(logMessage, "err", err)
		errResp := mapError(err)
		c.JSON(errResp.Code, gin.H{"error": errResp.Message})
		c.Abort()
	}
}

func mapError(err error) ErrorResponse {
	var (
		bindErr  *BindJSONError
		validErr *ValidationError
	)
	switch {
	case errors.As(err, bindErr):
		return ErrorResponse{http.StatusBadRequest, invalidJSON}
	case errors.As(err, validErr):
		return ErrorResponse{http.StatusBadRequest, errors.Unp}
	}
	return ErrorResponse{http.StatusInternalServerError, unknownError}
}
