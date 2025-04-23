package errors

import (
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ValidationError represents input validation errors
	ValidationError ErrorType = "validation_error"
	// AuthenticationError represents authentication related errors
	AuthenticationError ErrorType = "authentication_error"
	// AuthorizationError represents authorization related errors
	AuthorizationError ErrorType = "authorization_error"
	// NotFoundError represents resource not found errors
	NotFoundError ErrorType = "not_found_error"
	// ConflictError represents resource conflict errors
	ConflictError ErrorType = "conflict_error"
	// InternalError represents internal server errors
	InternalError ErrorType = "internal_error"
)

// Error represents an application error
type Error struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Status  int       `json:"status"`
	Details []string  `json:"details,omitempty"`
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}

// HTTPStatusCode returns the appropriate HTTP status code for the error type
func (e *Error) HTTPStatusCode() int {
	return e.Status
}

// NewError creates a new Error
func NewError(errorType ErrorType, message string, status int, details ...string) *Error {
	return &Error{
		Type:    errorType,
		Message: message,
		Status:  status,
		Details: details,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(message string, details ...string) *Error {
	return NewError(ValidationError, message, http.StatusBadRequest, details...)
}

// NewAuthenticationError creates a new authentication error
func NewAuthenticationError(message string) *Error {
	return NewError(AuthenticationError, message, http.StatusUnauthorized)
}

// NewAuthorizationError creates a new authorization error
func NewAuthorizationError(message string) *Error {
	return NewError(AuthorizationError, message, http.StatusForbidden)
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) *Error {
	return NewError(NotFoundError, message, http.StatusNotFound)
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *Error {
	return NewError(ConflictError, message, http.StatusConflict)
}

// NewInternalError creates a new internal error
func NewInternalError(message string) *Error {
	return NewError(InternalError, message, http.StatusInternalServerError)
}

// IsError checks if an error is an Error
func IsError(err error) (*Error, bool) {
	if appErr, ok := err.(*Error); ok {
		return appErr, true
	}
	return nil, false
} 