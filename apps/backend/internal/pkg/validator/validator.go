package validator

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var phoneRegex = regexp.MustCompile(`^\+7[0-9]{10}$`)

// Validator wraps go-playground/validator
type Validator struct {
	validate *validator.Validate
}

// New creates a new Validator with custom validation rules
func New() *Validator {
	v := validator.New()

	// Register custom validation for Kazakhstan phone numbers
	v.RegisterValidation("kz_phone", validateKZPhone)

	return &Validator{validate: v}
}

// Struct validates a struct
func (v *Validator) Struct(s interface{}) error {
	return v.validate.Struct(s)
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatErrors converts validator errors to a structured format
func FormatErrors(err error) []ValidationError {
	if err == nil {
		return nil
	}

	var errors []ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: formatErrorMessage(e),
			})
		}
	}

	return errors
}

func formatErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "kz_phone":
		return "Phone must be in format +7XXXXXXXXXX"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email", e.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}

// validateKZPhone validates Kazakhstan phone number format
func validateKZPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return phoneRegex.MatchString(phone)
}
