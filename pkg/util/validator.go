package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	usernameRegexp     = regexp.MustCompile(`^[a-z0-9_-]{3,20}$`)
	passwordRegexp     = regexp.MustCompile(`^.{8,}$`)
	emailRegexp        = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	slugRegexp         = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	fileNameRegexp     = regexp.MustCompile(`^[a-zA-Z0-9_-]+\.[a-zA-Z0-9]+$`)
	allowedMimeTypes   = map[string]bool{"image/jpeg": true, "image/png": true, "image/gif": true, "application/pdf": true}
	maxFileSizeBytes   = int64(5 * 1024 * 1024) // 5MB
	minPasswordLength  = 8
	maxPasswordLength  = 72 // bcrypt max
	maxUsernameLength  = 20
	maxTitleLength     = 200
	maxDescLength      = 500
	maxContentLength   = 50000
	maxFirstNameLength = 50
	maxLastNameLength  = 50
	maxBioLength       = 500
)

// ValidateUsername validates a username
func ValidateUsername(username string) error {
	username = strings.TrimSpace(strings.ToLower(username))
	if len(username) == 0 {
		return fmt.Errorf("username cannot be empty")
	}
	if len(username) > maxUsernameLength {
		return fmt.Errorf("username cannot be longer than %d characters", maxUsernameLength)
	}
	if !usernameRegexp.MatchString(username) {
		return fmt.Errorf("username can only contain lowercase letters, numbers, hyphens, and underscores, and must be between 3-20 characters")
	}
	return nil
}

// ValidatePassword validates a password
func ValidatePassword(password string) error {
	if len(password) < minPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", minPasswordLength)
	}
	if len(password) > maxPasswordLength {
		return fmt.Errorf("password cannot be longer than %d characters", maxPasswordLength)
	}
	if !passwordRegexp.MatchString(password) {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if len(email) == 0 {
		return fmt.Errorf("email cannot be empty")
	}
	if !emailRegexp.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidateSlug validates a slug
func ValidateSlug(slug string) error {
	slug = strings.TrimSpace(strings.ToLower(slug))
	if len(slug) == 0 {
		return fmt.Errorf("slug cannot be empty")
	}
	if !slugRegexp.MatchString(slug) {
		return fmt.Errorf("slug can only contain lowercase letters, numbers, and hyphens")
	}
	return nil
}

// ValidateFile validates a file upload
func ValidateFile(filename string, size int64, mimeType string) error {
	if !fileNameRegexp.MatchString(filename) {
		return fmt.Errorf("invalid filename format")
	}
	if size > maxFileSizeBytes {
		return fmt.Errorf("file size exceeds the maximum limit of 5MB")
	}
	if !allowedMimeTypes[mimeType] {
		return fmt.Errorf("unsupported file type: %s", mimeType)
	}
	return nil
}

// ValidateTitle validates a title
func ValidateTitle(title string) error {
	if len(title) == 0 {
		return fmt.Errorf("title cannot be empty")
	}
	if len(title) > maxTitleLength {
		return fmt.Errorf("title cannot be longer than %d characters", maxTitleLength)
	}
	return nil
}

// ValidateDescription validates a description
func ValidateDescription(description string) error {
	if len(description) > maxDescLength {
		return fmt.Errorf("description cannot be longer than %d characters", maxDescLength)
	}
	return nil
}

// ValidateContent validates content
func ValidateContent(content string) error {
	if len(content) == 0 {
		return fmt.Errorf("content cannot be empty")
	}
	if len(content) > maxContentLength {
		return fmt.Errorf("content cannot be longer than %d characters", maxContentLength)
	}
	return nil
}

// ValidateFirstName validates a first name
func ValidateFirstName(firstName string) error {
	if len(firstName) > maxFirstNameLength {
		return fmt.Errorf("first name cannot be longer than %d characters", maxFirstNameLength)
	}
	return nil
}

// ValidateLastName validates a last name
func ValidateLastName(lastName string) error {
	if len(lastName) > maxLastNameLength {
		return fmt.Errorf("last name cannot be longer than %d characters", maxLastNameLength)
	}
	return nil
}

// ValidateBio validates a bio
func ValidateBio(bio string) error {
	if len(bio) > maxBioLength {
		return fmt.Errorf("bio cannot be longer than %d characters", maxBioLength)
	}
	return nil
}

// ValidateStruct validates a struct using validator tags
func ValidateStruct(s interface{}) error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return err
	}
	return nil
}
