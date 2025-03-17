package errorx

import (
	"fmt"
)

// HTTP Status Code'ları
const (
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusConflict            = 409
	StatusUnprocessableEntity = 422
	StatusInternalServerError = 500
)

// Error yapısı
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error interface'ini implement et
func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Önceden tanımlanmış hatalar
var (
	ErrValidation = &Error{
		Code:    StatusUnprocessableEntity,
		Message: "Validation error",
	}

	ErrUnauthorized = &Error{
		Code:    StatusUnauthorized,
		Message: "Unauthorized access",
	}

	ErrForbidden = &Error{
		Code:    StatusForbidden,
		Message: "Permission denied",
	}

	ErrNotFound = &Error{
		Code:    StatusNotFound,
		Message: "Resource not found",
	}

	ErrInternal = &Error{
		Code:    StatusInternalServerError,
		Message: "Internal server error",
	}

	ErrDuplicate = &Error{
		Code:    StatusConflict,
		Message: "Resource already exists",
	}

	ErrInvalidRequest = &Error{
		Code:    StatusBadRequest,
		Message: "Invalid request",
	}

	ErrDatabaseOperation = &Error{
		Code:    StatusInternalServerError,
		Message: "Database operation failed",
	}

	ErrInvalidCredentials = &Error{
		Code:    StatusUnauthorized,
		Message: "Invalid credentials",
	}

	ErrAccountInactive = &Error{
		Code:    StatusForbidden,
		Message: "Account is inactive",
	}

	ErrPasswordHash = &Error{
		Code:    StatusInternalServerError,
		Message: "Password hashing failed",
	}

	ErrDuplicateEmail = &Error{
		Code:    StatusConflict,
		Message: "Email already exists",
	}
	ErrCacheNotInitialized = &Error{
		Code:    StatusInternalServerError,
		Message: "Cache is not initialized",
	}
	ErrKeyNotFound = &Error{
		Code:    StatusNotFound,
		Message: "Key not found in cache",
	}
	ErrInvalidValue = &Error{
		Code:    StatusUnprocessableEntity,
		Message: "Invalid value type",
	}
)

// Hata detayı eklemek için yardımcı fonksiyon
func WithDetails(err *Error, details string) *Error {
	return &Error{
		Code:    err.Code,
		Message: fmt.Sprintf("%s - %s", err.Message, details),
	}
}

// Yeni hata oluşturmak için yardımcı fonksiyon
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
