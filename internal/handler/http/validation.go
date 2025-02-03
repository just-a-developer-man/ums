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
)

func validateUser(user models.User) error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(&user); err != nil {
		return err
	}
	return nil
}

func getValidationErrorsRepresentation(err error) string {
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
