package v1

import (
	goerrors "errors"
	"net/http"
	"ums/internal/validation"
)

const (
	internalError  string = "internal error occured"
	invalidJSON    string = "invalid JSON format"
	invalidUID     string = "invalid UID format"
	validationFail string = "request validation failed"
)

// errorDesc represents the type that describes an error occured while handling
// request and is to send as response.
type errorDesc struct {
	Message     string `json:"message,omitempty"`
	Description any    `json:"description,omitempty"`
}

// newErrorDesc creates the new instance of ErrorDesc.
func newErrorDesc(message string, desc any) *errorDesc {
	errDesc := &errorDesc{}
	errDesc.Message = message
	if desc == nil {
		errDesc.Description = map[string]string(nil)
	} else {
		errDesc.Description = desc
	}
	return errDesc
}

// newInternalErrorDesc creates ErrorDesc with Internal Server Error message and
// no description.
func newInternalErrorDesc() *errorDesc {
	return newErrorDesc(internalError, nil)
}

// errorResponse represents the type that describes a status code and error
// description to set as the request response.
type errorResponse struct {
	Code int
	Desc *errorDesc `json:"desc,omitempty"`
}

// mapError maps the provided error to an appropriate error response.
func mapError(err error) errorResponse {
	var (
		jsonErr  *JSONError
		validErr *ValidateError
		uidErr   *UIDError
	)

	switch {
	case goerrors.As(err, &jsonErr):
		errDesc := newErrorDesc(invalidJSON, nil)
		return errorResponse{http.StatusBadRequest, errDesc}

	case goerrors.As(err, &validErr):
		validationErrCtx := validation.UnwrapValidationErrCtx(goerrors.Unwrap(err))
		errDesc := newErrorDesc(validationFail, validationErrCtx)
		return errorResponse{http.StatusBadRequest, errDesc}

	case goerrors.As(err, &uidErr):
		validationErrCtx := validation.UnwrapValidationErrCtx(goerrors.Unwrap(err))
		errDesc := newErrorDesc(invalidUID, validationErrCtx)
		return errorResponse{http.StatusBadRequest, errDesc}
	}

	errDesc := newInternalErrorDesc()
	return errorResponse{http.StatusInternalServerError, errDesc}
}
