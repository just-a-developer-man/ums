package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	govalidator "github.com/go-playground/validator/v10"
)

// ValidationError represents a custom error type for validation failures.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

var (
	unknownValidationError = "failed to extract validation error details"
	tagValueTemplate       = "%s=%v"

	// Precompiled regular expressions for password validation.
	spaceRegexp          = regexp.MustCompile(`\s`)
	printableASCIIRegexp = regexp.MustCompile(`^[\x20-\x7E]+$`)
	specialCharRegexp    = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	uppercaseRegexp      = regexp.MustCompile(`[A-Z]`)
	lowercaseRegexp      = regexp.MustCompile(`[a-z]`)
	digitRegexp          = regexp.MustCompile(`[0-9]`)
)

// Validator is a wrapper around go-playground/validator for validating project-related data.
type Validator struct {
	validate *govalidator.Validate
}

// NewValidator creates a new instance of Validator.
func NewValidator() (*Validator, error) {
	validate := govalidator.New(govalidator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("password", validatePassword)
	if err != nil {
		return nil, fmt.Errorf("validate.RegisterValidation: %w", err)
	}
	return &Validator{
		validate: validate,
	}, nil
}

// ValidateUID validates a user ID against the UUIDv5 standard.
func (v *Validator) ValidateUID(uid string) error {
	if err := v.validate.Var(uid, "required,uuid5"); err != nil {
		return &ValidationError{Message: "invalid UID"}
	}
	return nil
}

// ValidateStruct is a method to validate structs and return JSON formatted error messages, wrapped into error.
func (v *Validator) ValidateStruct(data interface{}) error {
	if data == nil {
		return &ValidationError{Message: "input data is nil"}
	}
	if err := v.validate.Struct(data); err != nil {
		return convertValidationErrors(err)
	}
	return nil
}

// validatePassword is a custom validator for password complexity.
func validatePassword(fl govalidator.FieldLevel) bool {
	password := fl.Field().String()

	// Check for password length
	if len(password) > 64 || len(password) < 6 {
		return false
	}

	// Check for spaces in the password
	if spaceRegexp.MatchString(password) {
		return false
	}

	// Check for non-printable or non-ASCII characters
	if !printableASCIIRegexp.MatchString(password) {
		return false
	}

	// Check for at least one special character
	if !specialCharRegexp.MatchString(password) {
		return false
	}

	// Check for at least one uppercase letter
	if !uppercaseRegexp.MatchString(password) {
		return false
	}

	// Check for at least one lowercase letter
	if !lowercaseRegexp.MatchString(password) {
		return false
	}

	// Check for at least one digit
	return digitRegexp.MatchString(password)
}

// convertValidationErrors converts validation errors into a formatted JSON string.
func convertValidationErrors(err error) error {
	var vErrors govalidator.ValidationErrors
	if errors.As(err, &vErrors) {
		errMap := make(map[string]string)
		for _, tagError := range vErrors {
			tagString := fmt.Sprintf(tagValueTemplate, tagError.ActualTag(), tagError.Value())
			errMap[tagString] = tagError.Error()
		}
		outputJSON, err := json.Marshal(errMap)
		if err != nil {
			return &ValidationError{Message: unknownValidationError}
		}
		return &ValidationError{Message: string(outputJSON)}
	}
	return &ValidationError{Message: unknownValidationError}
}
