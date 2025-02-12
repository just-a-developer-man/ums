package errors

import "fmt"

// JSONError is the custom error occured when unmarshalling JSON.
type JSONError struct {
	Err error
}

func (e *JSONError) Error() string {
	return fmt.Sprintf("failed to unmarshal from JSON: %s", e.Err.Error())
}

func (e *JSONError) Unwrap() error {
	return e.Err
}

// ValidateError is the custom error occured when validating structs and vars.
type ValidateError struct {
	Err error
}

func (e *ValidateError) Error() string {
	return fmt.Sprintf("failed to validate request: %s", e.Err.Error())
}

func (e *ValidateError) Unwrap() error {
	return e.Err
}

// UIDError is the custom error occured when validating and parsing UUID from
// string.
type UIDError struct {
	Err error
}

func (e *UIDError) Error() string {
	return fmt.Sprintf("bad UUID: %s", e.Err.Error())
}

func (e *UIDError) Unwrap() error {
	return e.Err
}
