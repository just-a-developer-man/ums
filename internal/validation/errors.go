package validation

import (
	"errors"
)

// ValidationErrCtx represents a validation error context.
type ValidationErrCtx struct {
	Field string `json:"field"`
	Value any    `json:"value"`
	Tag   string `json:"failed_tag"`
	Type  string `json:"type"`
}

// ValidateError represents a custom error type incapsulating error context and
// basic error.
type ValidateError struct {
	ctx  []ValidationErrCtx
	next error
}

func (e *ValidateError) Error() string {
	return e.next.Error()
}

func UnwrapValidationErrCtx(err error) []ValidationErrCtx {
	var validErr *ValidateError
	if errors.As(err, &validErr) {
		ctx := validErr.ctx
		if ctx != nil {
			return ctx
		}
	}
	return []ValidationErrCtx{}
}
