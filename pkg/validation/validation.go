package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/hacKRD0/trikona_go/pkg/errors"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	if email == "" {
		return errors.NewValidationError("email is required")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return errors.NewValidationError("invalid email format")
	}

	return nil
}

// ValidatePassword validates a password
func ValidatePassword(password string) error {
	if password == "" {
		return errors.NewValidationError("password is required")
	}

	if len(password) < 8 {
		return errors.NewValidationError("password must be at least 8 characters long")
	}

	// Check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.NewValidationError("password must contain at least one uppercase letter")
	}

	// Check for at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.NewValidationError("password must contain at least one lowercase letter")
	}

	// Check for at least one number
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.NewValidationError("password must contain at least one number")
	}

	// Check for at least one special character
	if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
		return errors.NewValidationError("password must contain at least one special character (!@#$%^&*)")
	}

	return nil
}

// ValidateName validates a name (first or last)
func ValidateName(name string, fieldName string) error {
	if name == "" {
		return errors.NewValidationError(fmt.Sprintf("%s is required", fieldName))
	}

	if len(name) < 2 {
		return errors.NewValidationError(fmt.Sprintf("%s must be at least 2 characters long", fieldName))
	}

	if !regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(name) {
		return errors.NewValidationError(fmt.Sprintf("%s can only contain letters", fieldName))
	}

	return nil
}

// ValidateLinkedInURL validates a LinkedIn URL
func ValidateLinkedInURL(url string) error {
	if url == "" {
		return nil // LinkedIn URL is optional
	}

	if !strings.HasPrefix(url, "https://www.linkedin.com/") {
		return errors.NewValidationError("invalid LinkedIn URL format")
	}

	return nil
}

// ValidateRole validates a user role
func ValidateRole(role string) error {
	validRoles := map[string]bool{
		"student":      true,
		"professional": true,
		"company":      true,
		"moderator":    true,
		"admin":        true,
		"guest":        true,
	}

	if !validRoles[role] {
		return errors.NewValidationError("invalid role")
	}

	return nil
} 