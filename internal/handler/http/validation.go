package handler_http

import (
	"encoding/json"
	"errors"
	"fmt"
	"ums/internal/models"

	"github.com/go-playground/validator/v10"
)

var (
	unknownValidationError = "unknown validation error"
	tagValueTemplate       = "%s=%v"
	validate               *validator.Validate
)

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// validateId validates user id with rfc4122 standard
func validateId(id string) error {
	return validate.Var(id, "required,uuid_rfc4122")
}

// validateUser validates models.User
func validateUser(user models.User) error {
	return validate.Struct(&user)
}

// getValidationErrorMessage return string representation of errors occured during models.User validation
func getValidationErrorMessage(err error) string {
	var vErrors validator.ValidationErrors
	if errors.As(err, &vErrors) {
		errMap := make(map[string]string)
		for _, tagError := range vErrors {
			tagString := fmt.Sprintf(tagValueTemplate, tagError.ActualTag(), tagError.Value())
			errMap[tagString] = tagError.Error()
		}
		outputJSON, err := json.Marshal(errMap)
		if err != nil {
			return unknownValidationError
		}
		return string(outputJSON)
	}

	return unknownValidationError
}
