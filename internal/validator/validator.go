// Package validator provides utilities for validating data using the go-playground/validator library.
// It includes functions for validating user IDs, DTOs, and extracting detailed error messages.
package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"ums/internal/dto"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents a custom error type for validation failures.
// It encapsulates a human-readable message describing the validation issue.
type ValidationError struct {
	Message string
}

// Error implements the error interface for ValidationError.
func (e *ValidationError) Error() string {
	return e.Message
}

var (
	// unknownValidationError is a fallback message when validation error details cannot be extracted.
	unknownValidationError = "failed to extract validation error details"
	// tagValueTemplate defines the format for representing validation error tags and their values.
	tagValueTemplate = "%s=%v"

	// Precompiled regular expressions for password validation.
	spaceRegexp          = regexp.MustCompile(`\s`)
	printableASCIIRegexp = regexp.MustCompile(`^[\x20-\x7E]+$`)
	specialCharRegexp    = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	uppercaseRegexp      = regexp.MustCompile(`[A-Z]`)
	lowercaseRegexp      = regexp.MustCompile(`[a-z]`)
	digitRegexp          = regexp.MustCompile(`[0-9]`)
)

// Validator is a wrapper around go-playground/validator for validating data.
// It provides methods for validating user IDs and DTOs.
type Validator struct {
	// validate is an instance of go-playground/validator used for performing validations.
	validate *validator.Validate
}

// NewValidator creates a new instance of Validator.
// It initializes the underlying go-playground/validator with required struct validation enabled
// and registers custom validation rules such as password complexity.
func NewValidator() (*Validator, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("password", validatePassword)
	if err != nil {
		return nil, fmt.Errorf("validate.RegisterValidation: %w", err)
	}
	return &Validator{
		validate: validate,
	}, nil
}

// ValidateUID validates a user ID against the RFC4122 standard.
// If the UID is invalid, it returns a ValidationError with a descriptive message.
func (v *Validator) ValidateUID(uid string) error {
	if err := v.validate.Var(uid, "required,uuid_rfc4122"); err != nil {
		return &ValidationError{Message: "invalid UID"}
	}
	return nil
}

// validateStruct is a helper method to validate any struct and convert errors into JSON format.
func (v *Validator) validateStruct(data interface{}) error {
	if data == nil {
		return &ValidationError{Message: "input data is nil"}
	}
	if err := v.validate.Struct(data); err != nil {
		return convertValidationErrors(err)
	}
	return nil
}

// ValidateCreateUserReq validates a CreateUserRequest DTO.
// It checks the request against predefined validation rules.
// If validation fails, it returns a detailed error message in JSON format.
func (v *Validator) ValidateCreateUserReq(createUser dto.CreateUserRequest) error {
	return v.validateStruct(&createUser)
}

// ValidateCreateUserAdminReq validates a CreateUserRequestAdmin DTO.
// It checks the request against predefined validation rules.
// If validation fails, it returns a detailed error message in JSON format.
func (v *Validator) ValidateCreateUserAdminReq(createUser dto.CreateUserAdminRequest) error {
	return v.validateStruct(&createUser)
}

// ValidateUpdateUserPasswordReq validates an UpdateUserPasswordRequest DTO.
// It checks the request against predefined validation rules.
// If validation fails, it returns a detailed error message in JSON format.
func (v *Validator) ValidateUpdateUserPasswordReq(updateUserPassword dto.UpdateUserPasswordRequest) error {
	return v.validateStruct(&updateUserPassword)
}

// ValidateUpdateUserNameReq validates an UpdateUserNameRequest DTO.
// It checks the request against predefined validation rules.
// If validation fails, it returns a detailed error message in JSON format.
func (v *Validator) ValidateUpdateUserNameReq(updateUserName dto.UpdateUserNameRequest) error {
	return v.validateStruct(&updateUserName)
}

// ValidateUpdateUserEmailReq validates an UpdateUserEmailRequest DTO.
// It checks the request against predefined validation rules.
// If validation fails, it returns a detailed error message in JSON format.
func (v *Validator) ValidateUpdateUserEmailReq(updateUserEmail dto.UpdateUserEmailRequest) error {
	return v.validateStruct(&updateUserEmail)
}

// ValidateUpdateUserDataReq validates an UpdateUserDataRequest DTO.
// It checks the request against predefined validation rules.
// If validation fails, it returns a detailed error message in JSON format.
func (v *Validator) ValidateUpdateUserDataReq(updateUserData dto.UpdateUserDataRequest) error {
	return v.validateStruct(&updateUserData)
}

// ValidateUpdateUserRoleReq validates an UpdateUserRoleRequest DTO.
// It checks the request against predefined validation rules.
// If validation fails, it returns a detailed error message in JSON format.
func (v *Validator) ValidateUpdateUserRoleReq(updateUserRole dto.UpdateUserRoleRequest) error {
	return v.validateStruct(&updateUserRole)
}

// validatePassword is a custom validator for password complexity.
// It ensures that the password contains at least one special character, one uppercase letter,
// one lowercase letter, and one digit. Spaces are not allowed, and only printable ASCII characters are permitted.
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

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

// convertValidationErrors converts validation errors into a JSON string.
// It processes validation errors from go-playground/validator and formats them into a map,
// which is then serialized into JSON. If serialization fails, it returns a fallback error message.
func convertValidationErrors(err error) error {
	var vErrors validator.ValidationErrors
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
