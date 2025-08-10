// HTTP обработчики (регистрация, логин)
package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"auth-service/internal/service"

	"github.com/BSanjik/LUNARY/services/common/shared"
)

type Handler struct {
	Service service.AuthServiceInterface
}

type RegisterRequest struct {
	Phone    string `json:"phone" example:"+77001234567"`
	Password string `json:"password" example:"MySecurePassword123"`
	Email    string `json:"email" example:"user@example.com"`
}

type LoginRequest struct {
	Phone    string `json:"phone" example:"+77001234567"`
	Password string `json:"password" example:"MySecurePassword123"`
}

const maxRequestSize = 1 << 20 // 1MB

// Простая функция для декодирования JSON
func (h *Handler) decodeRequest(r *http.Request, dst interface{}) error {
	// Лимит размера запроса
	r.Body = http.MaxBytesReader(nil, r.Body, maxRequestSize)

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return shared.NewValidationError("Неверный JSON")
	}

	// Применяем trim для строковых полей
	h.trimStringFields(dst)

	return nil
}

func (h *Handler) trimStringFields(v interface{}) {
	switch req := v.(type) {
	case *RegisterRequest:
		req.Phone = strings.TrimSpace(req.Phone)
		req.Password = strings.TrimSpace(req.Password)
		req.Email = strings.TrimSpace(req.Email)
	case *LoginRequest:
		req.Phone = strings.TrimSpace(req.Phone)
		req.Password = strings.TrimSpace(req.Password)
	}
}

// Регистрация
// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя и возвращает JWT токен с данными пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "Данные пользователя"
// @Success 201 {object} object
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := h.decodeRequest(r, &req); err != nil {
		shared.WriteErrorResponse(w, err)
		return
	}

	authResponse, err := h.Service.Register(req.Phone, req.Password)
	if err != nil {
		shared.WriteErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authResponse)
}

// Авторизация
// Login godoc
// @Summary Авторизация пользователя
// @Description Принимает телефон и пароль, возвращает JWT токен с данными пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Данные входа"
// @Success 200 {object} object
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := h.decodeRequest(r, &req); err != nil {
		shared.WriteErrorResponse(w, err)
		return
	}

	authResponse, err := h.Service.Login(req.Phone, req.Password)
	if err != nil {
		shared.WriteErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse)
}
