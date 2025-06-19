package utils

import (
	"housing-api/internal/models"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) []models.ValidationError {
	var errors []models.ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var validationError models.ValidationError
			validationError.Field = err.Field()
			validationError.Message = getValidationMessage(err)
			validationError.Value = err.Value().(string)
			errors = append(errors, validationError)
		}
	}

	return errors
}

// getValidationMessage returns a human-readable validation message
func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Please provide a valid email address"
	case "min":
		return "This field must be at least " + err.Param() + " characters"
	case "max":
		return "This field must be at most " + err.Param() + " characters"
	default:
		return "This field is invalid"
	}
}