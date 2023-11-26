package utils

import "github.com/go-playground/validator/v10"

var v *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateBody(dst any) map[string]string {
	err := v.Struct(dst)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return map[string]string{"error": err.Error()}
		}

		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Error()
		}

		return errors
	}

	return nil
}
