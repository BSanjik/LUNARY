// HTTP обработчики (регистрация, логин)
package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"auth-service/internal/service"
)

type Handler struct {
	Service *service.AuthService
}

type RegisterRequest struct {
	Phone    string `json:"phone" example:"+77001234567"`
	Password string `json:"password" example:"MySecurePassword123"`
	Email    string `json:"email" example:"user@example.com"`
}

type RegisterResponse struct {
	Token string `json:"token" example:"jwt.token.here"`
}

type LoginRequest struct {
	Phone    string `json:"phone" example:"+77001234567"`
	Password string `json:"password" example:"MySecurePassword123"`
}

type LoginResponse struct {
	Token string `json:"token" example:"jwt.token.here"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

// Регистрация
// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "Данные пользователя"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[Register] Ошибка декодирования JSON:", err)
		respondWithError(w, http.StatusBadRequest, "Неверный JSON")
		return
	}

	req.Phone = strings.TrimSpace(req.Phone)
	req.Password = strings.TrimSpace(req.Password)
	req.Email = strings.TrimSpace(req.Email)

	//phone and pass validation
	if !phoneRegex.MatchString(req.Phone) {
		respondWithError(w, http.StatusBadRequest, "Неверный формат телефона")
		return
	}
	if len(req.Password) < 12 {
		respondWithError(w, http.StatusBadRequest, "Пароль должен содержать минимум 12 символов")
		return
	}

	//JWT generation
	tokenStr, err := h.Service.Register(req.Phone, req.Password)
	if err != nil {
		log.Println("[Register] Service error:", err)
		if strings.Contains(err.Error(), "пользователь с таким номером уже существует") {
			respondWithError(w, http.StatusInternalServerError, "Ошибка регистрации пользователя")
			return
		}
	}

	//Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RegisterResponse{Token: tokenStr})

}

// Авторизация
// Login godoc
// @Summary Авторизация пользователя
// @Description Принимает телефон и пароль, возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Данные входа"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[Login] Ошибка декодирования JSON:", err)
		respondWithError(w, http.StatusBadRequest, "Неверный JSON")
		return
	}

	req.Phone = strings.TrimSpace(req.Phone)
	req.Password = strings.TrimSpace(req.Password)

	if !phoneRegex.MatchString(req.Phone) {
		respondWithError(w, http.StatusBadRequest, "Неверный формат телефона")
		return
	}

	tokenStr, err := h.Service.Login(req.Phone, req.Password)
	if err != nil {
		log.Println("[Login] Service error:", err)
		respondWithError(w, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenStr})
}
