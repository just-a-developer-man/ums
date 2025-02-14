package validation

import (
	"errors"
	"fmt"
	"regexp"

	govalidator "github.com/go-playground/validator/v10"
)

var (
	// Precompiled regular expressions for password validation.
	spaceRegexp          = regexp.MustCompile(`\s`)
	printableASCIIRegexp = regexp.MustCompile(`^[\x20-\x7E]+$`)
	specialCharRegexp    = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	uppercaseRegexp      = regexp.MustCompile(`[A-Z]`)
	lowercaseRegexp      = regexp.MustCompile(`[a-z]`)
	digitRegexp          = regexp.MustCompile(`[0-9]`)
	usernameRegexp       = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
)

// Validator is a wrapper around go-playground/validator for validating project-related data.
type Validator struct {
	validate *govalidator.Validate
}

// NewValidator creates a new instance of Validator.
func NewValidator() (*Validator, error) {
	validate := govalidator.New(govalidator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("passwordComplexity", validatePassword)
	if err != nil {
		return nil, fmt.Errorf("password validate.RegisterValidation: %w", err)
	}
	err = validate.RegisterValidation("usernameFormat", validateUserName)
	if err != nil {
		return nil, fmt.Errorf("username validate.RegisterValidation: %w", err)
	}
	return &Validator{
		validate: validate,
	}, nil
}

// ValidateUID validates a user ID against the UUIDv5 standard.
func (v *Validator) ValidateUID(uid string) error {
	if err := v.validate.Var(uid, "required,uuid5"); err != nil {
		return &ValidateError{ctx: unwrapErrors(err), next: err}
	}
	return nil
}

// ValidateStruct is a method to validate structs and return JSON formatted error messages, wrapped into error.
func (v *Validator) ValidateStruct(data interface{}) error {
	if data == nil {
		return &ValidateError{next: fmt.Errorf("provided data is nil")}
	}
	if err := v.validate.Struct(data); err != nil {
		return &ValidateError{ctx: unwrapErrors(err), next: err}
	}
	return nil
}

// validatePassword is a custom validator for username format.
func validateUserName(fl govalidator.FieldLevel) bool {
	username := fl.Field().String()

	// Check for name length
	if len(username) > 64 || len(username) < 3 {
		return false
	}

	return usernameRegexp.MatchString(username)
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

// unwrapErrors extracts validation errors to the slice of ErrDesc structs.
func unwrapErrors(err error) []ValidationErrCtx {
	var vErrors govalidator.ValidationErrors
	errDescs := make([]ValidationErrCtx, 0)
	if errors.As(err, &vErrors) {
		for _, tagError := range vErrors {
			errDescs = append(errDescs, ValidationErrCtx{
				Field: tagError.Field(),
				Value: tagError.Value(),
				Tag: func() string {
					if tagError.Param() == "" {
						return tagError.ActualTag()
					}
					return tagError.ActualTag() + "=" + tagError.Param()
				}(),
				Type: tagError.Type().String(),
			})
		}
		return errDescs
	}
	return []ValidationErrCtx{}
}
