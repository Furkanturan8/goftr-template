package response

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error interface'ini implement et
func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Başarı yanıtları için yardımcı fonksiyonlar
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Success: true,
		Data:    data,
	})
}

func SuccessWithMessage(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Hata yanıtları için yardımcı fonksiyonlar
func Errorx(c *fiber.Ctx, statusCode int, err *Error) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error:   err,
	})
}

func ErrorWithMessage(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error: &Error{
			Code:    statusCode,
			Message: message,
		},
	})
}

// HTTP Status Code'ları
const (
	StatusBadRequest          = 400 // Validation, Invalid Request
	StatusUnauthorized        = 401 // Authentication Required
	StatusForbidden           = 403 // Permission Denied
	StatusNotFound            = 404 // Resource Not Found
	StatusConflict            = 409 // Duplicate Entry
	StatusUnprocessableEntity = 422 // Validation Error
	StatusInternalServerError = 500 // Internal Server Error
)

// Önceden tanımlanmış hata yanıtları
var (
	ErrValidation = &Error{
		Code:    StatusUnprocessableEntity,
		Message: "Doğrulama hatası",
	}

	ErrUnauthorized = &Error{
		Code:    StatusUnauthorized,
		Message: "Yetkisiz erişim",
	}

	ErrForbidden = &Error{
		Code:    StatusForbidden,
		Message: "Bu işlem için yetkiniz yok",
	}

	ErrNotFound = &Error{
		Code:    StatusNotFound,
		Message: "Kayıt bulunamadı",
	}

	ErrInternal = &Error{
		Code:    StatusInternalServerError,
		Message: "Sunucu hatası",
	}

	ErrDuplicate = &Error{
		Code:    StatusConflict,
		Message: "Kayıt zaten mevcut",
	}

	ErrInvalidRequest = &Error{
		Code:    StatusBadRequest,
		Message: "Geçersiz istek",
	}

	ErrDatabaseOperation = &Error{
		Code:    StatusInternalServerError,
		Message: "Veritabanı işlem hatası",
	}
)

// Hata detayı eklemek için yardımcı fonksiyon
func WithDetails(err *Error, details string) *Error {
	return &Error{
		Code:    err.Code,
		Message: err.Message,
		Details: details,
	}
}
