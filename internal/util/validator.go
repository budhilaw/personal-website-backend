package util

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// GlobalValidator is the globally shared validator instance
	GlobalValidator = newValidator()
)

// Validation error messages
var validationMessages = map[string]string{
	"required":     "Field is required",
	"email":        "Must be a valid email address",
	"min":          "Must be at least %s characters long",
	"max":          "Must be at most %s characters long",
	"alphanum":     "Must contain only alphanumeric characters",
	"oneof":        "Must be one of the allowed values",
	"url":          "Must be a valid URL",
	"uuid":         "Must be a valid UUID",
	"password":     "Password must contain at least 8 characters, one uppercase letter, one lowercase letter, one number, and one special character",
	"nohtml":       "HTML code is not allowed",
	"image":        "Must be a valid image file (jpg, jpeg, png, gif)",
	"alphanumdash": "Must contain only alphanumeric characters, hyphens, or underscores",
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents a collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error implements the error interface
func (v ValidationErrors) Error() string {
	if len(v.Errors) == 0 {
		return ""
	}

	messages := make([]string, len(v.Errors))
	for i, err := range v.Errors {
		messages[i] = fmt.Sprintf("%s: %s", err.Field, err.Message)
	}

	return strings.Join(messages, "; ")
}

// newValidator creates a new validator instance with custom validations
func newValidator() *validator.Validate {
	v := validator.New()

	// Register custom validation tags
	_ = v.RegisterValidation("password", validatePassword)
	_ = v.RegisterValidation("nohtml", validateNoHTML)
	_ = v.RegisterValidation("alphanumdash", validateAlphanumDash)

	// Use JSON tag names instead of struct field names
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v
}

// Validate validates a struct against its validation tags
func Validate(s interface{}) error {
	err := GlobalValidator.Struct(s)
	if err != nil {
		// Convert validation errors to our custom format
		var validationErrors ValidationErrors

		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()
			param := err.Param()

			// Get the error message from our map
			message, exists := validationMessages[tag]
			if !exists {
				message = fmt.Sprintf("Failed validation for %s", tag)
			}

			// Replace placeholders in the message if needed
			if param != "" && strings.Contains(message, "%s") {
				message = fmt.Sprintf(message, param)
			}

			validationErrors.Errors = append(validationErrors.Errors, ValidationError{
				Field:   field,
				Message: message,
			})
		}

		return validationErrors
	}

	return nil
}

// validatePassword checks if a password meets security requirements
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateNoHTML ensures a string doesn't contain HTML tags
func validateNoHTML(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	htmlTagPattern := regexp.MustCompile(`<[^>]*>`)
	return !htmlTagPattern.MatchString(value)
}

// validateAlphanumDash ensures a string only contains alphanumeric characters, hyphens, or underscores
func validateAlphanumDash(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return pattern.MatchString(value)
}

// SanitizeHTML removes HTML tags from a string
func SanitizeHTML(input string) string {
	htmlTagPattern := regexp.MustCompile(`<[^>]*>`)
	return htmlTagPattern.ReplaceAllString(input, "")
}

// ValidateRequestBody validates a request body and sanitizes HTML if needed
func ValidateRequestBody(body interface{}) error {
	if body == nil {
		return errors.New("request body is required")
	}

	return Validate(body)
}
