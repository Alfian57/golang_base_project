package errs

import (
	"fmt"
	"net/http"
)

// FieldError represents an error on a specific field during validation.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ValidationError represents a collection of validation errors.
type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

// AppError is a general application error with an HTTP code and message.
type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implements the error interface for AppError.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap allows error wrapping for AppError.
func (e *AppError) Unwrap() error {
	return e.Err
}

// Error implements the error interface for ValidationError.
func (e *ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "validation errors"
	}
	return e.Errors[0].Error
}

// Helper function to create a new AppError.
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Helper function to create a new ValidationError.
func NewValidationError(fieldErrors []FieldError) *ValidationError {
	return &ValidationError{Errors: fieldErrors}
}

// Helper function to create a new FieldError.
func NewFieldError(field, message string) FieldError {
	return FieldError{Field: field, Error: message}
}

// Common errors used in the application.
var (
	ErrRefreshTokenNotFound = &AppError{Code: http.StatusNotFound, Message: "refresh token not found"}
	ErrUserNotFound         = &AppError{Code: http.StatusNotFound, Message: "user not found"}
	ErrUsernameExist        = &AppError{Code: http.StatusUnprocessableEntity, Message: "username already exists"}
	ErrTodoNotFound         = &AppError{Code: http.StatusNotFound, Message: "todo not found"}
	ErrInternalServer       = &AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrBadRequest           = &AppError{Code: http.StatusBadRequest, Message: "bad request"}
	ErrUnauthorized         = &AppError{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrForbidden            = &AppError{Code: http.StatusForbidden, Message: "forbidden"}
)
