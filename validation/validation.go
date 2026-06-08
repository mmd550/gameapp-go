package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var iranianMobileRegex = regexp.MustCompile(`^(\+98|0)?9\d{9}$`)

type ValidationError struct {
	Message     string            `json:"message"`
	FieldErrors map[string]string `json:"field_errors,omitempty"`
}

func New() *validator.Validate {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	validate.RegisterValidation("persianphone", validatePhoneNumber)

	return validate
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	return iranianMobileRegex.MatchString(value)
}

func FieldErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "persianphone":
		return fmt.Sprintf("%s is not valid", e.Field())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}
