// Общие обработчики ошибок
package shared

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewValidationError(message string) *AppError {
	return &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewAuthError(message string) *AppError {
	return &AppError{
		Code:       "AUTH_ERROR",
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s не найден", resource),
		HTTPStatus: http.StatusNotFound,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       "CONFLICT",
		Message:    message,
		HTTPStatus: http.StatusConflict,
	}
}

func NewInternalError(message string) *AppError {
	return &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
	}
}

// NewRateLimitError создаёт ошибку превышения лимита запросов
func NewRateLimitError(message string) *AppError {
	return &AppError{
		Code:       "RATE_LIMIT_ERROR",
		Message:    message,
		HTTPStatus: http.StatusTooManyRequests,
	}
}

// Функция для отправки ошибки в HTTP ответе
func WriteErrorResponse(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*AppError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.HTTPStatus)

		response := map[string]interface{}{
			"error": appErr,
		}

		json.NewEncoder(w).Encode(response)
	} else {
		// Для неизвестных ошибок
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		internalErr := NewInternalError("Внутренняя ошибка сервера")
		response := map[string]interface{}{
			"error": internalErr,
		}
		json.NewEncoder(w).Encode(response)
	}
}
