package v1

// JSONError is the custom error occured when unmarshalling JSON.
type JSONError struct {
	next error
}

func (e *JSONError) Error() string {
	return e.next.Error()
}

func (e *JSONError) Unwrap() error {
	return e.next
}

// ValidateError is the custom error occured when validating structs and vars.
type ValidateError struct {
	next error
}

func (e *ValidateError) Error() string {
	return e.next.Error()
}

func (e *ValidateError) Unwrap() error {
	return e.next
}

// UIDError is the custom error occured when validating and parsing UUID from
// string.
type UIDError struct {
	next error
}

func (e *UIDError) Error() string {
	return e.next.Error()
}

func (e *UIDError) Unwrap() error {
	return e.next
}
