package response

import (
	goerrors "errors"
	"net/http"
	"ums/internal/errors"
)

const (
	internalError  string = "internal error occured"
	invalidJSON    string = "invalid JSON format"
	invalidUID     string = "invalid UID format"
	validationFail string = "request validation failed"
)

// ErrorDesc represents the type that describes an error occured while handling
// request and is to send as response.
type ErrorDesc struct {
	Message     string `json:"message,omitempty"`
	Description any    `json:"description,omitempty"`
}

// NewErrorDesc creates the new instance of ErrorDesc.
func NewErrorDesc(message string, desc any) *ErrorDesc {
	errDesc := &ErrorDesc{}
	errDesc.Message = message
	if desc == nil {
		errDesc.Description = map[string]string(nil)
	} else {
		errDesc.Description = desc
	}
	return errDesc
}

// NewInternalErrorDesc creates ErrorDesc with Internal Server Error message and
// no description.
func NewInternalErrorDesc() *ErrorDesc {
	return NewErrorDesc(internalError, nil)
}

// ErrorResponse represents the type that describes a status code and error
// description to set as the request response.
type ErrorResponse struct {
	Code int
	Desc *ErrorDesc `json:"desc,omitempty"`
}

// MapError maps the provided error to an appropriate error response.
func MapError(err error) ErrorResponse {
	var (
		jsonErr  *errors.JSONError
		validErr *errors.ValidateError
		uidErr   *errors.UIDError
	)

	switch {
	case goerrors.As(err, &jsonErr):
		errDesc := NewErrorDesc(invalidJSON, nil)
		return ErrorResponse{http.StatusBadRequest, errDesc}

	case goerrors.As(err, &validErr):
		errDesc := NewErrorDesc(validationFail, goerrors.Unwrap(validErr))
		return ErrorResponse{http.StatusBadRequest, errDesc}

	case goerrors.As(err, &uidErr):
		errDesc := NewErrorDesc(invalidUID, nil)
		return ErrorResponse{http.StatusBadRequest, errDesc}
	}

	errDesc := NewInternalErrorDesc()
	return ErrorResponse{http.StatusInternalServerError, errDesc}
}
