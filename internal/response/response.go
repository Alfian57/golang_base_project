package response

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Alfian57/belajar-golang/internal/dto"
	errs "github.com/Alfian57/belajar-golang/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
	Error      any    `json:"error,omitempty"`
	Pagination any    `json:"pagination,omitempty"`
}

func WritePaginatedResponse[T any](ctx *gin.Context, statusCode int, result dto.PaginatedResult[T]) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(statusCode, Response{
		Success:    true,
		Data:       result.Data,
		Pagination: result.Pagination,
	})
}

func WriteDataResponse(ctx *gin.Context, statusCode int, data any) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func WriteMessageResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(statusCode, Response{
		Success: true,
		Message: message,
	})
}

// Improved error response handler
func WriteErrorResponse(ctx *gin.Context, err error) {
	ctx.Header("Content-Type", "application/json")

	// Handle custom AppError
	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		ctx.JSON(appErr.Code, Response{
			Success: false,
			Error:   appErr.Message,
		})
		return
	}

	// Handle validation errors
	var validationErr *errs.ValidationError
	if errors.As(err, &validationErr) {
		ctx.JSON(http.StatusUnprocessableEntity, Response{
			Success: false,
			Error:   validationErr.Errors,
		})
		return
	}

	// Handle gin validation errors
	var ginValidationErr validator.ValidationErrors
	if errors.As(err, &ginValidationErr) {
		fieldErrors := make([]errs.FieldError, len(ginValidationErr))
		for i, fe := range ginValidationErr {
			fieldErrors[i] = errs.FieldError{
				Field: toSnakeCase(fe.Field()),
				Error: validationMessage(fe),
			}
		}
		ctx.JSON(http.StatusUnprocessableEntity, Response{
			Success: false,
			Error:   fieldErrors,
		})
		return
	}

	// Default error response
	ctx.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error:   "internal server error",
	})
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return toSnakeCase(fe.Field()) + " is required"
	case "min":
		return toSnakeCase(fe.Field()) + " must be at least " + fe.Param() + " characters"
	case "max":
		return toSnakeCase(fe.Field()) + " must be at most " + fe.Param() + " characters"
	case "eqfield":
		return toSnakeCase(fe.Field()) + " must be equal to " + toSnakeCase(fe.Param())
	default:
		return toSnakeCase(fe.Field()) + " is invalid"
	}
}

func toSnakeCase(str string) string {
	var sb strings.Builder
	runes := []rune(str)

	for i := 0; i < len(runes); i++ {
		if i > 0 && isUpper(runes[i]) && (isLower(runes[i-1]) || (i+1 < len(runes) && isLower(runes[i+1]))) {
			sb.WriteRune('_')
		}
		sb.WriteRune(runes[i])
	}

	return strings.ToLower(sb.String())
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func isLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}
